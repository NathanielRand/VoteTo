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

var version string = "1.0.1"

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

		commandHelpTitle := "Looks like you need a hand. Check out my goodies below... \n \n"
		commandHelp := "❔  !vthelp : Provides a list of my commands. \n"
		commandKick := "🦶🏽  !vtk @User : Vote to kick (disconnect, not ban) a user. \n"
		commandMute := "🎙️  !vtm @User : Vote to mute a user. \n"
		commandUnMute := "🎙️  !vtum @User : Vote to unmute a user. \n"
		commandDeafen := "🎧  !vtd @User : Vote to deafen a user. \n"
		commandUnDeafen := "🎧  !vtud @User : Vote to undeafen a user. \n"
		commandMuteDeafen := "🔇  !vtx @User : Vote to mute & deafen a user. \n"
		commandUnMuteDeafen := "🔇  !vtux @User : Vote to unmute & undeafen a user. \n"
		commandKiss := "💋  !vtkiss @User : Vote to kiss a user ❤️. \n"
		commandSite := "🔗  !vtsite : Link to the VoteTo website \n"
		commandSupport := "✨  !vtsupport : Link to the VoteTo Patreon. \n"
		commandVersion := "🤖  !vtversion : Current VoteTo version. \n"

		message := "Whats up " + author + "\n \n" + commandHelpTitle + "COMMANDS: \n \n" + commandHelp + commandKick + commandMute + commandUnMute + commandDeafen + commandUnDeafen + commandMuteDeafen + commandUnMuteDeafen + commandKiss + "\n" + "OTHER: \n \n" + commandSite + commandSupport + commandVersion + "\n \n" + "https://www.patreon.com/BotVoteTo"

		// Reply to help request with build message above.
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!vtsite") == true {
		// Build start vote message
		author := m.Author.Username
		message := "Here ya go " + author + "..." + "\n" + "https://discordbots.dev/VoteTo"

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!vtsupport") == true {
		// Build start vote message
		author := m.Author.Username
		message := "Thanks for thinking of me " + author + " 💖." + "\n" + "https://www.patreon.com/BotVoteTo"

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, "!vtversion") == true {
		// Build start vote message
		message := "VoteBot is currently running version " + version

		// Send start vote message
		_, err := s.ChannelMessageSend(m.ChannelID, message)
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
			// GuildMemberMute(guildID string, userID string, channelID *string) (err error)
			err = s.GuildMemberMute(guildID, trimmedUser, true)
			if err != nil {
				fmt.Println(err)
			}

			voteMessage := "Words are hard " + trimmed + ". You suck."
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

	if strings.Contains(content, "!vtum") == true {
		// Trim bot command from string to grab User tagged
		trimmed := strings.TrimPrefix(content, "!vtum ")
		trimmedUser := strings.Trim(trimmed, "<@!>")

		// Build start vote message
		author := m.Author.Username
		message := author + " is voting to UNMUTE " + trimmed + ". You have 15 seconds to vote starting now..."

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
			// GuildMemberMute(guildID string, userID string, channelID *string) (err error)
			err = s.GuildMemberMute(guildID, trimmedUser, false)
			if err != nil {
				fmt.Println(err)
			}

			voteMessage := "Here, have your words " + trimmed
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) < len(no) {
			voteMessage := "Big sad" + trimmed + ", you still can't talk."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == 1 && len(no) == 1 {
			voteMessage := "No one cares enough to vote, " + trimmed + ". Almost worse than still being muted..."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == len(no) {
			voteMessage := "Woah, a tie! " + trimmed + " still can't speak though..."
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

	if strings.Contains(content, "!vtud") == true {
		// Trim bot command from string to grab User tagged
		trimmed := strings.TrimPrefix(content, "!vtud ")
		trimmedUser := strings.Trim(trimmed, "<@!>")

		// Build start vote message
		author := m.Author.Username
		message := author + " is voting to UNDEAFEN " + trimmed + ". You have 15 seconds to vote starting now..."

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
			err = s.GuildMemberDeafen(guildID, trimmedUser, false)
			if err != nil {
				fmt.Println(err)
			}

			voteMessage := "Shhh shhh be quiet... " + trimmed + " is back."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) < len(no) {
			voteMessage := trimmed + " still hears nada."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == 1 && len(no) == 1 {
			voteMessage := "No one cares enough to vote, " + trimmed + ". Almost worse than still being deafened..."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == len(no) {
			voteMessage := "Woah, a tie! " + trimmed + " still can't hear though..."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if strings.Contains(content, "!vtx") == true {
		// Trim bot command from string to grab User tagged
		trimmed := strings.TrimPrefix(content, "!vtx ")
		trimmedUser := strings.Trim(trimmed, "<@!>")

		// Build start vote message
		author := m.Author.Username
		message := author + " is voting to MUTE & DEAFEN " + trimmed + ". You have 15 seconds to vote starting now..."

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

			err = s.GuildMemberMute(guildID, trimmedUser, true)
			if err != nil {
				fmt.Println(err)
			}

			voteMessage := trimmed + " cant hear or say nuttn. What a loser..."
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
			voteMessage := "No one cares enough to vote, " + trimmed + ". Almost worse than getting muted and deafened..."
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

	if strings.Contains(content, "!vtux") == true {
		// Trim bot command from string to grab User tagged
		trimmed := strings.TrimPrefix(content, "!vtux ")
		trimmedUser := strings.Trim(trimmed, "<@!>")

		// Build start vote message
		author := m.Author.Username
		message := author + " is voting to UNMUTE & UNDEAFEN " + trimmed + ". You have 15 seconds to vote starting now..."

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
			err = s.GuildMemberDeafen(guildID, trimmedUser, false)
			if err != nil {
				fmt.Println(err)
			}

			err = s.GuildMemberMute(guildID, trimmedUser, false)
			if err != nil {
				fmt.Println(err)
			}

			voteMessage := trimmed + " is back baby!"
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) < len(no) {
			voteMessage := trimmed + " still can't speak or hear."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == 1 && len(no) == 1 {
			voteMessage := "No one cares enough to vote, " + trimmed + ".  Almost worse than still being muted and deafened..."
			_, err := s.ChannelMessageSendReply(m.ChannelID, voteMessage, m.Reference())
			if err != nil {
				fmt.Println(err)
			}
		} else if len(yes) == len(no) {
			voteMessage := "Woah, a tie! " + trimmed + " still can't speak or hear though..."
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
