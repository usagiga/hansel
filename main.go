package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"tetona/config"

	"github.com/bwmarrin/discordgo"
)

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

func receive(s *discordgo.Session, event *discordgo.MessageCreate) {
	messages, err := config.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	if event.Content == messages.StartTriggerMessage {
		// 起動時
		log.Println("開始 : インスタンス起動...")
		outputJSON, err := exec.Command("aws", "ec2", "start-instances", "--instance-ids", os.Getenv("INSTANCE_ID")).Output()
		if err != nil {
			log.Println("起動に失敗した :", err)
			return
		}

		startResponse := StartResponse{}
		if err := json.Unmarshal(outputJSON, &startResponse); err != nil {
			log.Println("起動時のレスポンスに異常 :", err)
			return
		}
		currentState := startResponse.StartingInstances[0].CurrentState.Name
		if currentState == "running" {
			log.Println("既に起動している")
			return
		}

		previousState := startResponse.StartingInstances[0].PreviousState.Name
		if currentState == "pending" && previousState == "pending" {
			log.Println("起動処理実行中")
			return
		}

		// 開始待ち
		if _, err := exec.Command("aws", "ec2", "wait", "instance-running", "--instance-ids", os.Getenv("INSTANCE_ID")).Output(); err != nil {
			log.Println("起動待ちに失敗した :", err)
			return
		}
		log.Println("正常終了 : インスタンス起動")

	} else if event.Content == messages.HibernateTriggerMessage {
		// 停止時
		log.Println("開始 : インスタンス停止...")
		outputJSON, err := exec.Command("aws", "ec2", "stop-instances", "--instance-ids", os.Getenv("INSTANCE_ID")).Output()
		if err != nil {
			log.Println("停止に失敗した :", err)
			return
		}

		stopResponse := StopResponse{}
		if err := json.Unmarshal(outputJSON, &stopResponse); err != nil {
			log.Println("停止時のレスポンスに異常 :", err)
			return
		}

		currentState := stopResponse.StoppingInstances[0].CurrentState.Name
		if currentState == "stopped" {
			log.Println("既に停止している")
			return
		}

		previousState := stopResponse.StoppingInstances[0].PreviousState.Name
		if currentState == "stopping" && previousState == "stopping" {
			log.Println("停止処理実行中")
			return
		}

		// 停止待ち
		if _, err := exec.Command("aws", "ec2", "wait", "instance-stopped", "--instance-ids", os.Getenv("INSTANCE_ID")).Output(); err != nil {
			log.Println("停止待ちに失敗した :", err)
			return
		}
		log.Println("正常終了 : インスタンス停止")
	}
}

func runDiscordBot() error {
	session, err := discordgo.New()
	if err != nil {
		return err
	}

	session.Token = "Bot " + os.Getenv("BOT_ID")

	session.AddHandler(receive)
	err = session.Open()

	if err != nil {
		log.Println("Failed : Start Bot")
		return err
	}
	log.Println("Succeeded : Start Bot")

	return nil
}

var stopBot = make(chan bool)

func main() {
	err := runDiscordBot()
	if err != nil {
		panic(err)
	}

	<-stopBot
	return
}
