package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	// Load .env file.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Return value from key provided.
	return os.Getenv(key)
}

func main() {

	// Grab bot token env var.
	botToken := goDotEnvVariable("BOT_TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	guildID := m.Message.GuildID

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Grab message content from guild.
	content := m.Content

	if strings.Contains(content, "!vthelp") == true {
		// Build help message
		author := m.Author.Username

		commandHelpTitle := "Looks like you need a hand. Check out my stuff below... \n \n"
		commandHelp := "❔ - !vthelp : Provides a list of my commands. \n"
		commandKick := "🦶🏽 - !vtk @User : Vote to kick (disconnect from voice, not ban) the tagged user. \n"
		commandMute := "🎙️ - !vtm @User : Vote to mute the tagged user. \n"
		commandDeafen := "🎧 - !vtd @User : Vote to deafen the tagged user. \n"
		commandKiss := "💋 - !vtkiss @User : Vote to kiss the tagged user ❤️. \n"
		commandSite := "🔗 - !vtsite : Link to the VoteTo website \n"
		commandSupport := "✨ - !vtsupport : Link to the VoteTo Patreon. \n"
		commandVersion := "🤖 - !vtversion : Current VoteTo version. \n"

		message := "Whats up " + author + "\n \n" + commandHelpTitle + "COMMANDS: \n" + commandHelp + commandKick + commandMute + commandDeafen + commandKiss + "\n" + "OTHER: \n" + commandSite + commandSupport + commandVersion + "\n"

		// Reply to help request with build message above.
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!vtk") == true && strings.Contains(content, "!vtkiss") == false {
		// Trim bot command from string to grab User tagged
		trimmed := strings.TrimPrefix(content, "!vtk ")
		trimmedUser := strings.Trim(trimmed, "<@!>")

		// Build start vote message
		author := m.Author.Username
		message := author + " is voting to KICK " + trimmed + ". You have 15 seconds to vote starting now..."

		// Send start vote message
		messageVote, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}

		// Add yes reaction to vote message
		err = s.MessageReactionAdd(m.ChannelID, messageVote.ID, "✔️")
		if err != nil {
			fmt.Println(err)
		}

		// Add no reaction to vote message
		err = s.MessageReactionAdd(m.ChannelID, messageVote.ID, "❌")
		if err != nil {
			fmt.Println(err)
		}

		// Wait 15 seconds before counting the votes
		time.Sleep(15 * time.Second)

		// Count yes reactions from vote message
		yes, err2 := s.MessageReactions(m.ChannelID, messageVote.ID, "✔️", 100, "", "")
		if err2 != nil {
			fmt.Println(err)
		}

		// Count no reactions from vote message
		no, err3 := s.MessageReactions(m.ChannelID, messageVote.ID, "❌", 100, "", "")
		if err3 != nil {
			fmt.Println(err)
		}

		// Check reaction counts and return action/message based on results
		if len(yes) > len(no) {
			// GuildMemberMove(guildID string, userID string, channelID *string) (err error)
			err = s.GuildMemberMove(guildID, trimmedUser, nil)
			if err != nil {
				fmt.Println(err)
			}

			voteMessage := trimmed + " got clapped"
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) < len(no) {
			voteMessage := trimmed + " is loved"
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == 1 && len(no) == 1 {
			voteMessage := "No one cares enough to vote, " + trimmed + ". Almost worse than getting kicked..."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == len(no) {
			voteMessage := "Woah, a tie! " + trimmed + " got lucky this time..."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		}

	}

	if strings.Contains(content, "!vtm") == true {
		// Trim bot command from string to grab User tagged
		trimmed := strings.TrimPrefix(content, "!vtm ")
		trimmedUser := strings.Trim(trimmed, "<@!>")

		// Build start vote message
		author := m.Author.Username
		message := author + " is voting to MUTE " + trimmed + ". You have 15 seconds to vote starting now..."

		// Send start vote message
		messageVote, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}

		// Add yes reaction to vote message
		err = s.MessageReactionAdd(m.ChannelID, messageVote.ID, "✔️")
		if err != nil {
			fmt.Println(err)
		}

		// Add no reaction to vote message
		err = s.MessageReactionAdd(m.ChannelID, messageVote.ID, "❌")
		if err != nil {
			fmt.Println(err)
		}

		// Wait 15 seconds before counting the votes
		time.Sleep(15 * time.Second)

		// Count yes reactions from vote message
		yes, err2 := s.MessageReactions(m.ChannelID, messageVote.ID, "✔️", 100, "", "")
		if err2 != nil {
			fmt.Println(err)
		}

		// Count no reactions from vote message
		no, err3 := s.MessageReactions(m.ChannelID, messageVote.ID, "❌", 100, "", "")
		if err3 != nil {
			fmt.Println(err)
		}

		// Check reaction counts and return action/message based on results
		if len(yes) > len(no) {
			// GuildMemberMove(guildID string, userID string, channelID *string) (err error)
			err = s.GuildMemberMute(guildID, trimmedUser, true)
			if err != nil {
				fmt.Println(err)
			}

			voteMessage := "Shut up " + trimmed + ". You suck."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) < len(no) {
			voteMessage := trimmed + " sounds smexy"
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == 1 && len(no) == 1 {
			voteMessage := "No one cares enough to vote, " + trimmed + ". Almost worse than getting muted..."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == len(no) {
			voteMessage := "Woah, a tie! " + trimmed + " can still speak..."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if strings.Contains(content, "!vtd") == true {
		// Trim bot command from string to grab User tagged
		trimmed := strings.TrimPrefix(content, "!vtd ")
		trimmedUser := strings.Trim(trimmed, "<@!>")

		// Build start vote message
		author := m.Author.Username
		message := author + " is voting to DEAFEN " + trimmed + ". You have 15 seconds to vote starting now..."

		// Send start vote message
		messageVote, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}

		// Add yes reaction to vote message
		err = s.MessageReactionAdd(m.ChannelID, messageVote.ID, "✔️")
		if err != nil {
			fmt.Println(err)
		}

		// Add no reaction to vote message
		err = s.MessageReactionAdd(m.ChannelID, messageVote.ID, "❌")
		if err != nil {
			fmt.Println(err)
		}

		// Wait 15 seconds before counting the votes
		time.Sleep(15 * time.Second)

		// Count yes reactions from vote message
		yes, err2 := s.MessageReactions(m.ChannelID, messageVote.ID, "✔️", 100, "", "")
		if err2 != nil {
			fmt.Println(err)
		}

		// Count no reactions from vote message
		no, err3 := s.MessageReactions(m.ChannelID, messageVote.ID, "❌", 100, "", "")
		if err3 != nil {
			fmt.Println(err)
		}

		// Check reaction counts and return action/message based on results
		if len(yes) > len(no) {
			// GuildMemberMove(guildID string, userID string, channelID *string) (err error)
			err = s.GuildMemberDeafen(guildID, trimmedUser, true)
			if err != nil {
				fmt.Println(err)
			}

			voteMessage := trimmed + " cant hear nuttn. You can talk shit freely now..."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) < len(no) {
			voteMessage := trimmed + " hears everything"
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == 1 && len(no) == 1 {
			voteMessage := "No one cares enough to vote, " + trimmed + ". Almost worse than getting deafened..."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == len(no) {
			voteMessage := "Woah, a tie! " + trimmed + " can still hear us..."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if strings.Contains(content, "!vtkiss") == true {
		// Trim bot command from string to grab User tagged
		trimmed := strings.TrimPrefix(content, "!vtkiss ")
		// trimmedUser := strings.Trim(trimmed, "<@!>")

		// Build start vote message
		author := m.Author.Username
		message := author + " is voting to KISS " + trimmed + ". You have 15 seconds to vote starting now..."

		// Send start vote message
		messageVote, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}

		// Add yes reaction to vote message
		err = s.MessageReactionAdd(m.ChannelID, messageVote.ID, "😘")
		if err != nil {
			fmt.Println(err)
		}

		// Add no reaction to vote message
		err = s.MessageReactionAdd(m.ChannelID, messageVote.ID, "🤢")
		if err != nil {
			fmt.Println(err)
		}

		// Wait 15 seconds before counting the votes
		time.Sleep(15 * time.Second)

		// Count yes reactions from vote message
		yes, err2 := s.MessageReactions(m.ChannelID, messageVote.ID, "😘", 100, "", "")
		if err2 != nil {
			fmt.Println(err)
		}

		// Count no reactions from vote message
		no, err3 := s.MessageReactions(m.ChannelID, messageVote.ID, "🤢", 100, "", "")
		if err3 != nil {
			fmt.Println(err)
		}

		// Check reaction counts and return action/message based on results
		if len(yes) > len(no) {
			voteMessage := trimmed + " got slobbered 😘 "
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) < len(no) {
			voteMessage := trimmed + " was stood up 😔"
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == 1 && len(no) == 1 {
			voteMessage := "No one cares enough to vote, " + trimmed + ". Almost worse than not getting kissed..."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == len(no) {
			voteMessage := "Woah, a tie! " + trimmed + " Here ya go anyways 😘"
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		}

	}

}
