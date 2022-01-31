package main

import (
	"flag"
	"fmt"

	"os"
	"os/signal"

	"syscall"

	"github.com/bwmarrin/discordgo"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.StringVar(&replyChannel, "c", "", "channel to send responses to")
	flag.Parse()
}

var token string
var replyChannel string
var buffer = make([][]byte, 0)

func main() {

	if token == "" {
		fmt.Println("No token provided. Please run: airhorn -t <bot token>")
		return
	}
	if replyChannel != "" {
		fmt.Printf("[Configuration] Sending replies to channel ID: %s\n", replyChannel)
	}
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events.
	dg.AddHandler(ready)

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(MessageCreate)

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Printf("KoCStatBot is now running.  Press CTRL-C to exit.\n")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateGameStatus(0, "!airhorn")
}
