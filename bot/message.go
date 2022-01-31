package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type StatsResponse struct {
	Content   string
	Results   MatchResults
	ChannelID string
}

type MatchResults struct {
	Forfeit bool
	Records []Record
	Stats   [][]byte
}

type Record struct {
	Team string
	Wins string
}

func (s *StatsResponse) MatchEmbedTitle() string {
	if len(s.Results.Records) == 2 {
		return fmt.Sprintf("%s vs %s", s.Results.Records[0].Team, s.Results.Records[1].Team)
	}
	return ""
}

func (s *StatsResponse) MatchRecords() string {
	if len(s.Results.Records) == 2 {
		return fmt.Sprintf(
			"%s %s-%s %s",
			s.Results.Records[0].Team, s.Results.Records[0].Wins,
			s.Results.Records[1].Wins, s.Results.Records[1].Team,
		)
	}
	return ""
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	sr := StatsResponse{}

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Content) > 0 {
		parsed := strings.Split(m.Content, ":")
		if len(parsed) == 3 {
			sr.Results.Forfeit = true
		}
		for _, t := range parsed {
			parsedRecord := strings.Split(t, "-")
			if len(parsedRecord) == 2 {
				sr.Results.Records = append(
					sr.Results.Records,
					Record{
						Team: parsedRecord[0],
						Wins: parsedRecord[1],
					},
				)
			}
		}
	}

	fmt.Printf("Match submission: %s [%d attachments]\n", sr.MatchEmbedTitle(), len(m.Attachments))

	if sr.Results.Forfeit {
		fmt.Printf("%s resulted in a Forfiet\n", sr.MatchEmbedTitle())
		msg := discordgo.MessageSend{
			Content: fmt.Sprintf("%s (Match Forfeit)\n", sr.MatchRecords()),
		}
		if replyChannel == "" {
			sr.ChannelID = m.ChannelID
		} else {
			sr.ChannelID = replyChannel
		}
		_, err := s.ChannelMessageSendComplex(sr.ChannelID, &msg)
		if err != nil {
			fmt.Printf("Error sending message: %\n", err)
		}
		return
	}

	if len(m.Attachments) > 0 {
		for _, file := range m.Attachments {
			fmt.Printf("Processing attachment: %s\n", file.Filename)
			attachment, err := http.Get(file.URL)
			if err != nil {
				fmt.Printf("Error downloading attachment: %v\n", err)
			}
			defer attachment.Body.Close()
			attachmentData, err := ioutil.ReadAll(attachment.Body)
			if err != nil {
				fmt.Printf("Error reading attachment: %v\n", err)
			}
			extractedTable, err := ScanImage(attachmentData)
			tableCSV, _ := toCSV(extractedTable)
			if err != nil {
				fmt.Printf("Error extracting table from attachment: %v\n", err)
			}
			fmt.Printf("Finished processing: %s\n", file.Filename)
			sr.Results.Stats = append(sr.Results.Stats, []byte(tableCSV))
		}

		embedFields := []*discordgo.MessageEmbedField{}

		for i, matchStats := range sr.Results.Stats {
			emb := discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("Match %d", i+1),
				Value: fmt.Sprintf("```%s```", string(matchStats)),
			}
			embedFields = append(embedFields, &emb)
		}

		msg := discordgo.MessageSend{
			Content: "Match Stats",
			Embed: &discordgo.MessageEmbed{
				Title:  sr.MatchEmbedTitle(),
				Fields: embedFields,
			},
		}
		if replyChannel == "" {
			sr.ChannelID = m.ChannelID
		} else {
			sr.ChannelID = replyChannel
		}
		_, err := s.ChannelMessageSendComplex(sr.ChannelID, &msg)
		if err != nil {
			fmt.Printf("Error sending message: %v", err)
		}
		return
	}
}
