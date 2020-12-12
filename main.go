package main

import (
	"fmt"
	"log"
	"os"
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

func receive(s *discordgo.Session, event *discordgo.MessageCreate) {
	messages, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	if event.Content == messages.StartTriggerMessage {
		log.Println("Start Instance")
		return

	} else if event.Content == messages.HibernateTriggerMessage {
		log.Println("Hibernate Instance")
		return

	}
	fmt.Println(event.Content)
}
