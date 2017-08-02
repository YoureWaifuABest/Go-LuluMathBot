package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

func loadSound() error {
	file, err := os.Open("missing.dca")
	if err != nil {
		fmt.Println("Error opening dca file:", err)
		return err
	}

	var opuslen int16

	for {
		// read opus frame length
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// if EOF, return
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file:", err)
			return err
		}

		// Read encoded pcm from dca
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Not EOF
		if err != nil {
			fmt.Println("Error reading from dca file:", err)
			return err
		}

		buffer = append(buffer, InBuf)
	}
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

func playSound(s *discordgo.Session, guildID, channelID string, amt int) (err error) {
	// use gid in pointer to determine guild; use to figure out if in channel or no
	/* Check if VoiceConnections[guildID] exists */
	if _, ok := s.VoiceConnections[guildID]; ok {
		if s.VoiceConnections[guildID].ChannelID == channelID {
			vc := s.VoiceConnections[guildID]
			vc.Speaking(true)

			for i := 0; i != amt; i++ {
				for _, buff := range buffer {
					vc.OpusSend <- buff
				}
			}

			vc.Speaking(false)
			return nil
		}
	}

	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)

	// Start speaking
	vc.Speaking(true)

	// Send buffer data
	for i := 0; i != amt; i++ {
		for _, buff := range buffer {
			vc.OpusSend <- buff
		}
	}

	// Stop speaking
	vc.Speaking(false)

	time.Sleep(250 * time.Millisecond)

	// Disconnect from voice channel
	vc.Disconnect()

	return nil
}
