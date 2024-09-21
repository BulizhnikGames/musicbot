package Interations

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/lithdew/youtube"
	"log"
	"strings"
	"time"
)

func Play(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		data := interaction.ApplicationCommandData()
		err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: data.Options[0].StringValue(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Printf("Errro responding to interaction: %s \n", err)
			return
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		data := interaction.ApplicationCommandData()
		input := data.Options[0].StringValue()
		if input == "" {
			return
		}
		choices := make([]*discordgo.ApplicationCommandOptionChoice, 0)
		if strings.HasPrefix(input, "http") || strings.HasPrefix("http", input) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  input,
				Value: input,
			})
		} else {
			names, _, err := getVideos(input, 5, 13*time.Minute)
			if err != nil {
				log.Printf("Error getting YT videos by with name %s: %s \n", input, err)
				return
			}
			log.Printf("Got %v names from search", len(*names))
			for _, name := range *names {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  name,
					Value: name,
				})
			}
		}

		err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
		if err != nil {
			log.Printf("Errro responding to interaction: %s \n", err)
			return
		}
	}
}

func getVideos(input string, amountLimit int, durationLimit time.Duration) (*[]string, *[]youtube.StreamID, error) {
	test()

	if input == "" {
		return nil, nil, errors.New("input is empty")
	}
	results, err := youtube.Search("animus vox", 0)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Got %v search results \n", results.Hits)

	names := make([]string, 0)
	ids := make([]youtube.StreamID, 0)

	cnt := 0
	for _, result := range results.Items {
		if cnt >= amountLimit {
			break
		}
		if result.LengthSeconds <= durationLimit {
			names = append(names, result.Title)
			ids = append(ids, result.ID)
			cnt++
		}
	}

	return &names, &ids, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func test() {
	results, err := youtube.Search("animus vox", 0)
	check(err)

	fmt.Printf("Got %d search result(s).\n\n", results.Hits)

	if len(results.Items) == 0 {
		check(fmt.Errorf("got zero search results"))
	}

	// Get the first search result and print out its details.

	details := results.Items[0]

	fmt.Printf(
		"ID: %q\n\nTitle: %q\nAuthor: %q\nDuration: %q\n\nView Count: %q\nLikes: %d\nDislikes: %d\n\n",
		details.ID,
		details.Title,
		details.Author,
		details.Duration,
		details.Views,
		details.Likes,
		details.Dislikes,
	)
}
