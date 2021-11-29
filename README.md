# NoseGoes
### Discord Bot
#### Author: Nathaniel Rand

#### TODO:
- Core: Normalize the content before checking the command to make it case-insensitive.
- Core: Add !vthelp to return aviable commands.
- Free: Tag everyone in a specific VC with author is currently active
- Premium: 

#### Run Bot (Locally):
##### - Export Bot Token Env Var

    export BOT_TOKEN=<INSERT TOKEN HERE>

##### - Run go app calling env var

    go run main.go -t $BOT_TOKEN


#### Deployment to GCP:
##### - Set current workspace to active GCP project

    gcloud config set project discordbots-01

#### Reference: 
`https://dev.to/aurelievache/learning-go-by-examples-part-4-create-a-bot-for-discord-in-go-43cf`