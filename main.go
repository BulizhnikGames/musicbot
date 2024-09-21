package main

import (
	"fmt"
	"github.com/BulizhnikGames/musicbot/Interations"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	BotToken := os.Getenv("BOT_TOKEN")
	if BotToken == "" {
		log.Fatal("Bot token not found")
	}

	AppID := os.Getenv("APP_ID")
	if AppID == "" {
		log.Fatal("AppID not found")
	}

	session, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %s", err)
	}

	_, err = session.ApplicationCommandBulkOverwrite(AppID, "", []*discordgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "play YT video by name or URL",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "name",
					Description:  "name or url of the video",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Error initializing application's slash commands: %s", err)
	}

	CommandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"play": Interations.Play,
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	err = session.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: %s", err)
	}
	defer session.Close()

	fmt.Println("Bot is now running.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
