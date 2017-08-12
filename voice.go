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

func sendBuffer(vc *discordgo.VoiceConnection) {
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}
}

func playSound(s *discordgo.Session, guildID, channelID string, amt int) (err error) {
	// use gid in pointer to determine guild; use to figure out if in channel or no
	/* Check if VoiceConnections[guildID] exists */
	if _, ok := s.VoiceConnections[guildID]; ok {
		if s.VoiceConnections[guildID].ChannelID == channelID {
			vc := s.VoiceConnections[guildID]
			vc.Speaking(true)

			for i := 0; i != amt; i++ {
				sendBuffer(vc)
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
		sendBuffer(vc)
	}

	// Stop speaking
	vc.Speaking(false)

	time.Sleep(250 * time.Millisecond)

	// Disconnect from voice channel
	vc.Disconnect()

	return nil
}

func spamPing(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	vc, err := s.ChannelVoiceJoin(v.GuildID, v.ChannelID, false, true)
	if err != nil {
		return
	}

	vc.Speaking(true)
	for inChannel.amount != 0 {
		sendBuffer(vc)
	}
	vc.Speaking(false)
	vc.Disconnect()
	return
}
