package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

/* declare buffer as thing sound is read into */
var buffer = make([][]byte, 0)
var token string

func init() {
	token = readFileToString("config", 59)
}

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("err creating Discord session,", err)
		return
	}

	// Load missing ping sound
	err = loadSound()
	if err != nil {
		fmt.Println("Error loading sound: ", err)
		return
	}

	// Register messageCreate func as a callback for MessageCreate events
	dg.AddHandler(messageCreate)

	// Open a connection to Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Loop until ctrl-c is received
	fmt.Println("bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close Discord
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var argv []string
	var argc int
	// Ignore all bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content     == "" {
		return
	}

	if m.Content[:1] == "!" {
		argv, argc = getArgs(m)
	}
	/* Fun stuff */
	if strings.EqualFold(m.Content, "can we fast travel?") {
		s.ChannelMessageSend(m.ChannelID, "Sam Sam is in a menu.")
		return
	}

	if strings.EqualFold(m.Content, "where are we going to go to?") {
		s.ChannelMessageSend(m.ChannelID, "Sanctuary.")
		return
	}

	if strings.EqualFold(m.Content, "that tasted") {
		s.ChannelMessageSend(m.ChannelID, "purple!")
		return
	}

	if strings.EqualFold(m.Content, "1v1 me") {
		s.ChannelMessageSend(m.ChannelID, "http://mrwgifs.com/wp-content/uploads/2013/07/Randy-Marsh-Ready-To-Fight-At-The-Baseball-Game-Gif-On-South-Park.gif")
	}

	if argc != 0 {
		if strings.EqualFold(argv[0], "!random") {
			if argc <= 1 {
				s.ChannelMessageSend(m.ChannelID, "Error! Include arguments pls")
				return
			}

			if argv[1] == "help" {
				var embed discordgo.MessageEmbed
				var command discordgo.MessageEmbedField
				embed.Color = 0xCC00CC
				embed.Title = "Help"
				embed.Description = "Finds a random whole number between two values"
				command.Name = "Usage"
				command.Value = "`!random number number`"
				embed.Fields = append(embed.Fields, &command)
				s.ChannelMessageSendEmbed(m.ChannelID, &embed)
			}

			if argc <= 2 {
				s.ChannelMessageSend(m.ChannelID, "Error! Include arguments pls")
				return
			}

			arg1, err := strconv.Atoi(argv[1])
			if checkErrorSend(err, m, s) {
				return
			}

			arg2, err := strconv.Atoi(argv[2])
			if checkErrorSend(err, m, s) {
				return
			}

			/* Bulk of work is done by randRangeInt function */
			rand, err := randRangeInt(arg1, arg2)
			if checkErrorSend(err, m, s) {
				return
			}

			s.ChannelMessageSend(m.ChannelID, strconv.Itoa(rand))
			return
		}

		if strings.EqualFold(argv[0], "!8ball") {
			rand, err := randInt(6)
			if checkErrorSend(err, m, s) {
				return
			}

			switch rand {
			case 0:
				s.ChannelMessageSend(m.ChannelID, "Definitely not.")
			case 1:
				s.ChannelMessageSend(m.ChannelID, "Answer is no.")
			case 2:
				s.ChannelMessageSend(m.ChannelID, "Absolutely.")
			case 3:
				s.ChannelMessageSend(m.ChannelID, "Answer is unclear.")
			case 4:
				s.ChannelMessageSend(m.ChannelID, "What?")
			case 5:
				s.ChannelMessageSend(m.ChannelID, "Undoubtedly.")
			}
			return
		}

		/* Join discord voice channel */
		if strings.EqualFold(argv[0], "!join") {
			g, err := guildFromMessage(m, s)
			if checkErrorSend(err, m, s) {
				return
			}

			cID, err := findUserChannel(m, s)
			if checkErrorSend(err, m, s) {
				return
			}

			_, err = s.ChannelVoiceJoin(g.ID, cID, false, true)
			checkErrorSend(err, m, s)
			return
		}

		if strings.EqualFold(argv[0], "!leave") {
			g, err := guildFromMessage(m, s)
			if checkErrorSend(err, m, s) {
				return
			}

			if _, ok := s.VoiceConnections[g.ID]; ok {
				s.VoiceConnections[g.ID].Disconnect()
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Error: Bot is not currently connected to a voice channel in your server!")
			return
		}

		/* Missing ping */
		if strings.EqualFold(argv[0], "!missing") {
			g, err := guildFromMessage(m, s)
			if checkErrorSend(err, m, s) {
				return
			}

			cID, err := findUserChannel(m, s)
			if checkErrorSend(err, m, s) {
				return
			}

			arg1 := 1
			if argc > 1 {
				arg1, err = strconv.Atoi(argv[1])
				if checkErrorSend(err, m, s) {
					return
				}
			}

			err = playSound(s, g.ID, cID, arg1)
			_ = checkErrorPrint(err)
			return
		}
	}

	/* Administrative Stuff */
	if argc != 0 {
		if strings.EqualFold(argv[0], "!help") {
			var embed discordgo.MessageEmbed
			var general discordgo.MessageEmbedField
			var calculations discordgo.MessageEmbedField
			var static discordgo.MessageEmbedField
			var footer discordgo.MessageEmbedFooter
			embed.Color = 0xCC00CC
			embed.Title = "Commands"
			general.Name = "General"
			general.Value = "`!help !source !license`"
			calculations.Name = "Calculations"
			calculations.Value = "`!reduction !lethality !damage`"
			static.Name = "Static Data"
			static.Value = "`!champ !item`"
			footer.Text = "Add help as an argument to any command to get help with it"
			embed.Fields = append(embed.Fields, &general)
			embed.Fields = append(embed.Fields, &calculations)
			embed.Fields = append(embed.Fields, &static)
			embed.Footer = &footer
			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
			return
		}

		if strings.EqualFold(argv[0], "!source") {
			s.ChannelMessageSend(m.ChannelID, "Github: https://github.com/YoureWaifuABest/Go-LuluMathBot")
			return
		}

		if strings.EqualFold(argv[0], "!license") {
			license := readFileToString("LICENSE", 1058)
			var embed discordgo.MessageEmbed
			embed.Color = 0xCC00CC
			embed.Title = "License"
			embed.Description = license
			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
			// include rito attribution here
			embed.Title = "Attribution"
			embed.Description = "LuluMathBot isn't endorsed by Riot games and "+
			                    "doesn't reflect the views or opinions of Riot Games or anyone " +
													"officially involved in producing or managing League of Legends. " +
													"League of Legends and Riot Games are trademarks or registered " +
													"trademarks of Riot Games, Inc. League of Legends © Riot Games, " +
													"Inc."
			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
			return
		}

		// This will eventually be separated into another bot
		if strings.EqualFold(argv[0], "!clear") {
			perms, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
			/*
			 * Checks if permission to manage messages is set
			 * Permissions in Discord are set by ORing a bunch of bits together
			 * So, to check for a certain one, a bitwise AND has to be used.
			 * If the permission is set, perms & PermissionManageMessages = 8192 (the value of PermissionManageMessages)
			 * If unset, = 0
			 */
			if perms&discordgo.PermissionManageMessages == 0 {
				s.ChannelMessageSend(m.ChannelID, "ERROR! You do not have permission to manage messages!")
				return
			}

			if argc <= 1 {
				s.ChannelMessageSend(m.ChannelID, "ERROR! Incorrect number of arguments. Try `!clear (amount)`, where amount is an integer under 100")
				return
			}

			arg1, err := strconv.Atoi(argv[1])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "ERROR! Argument is not an integer!")
				return
			}

			if arg1 > 100 {
				s.ChannelMessageSend(m.ChannelID, "ERROR! Argument cannot be greater than 100")
				return
			}

			var mesIDs []string
			messages, err := s.ChannelMessages(m.ChannelID, arg1, "", "", "")

			for i, _ := range messages {
				mesIDs = append(mesIDs, messages[i].ID)
			}

			err = s.ChannelMessagesBulkDelete(m.ChannelID, mesIDs)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "ERROR! "+err.Error())
				return
			}
			return
		}
	}

	/* Math */
	if argc != 0 {
		if strings.EqualFold(argv[0], "!reduction") {
			if argc <= 1 {
				s.ChannelMessageSend(m.ChannelID, "ERROR! Too few arguments. See `!reduction help`")
				return
			} else if argc > 2 {
				s.ChannelMessageSend(m.ChannelID, "ERROR! Too many arguments. See `!reduction help`")
				return
			}

			if argv[1] == "help" {
				var embed discordgo.MessageEmbed
				var usage discordgo.MessageEmbedField
				embed.Color = 0xCC00CC
				embed.Title = "Help"
				embed.Description = "Prints the value damage is multiplied by with a given amount of armor / mr"
				usage.Name = "Usage"
				usage.Value = "`!reduction number`"
				embed.Fields = append(embed.Fields, &usage)
				s.ChannelMessageSendEmbed(m.ChannelID, &embed)
				return
			}

			arg1, err := strconv.ParseFloat(argv[1], 64)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "ERROR! Entered value is not a number")
				return
			}

			var reduction float64
			if arg1 < 0 {
				reduction = 2 - (100 / (100 - arg1))
			} else {
				reduction = 100 / (100 + arg1)
			}
			s.ChannelMessageSend(m.ChannelID, "Reduction Multiplier:"+strconv.FormatFloat(reduction, 'G', -1, 64))
			return
		}
	}

	/* Data */
	if argc != 0 {
		if strings.EqualFold(argv[0], "!item") {
			if argc <= 1 {
				s.ChannelMessageSend(m.ChannelID, "ERROR! No arguments present.")
				return
			}

			if argv[1] == "help" {
				var embed discordgo.MessageEmbed
				var usage discordgo.MessageEmbedField
				embed.Color = 0xCC00CC
				embed.Title = "Help"
				embed.Description = "Prints information on a given item."
				usage.Name = "Usage"
				usage.Value = "`!item item`"
				embed.Fields = append(embed.Fields, &usage)
				s.ChannelMessageSendEmbed(m.ChannelID, &embed)
				return
			}
		}

		if strings.EqualFold(argv[0], "!champ") {
			if argc <= 1 {
				s.ChannelMessageSend(m.ChannelID, "ERROR! No arguments present!")
				return
			}

			if argv[1] == "help" {
				var embed discordgo.MessageEmbed
				var usage discordgo.MessageEmbedField
				var datav discordgo.MessageEmbedField
				embed.Color = 0xCC00CC
				embed.Title = "Help"
				embed.Description = "Prints data on a certain champion"
				usage.Name = "Usage"
				usage.Value = "`!champ data champion options`\n Options vary. `!champ data champion help` to list options."
				datav.Name = "Possible Values for Data"
				datav.Value = "`lore, blurb, stats, skins, tags`"
				embed.Fields = append(embed.Fields, &usage)
				embed.Fields = append(embed.Fields, &datav)
				s.ChannelMessageSendEmbed(m.ChannelID, &embed)
				return
			}
			return
		}
	}
	return
}
