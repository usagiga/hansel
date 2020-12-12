package main

import (
	"fmt"
	"os"

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
	defer session.Close()

	<-stopBot

	return
}

func receive(s *discordgo.Session, event *discordgo.MessageCreate) {
	fmt.Println(event.Content)
}
