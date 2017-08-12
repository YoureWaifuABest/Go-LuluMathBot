package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func checkErrorSend(err error, m *discordgo.MessageCreate, s *discordgo.Session) bool {
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		return true
	}
	return false
}

func checkErrorPrint(err error) bool {
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return true
	}
	return false
}

func checkErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func checkRole(s *discordgo.Session, gID string, uID string, role string) bool {
	member, err := s.GuildMember(gID, uID)
	if err != nil {
		return false
	}

	/* check if user has specific role */
	for i := range member.Roles {
		if member.Roles[i] == role {
			return true
		}
	}
	return false
}

func guildFromMessage(m *discordgo.MessageCreate, s *discordgo.Session) (*discordgo.Guild, error) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return nil, err
	}

	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func findUserChannel(m *discordgo.MessageCreate, s *discordgo.Session) (string, error) {
	g, err := guildFromMessage(m, s)
	if err != nil {
		return "", err
	}

	/* Search for message sender in guild's voice states */
	for _, vstate := range g.VoiceStates {
		if vstate.UserID == m.Author.ID {
			return vstate.ChannelID, nil
		}
	}

	return "", fmt.Errorf("user is not in a channel")
}

func getArgs(m *discordgo.MessageCreate) (argv []string, argc int) {
	// necessary to write initial values into argv
	argv = make([]string, 1, 1)
	// temporary valuable to store byte; later converted into string and placed into argv
	var tread []byte
	// index starts at 0
	argc = 0
	for i := 0; i != len(m.Content); i++ {
		if m.Content[i] == ' ' {
			// Increase length & capacity of argv by 1
			t := make([]string, len(argv)+1, (cap(argv) + 1))
			copy(t, argv)
			argv = t
			// read tread into argv[argc]
			argv[argc] = string(tread[:])
			// clear tread.
			// this garbage collects all of tread, setting capacity to 0
			// might be inefficient for long arguments
			tread = nil
			argc++
		} else {
			tread = append(tread, m.Content[i])
		}
	}
	// get the last argument if there are no spaces at the end
	if m.Content[len(m.Content)-1] != ' ' {
		argv[argc] = string(tread[:])
	}
	// increment argc, so the first argument (the command) is counted
	argc++
	return
}
