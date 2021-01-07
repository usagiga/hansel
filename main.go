package main

import (
	"encoding/json"
	"hansel/config"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ServerStatusResponse EC2起動後のステータス確認レスポンス
type ServerStatusResponse struct {
	PublicIP string `json:"publicip"`
}

// StartResponse EC2起動指示時のレスポンス
type StartResponse struct {
	StartingInstances []InstanceStatus `json:"StartingInstances"`
}

// StopResponse EC2停止指示時のレスポンス
type StopResponse struct {
	StoppingInstances []InstanceStatus `json:"StoppingInstances"`
}

// InstanceStatus EC2指示時の共通レスポンス
type InstanceStatus struct {
	InstanceID   string `json:"InstanceId"`
	CurrentState struct {
		Code int    `json:"Code"`
		Name string `json:"Name"`
	} `json:"CurrentState"`
	PreviousState struct {
		Code int    `json:"Code"`
		Name string `json:"Name"`
	} `json:"PreviousState"`
}

// TargetChannel Botがメッセージを投稿するDiscordチャンネル
type TargetChannel struct {
	s     *discordgo.Session
	event *discordgo.MessageCreate
}

func (tc *TargetChannel) messageSend(message string) error {
	// コマンドが投稿されたチャンネル
	targetChannel, err := tc.s.State.Channel(tc.event.ChannelID)
	if err != nil {
		log.Println("チャンネルの取得に失敗 :", err)
		return err
	}

	// Botからメッセージ投稿
	if _, err := tc.s.ChannelMessageSend(targetChannel.ID, message); err != nil {
		log.Println("チャンネルメッセージの送信に失敗 :", err)
		return err
	}
	return nil
}

// getIPAddress インスタンスのIPアドレス取得
func getIPAddress() (string, error) {
	statusOutputJSON, err := exec.Command("aws", "ec2", "describe-instances", "--instance-ids", os.Getenv("INSTANCE_ID"), "--query", "Reservations[].Instances[].{publicip:PublicIpAddress}").Output()
	if err != nil {
		log.Println("IPアドレス取得時、コマンド実行に失敗 : ", err)
		return "", err
	}

	ssResponse := []ServerStatusResponse{}
	if err := json.Unmarshal(statusOutputJSON, &ssResponse); err != nil {
		log.Println("IPアドレス取得時のレスポンスに異常 :", err)
		return "", err
	}

	ipaddress := ssResponse[0].PublicIP
	if ipaddress != "" {
		log.Println("IPアドレス : ", ipaddress)
	}

	return ipaddress, nil
}

func receive(s *discordgo.Session, event *discordgo.MessageCreate) {
	targetChannel := TargetChannel{
		s:     s,
		event: event,
	}

	messages, err := config.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	if event.Content == messages.StartTriggerMessage {
		// 起動時
		log.Println("開始 : インスタンス起動...")
		targetChannel.messageSend("インスタンスの起動コマンドを検知")

		outputJSON, err := exec.Command("aws", "ec2", "start-instances", "--instance-ids", os.Getenv("INSTANCE_ID")).Output()
		if err != nil {
			log.Println("起動に失敗した :", err)
			targetChannel.messageSend("インスタンスの起動に失敗")
			return
		}

		startResponse := StartResponse{}
		if err := json.Unmarshal(outputJSON, &startResponse); err != nil {
			log.Println("起動時のレスポンスに異常 :", err)
			targetChannel.messageSend("インスタンスの起動に失敗")
			return
		}
		currentState := startResponse.StartingInstances[0].CurrentState.Name
		if currentState == "running" {
			log.Println("既に起動している")
			targetChannel.messageSend("インスタンスは起動済み")
			return
		}

		previousState := startResponse.StartingInstances[0].PreviousState.Name
		if currentState == "pending" && previousState == "pending" {
			log.Println("起動処理実行中")
			targetChannel.messageSend("インスタンスは既に起動準備中")
			return
		}

		// 開始待ち
		if _, err := exec.Command("aws", "ec2", "wait", "instance-running", "--instance-ids", os.Getenv("INSTANCE_ID")).Output(); err != nil {
			log.Println("起動待ちに失敗した :", err)
			targetChannel.messageSend("インスタンスの起動状態不明　再度のコマンド入力を要求")
			return
		}

		log.Println("正常終了 : インスタンス起動")
		targetChannel.messageSend("インスタンスの起動に成功")

		// IPアドレス通知
		log.Println("IPアドレス取得待機中...")
		time.Sleep(time.Second)

		ipaddress, err := getIPAddress()
		if err != nil {
			targetChannel.messageSend("IPアドレスの取得に失敗")
			return
		}

		targetChannel.messageSend("今回のIPアドレス : " + ipaddress)

	} else if event.Content == messages.HibernateTriggerMessage {
		// 停止時
		log.Println("開始 : インスタンス停止...")
		targetChannel.messageSend("インスタンスの停止コマンドを検知")

		outputJSON, err := exec.Command("aws", "ec2", "stop-instances", "--instance-ids", os.Getenv("INSTANCE_ID")).Output()
		if err != nil {
			log.Println("停止に失敗した :", err)
			targetChannel.messageSend("インスタンスの停止に失敗")
			return
		}

		stopResponse := StopResponse{}
		if err := json.Unmarshal(outputJSON, &stopResponse); err != nil {
			log.Println("停止時のレスポンスに異常 :", err)
			targetChannel.messageSend("インスタンスの停止に失敗")
			return
		}

		currentState := stopResponse.StoppingInstances[0].CurrentState.Name
		if currentState == "stopped" {
			log.Println("既に停止している")
			targetChannel.messageSend("インスタンスは停止済み")
			return
		}

		previousState := stopResponse.StoppingInstances[0].PreviousState.Name
		if currentState == "stopping" && previousState == "stopping" {
			log.Println("停止処理実行中")
			targetChannel.messageSend("インスタンスは既に停止準備中")
			return
		}

		// 停止待ち
		if _, err := exec.Command("aws", "ec2", "wait", "instance-stopped", "--instance-ids", os.Getenv("INSTANCE_ID")).Output(); err != nil {
			log.Println("停止待ちに失敗した :", err)
			targetChannel.messageSend("インスタンスの停止状態不明　再度のコマンド入力を要求")
			return
		}

		log.Println("正常終了 : インスタンス停止")
		targetChannel.messageSend("インスタンスの停止に成功")

	} else if event.Content == messages.GetStatusTriggerMessage {
		// 起動状態の確認(IPアドレスの取得)
		log.Println("開始 : インスタンスステータス確認")
		targetChannel.messageSend("インスタンスの確認コマンドを検知")

		ipaddress, err := getIPAddress()
		if err != nil {
			targetChannel.messageSend("インスタンスの確認に失敗")
			return
		}

		if ipaddress != "" {
			targetChannel.messageSend("インスタンスは起動済み :" + ipaddress)
		} else {
			targetChannel.messageSend("インスタンスは未起動")
		}

	}
}

func runDiscordBot(session *discordgo.Session) error {
	session.Token = "Bot " + os.Getenv("BOT_ID")

	session.AddHandler(receive)
	err := session.Open()
	if err != nil {
		log.Println("Failed : Start Bot")
		return err
	}
	log.Println("Succeeded : Start Bot")

	return nil
}

func main() {
	session, err := discordgo.New()
	if err != nil {
		panic(err)
	}

	err = runDiscordBot(session)
	if err != nil {
		panic(err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer session.Close()

	// 終了を待機
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	select {
	case <-signalChan:
		log.Println("bye")
		return
	}
}
