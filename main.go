package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"tetona/config"

	"github.com/bwmarrin/discordgo"
)

var (
	stopBot = make(chan bool)
)

func main() {
	session, err := discordgo.New()
	if err != nil {
		panic(err)
	}
	session.Token = "Bot " + os.Getenv("BOT_ID")

	session.AddHandler(receive)
	err = session.Open()
	if err != nil {
		fmt.Println(err)
	}
	log.Println("Start Bot")

	// defer func(session *discordgo.Session) {
	// 	session.Close()
	// 	log.Println("Stop Bot")
	// }(session)

	<-stopBot

	return
}

type StartResponse struct {
	StartingInstances []InstanceStatus `json:"StartingInstances"`
}

type StopResponse struct {
	StoppingInstances []InstanceStatus `json:"StoppingInstances"`
}

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
		panic(err)
	}

	if event.Content == messages.StartTriggerMessage {
		log.Println("Start Instance")
		outputJSON, err := exec.Command("aws", "ec2", "start-instances", "--instance-ids", os.Getenv("INSTANCE_ID")).Output()
		if err != nil {
			panic(err)
		}

		startResponse := StartResponse{}
		if err := json.Unmarshal(outputJSON, &startResponse); err != nil {
			panic(err)
		}
		return

	} else if event.Content == messages.HibernateTriggerMessage {
		log.Println("Hibernate Instance")
		outputJSON, err := exec.Command("aws", "ec2", "stop-instances", "--instance-ids", os.Getenv("INSTANCE_ID")).Output()
		if err != nil {
			panic(err)
		}

		stopResponse := StopResponse{}
		if err := json.Unmarshal(outputJSON, &stopResponse); err != nil {
			panic(err)
		}
		return

	}
	fmt.Println(event.Content)
}
