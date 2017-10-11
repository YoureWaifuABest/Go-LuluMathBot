package main

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

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

func runMath(mathExpr string) (float64, error) {
	/* See https://golang.org/pkg/text/scanner/ */
	var s scanner.Scanner
	tokens := make([]string, 0, 0)

	s.Init(strings.NewReader(mathExpr))
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		tokens = append(tokens, s.TokenText())
	}

	/* Run shuntingYard to convert infix to reverse polish */
	eval, err := shuntingYard(tokens)
	if err != nil {
		return -1, err
	}

	/* Calculate final result, return it */
	return calc(eval)
}

/*
 * Probably replace this with scanner eventually
 * Since it obviously handles tokenizing much better
 * Not sure about how quoting etc. is handled by it
 */
func getArgs(m *discordgo.MessageCreate) (argv []string, argc int, err error) {
	/* necessary to write initial values into argv */
	argv = make([]string, 1, 1)
	/* temporary variable to store byte; later converted into string and placed into argv */
	var tread []byte
	/* Variable to track if in quote */
	inQuote := false

	/* index starts at 0 */
	argc = 0
	for i := 0; i != len(m.Content); i++ {
		if m.Content[i] == ' ' && inQuote != true {
			/* Increase length & capacity of argv by 1 */
			t := make([]string, len(argv)+1, (cap(argv) + 1))
			copy(t, argv)
			argv = t
			/* read tread into argv[argc] */
			argv[argc] = string(tread[:])
			/*
			 * clear tread.
			 * this garbage collects all of tread, setting capacity to 0
			 * might be inefficient for long arguments
			 */
			tread = nil
			argc++
		} else {
			/*
			 * Not closing / opening quotes correctly is currently undefined behavior
			 * Largely subject to change
			 * If quotes aren't closed; whole command after quote is counted as one argument
			 * if ends in quote, ?
			 */
			if inQuote == false && m.Content[i] == '"' {
				inQuote = true
			} else if m.Content[i] == '"' {
				inQuote = false
			} else if m.Content[i] == '$' && len(m.Content) >= i+3 {
				if m.Content[i:i+3] == "$((" {
					for n, _ := range m.Content {
						if len(m.Content) > n+1 && m.Content[n] == ')' && m.Content[n+1] == ')' {
							var result float64
							result, err = runMath(m.Content[i+3 : n])
							if err != nil {
								argv = nil
								argc = -1
								return
							}
							resultBytes := []byte(strconv.FormatFloat(result, 'g', -1, 64))
							for _, c := range resultBytes {
								tread = append(tread, c)
							}
							i = n + 1
						}
					}
				}
			} else {
				tread = append(tread, m.Content[i])
			}
		}
	}
	/* get the last argument if there are no spaces at the end */
	if m.Content[len(m.Content)-1] != ' ' {
		argv[argc] = string(tread[:])
	}
	/* increment argc, so the first argument (the command) is counted */
	argc++
	return
}

func findBetween(s, startstr, endstr string) (string, int) {
	start := strings.Index(s, startstr)
	end := strings.Index(s, endstr)
	if start == -1 {
		return "", -1
	}
	if end == -1 {
		return "", -1
	}

	var holder []byte
	for i := start + 1; i != end; i++ {
		holder = append(holder, s[i])
	}
	return string(holder), start
}
