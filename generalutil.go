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
