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

func receive(s *discordgo.Session, event *discordgo.MessageCreate) {
	messages, err := config.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	if event.Content == messages.StartTriggerMessage {
		log.Println("Start Instance")
		outputJSON, err := exec.Command("aws", "ec2", "start-instances", "--instance-ids", os.Getenv("INSTANCE_ID")).Output()
		if err != nil {
			log.Fatalln(err)
		}

		startResponse := StartResponse{}
		if err := json.Unmarshal(outputJSON, &startResponse); err != nil {
			log.Fatalln(err)
		}

	} else if event.Content == messages.HibernateTriggerMessage {
		log.Println("Hibernate Instance")
		outputJSON, err := exec.Command("aws", "ec2", "stop-instances", "--instance-ids", os.Getenv("INSTANCE_ID")).Output()
		if err != nil {
			log.Fatalln(err)
		}

		stopResponse := StopResponse{}
		if err := json.Unmarshal(outputJSON, &stopResponse); err != nil {
			log.Fatalln(err)
		}
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
