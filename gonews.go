package main

import (
	"fmt"
	"strings"

	"log"
	"os"

	"github.com/thoj/go-ircevent"
)

const server = "irc.snoonet.org:6667"
const admin = "Moter8"
const prefixCharacter = "@"

func main() {

	joinChannel := ""
	if len(os.Args) > 1 {
		joinChannel = "#" + os.Args[1]
	} else {
		log.Fatal("Invalid or no channel name given")
	}

	ircnick := "Dachnews"
	irccon := irc.IRC(ircnick, ircnick)
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = false

	irccon.AddCallback("001", func(e *irc.Event) {
		irccon.Join(joinChannel)
	})

	irccon.AddCallback("366", func(e *irc.Event) {})

	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {

		// Only handle messages that start with the prefixCharacter
		if e.Message()[:1] != prefixCharacter {
			return
		}

		go func(e *irc.Event) {

			// Returns the channel name when message is recieved in a channel, otherwise returns the bot's nick,
			// which is why we have to manually set the destination in case of a query.
			channel := e.Arguments[0]
			if e.Arguments[0] == irccon.GetNick() {
				channel = e.Nick
			}

			commands := strings.Split(e.Arguments[1], " ")

			switch commands[0][1:] {
			case "news", "dachnews", "dnews", "tagesschau", "nachrichten":
				{
					irccon.Privmsg(channel, "Tra-la-la die News sind da: http://www.tagesschau.de/sendung/tagesschau/index.html")
				}
			case "live", "ard":
				{
					irccon.Privmsg(channel, "Tra-la-la die News sind live: http://www.tagesschau.de/multimedia/livestreams/livestream3/index.html")
				}
			case "zib", "anews":
				{
					irccon.Privmsg(channel, "Tra-la-la die News sind da: http://tvthek.orf.at/profile/ZIB-1/1203")
				}
			case "join":
				{
					if e.Nick == admin && len(commands) >= 2 {
						// todo: verify if it's a valid channel name?
						irccon.Join(commands[1])
					}
				}
			case "leave", "part":
				{
					if e.Nick == admin && channel != e.Nick {
						irccon.Part(channel)
					}
				}
			case "quit":
				{
					if e.Nick == admin {
						irccon.QuitMessage = "Ayylmao"
						irccon.Quit()
					}
				}
			default:
				{
					irccon.Privmsg(e.Nick, "Command not recognized.")
				}
			}
		}(e)
	})

	err := irccon.Connect(server)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	irccon.Loop()
}
