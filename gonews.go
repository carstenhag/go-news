package main

import (
	"fmt"
	"strings"

	"github.com/thoj/go-ircevent"
	"os"
	"log"
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

		if e.Message()[:1] != prefixCharacter {
			return
		}

		go func(e *irc.Event) {

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
			case "join":
				{
					if e.Nick == admin && len(commands) >= 2 {
						// verify if it's a valid channel name?
						irccon.Join(commands[1])
					}
				}
			case "leave":
				{
					if e.Nick == admin {
						// don't try to part queries
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
