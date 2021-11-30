# NoseGoes
### Discord Bot
#### Author: Nathaniel Rand

#### TODO:
- Core: Normalize the content before checking the command to make it case-insensitive.
- Core: Add !vtrequest to request a new command.
- Core: Consider replcaing "!" command prefix with "/" for autocomplete. Look into how to add message.
- Core: Add timer flag option (15s, 1h, 7d, etc) where it reverses the action in the vote if voted yes. No flag is infinite and requires admin unmute/undeafen or vote to unmute/und
- Free: Tag everyone in a specific VC with author is currently active
- Premium: n/a

#### Run Bot (Locally):
##### - Add .env file with bot token

    BOT_TOKEN=<INSERT TOKEN HERE>

##### - Run go app calling env var

    go run main.go

#### Reference: 
`https://dev.to/aurelievache/learning-go-by-examples-part-4-create-a-bot-for-discord-in-go-43cf`