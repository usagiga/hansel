package entity

import "github.com/bwmarrin/discordgo"

type ServerStatus int

const (
	Wake ServerStatus = 1 << iota
	Sleep
	Going

	// Usage
	// Going | Wake : Going to Sleep
)

func (status ServerStatus) GetPresence() (pres *discordgo.Presence) {
	switch status {
	case Wake:
		fallthrough
	case Going | Wake:
		return &discordgo.Presence{
			Status: "online",
			Activities: []*discordgo.Game{
				{
					Type: discordgo.GameTypeGame,
					Name: "Minecraft",
					ApplicationID: "356875570916753438",
				},
			},
		}
	case Sleep:
		fallthrough
	case Going | Sleep:
		return &discordgo.Presence{
			Status: "idle",
		}
	default:
		panic("ServerStatus.GetPresence(): invalid status")
	}
}