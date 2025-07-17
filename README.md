Nock Balance Monitor
A Go program that monitors blockchain address balances on the Nockblocks network, converts balances from nick to $NOCK (2^16 nick = 65,536 nick per $NOCK), and sends formatted notifications to a Slack channel. It checks balances every minute, sends alerts on changes, and provides a summary every 6 hours. Balances are stored locally in balances.json for comparison.
Features

Balance Monitoring: Queries balances for multiple addresses every minute via the Nockblocks RPC API (https://nockblocks.com/rpc).
Slack Notifications: Sends beautifully formatted alerts for balance changes and periodic summaries using Slack's block kit.
Balance Conversion: Converts balances from nick to $NOCK (1 $NOCK = 65,536 nick).
Local Storage: Persists balances in balances.json to track changes across runs.
Multi-Address Support: Monitors multiple addresses configured in the .env file.
Scheduling: Uses gocron for reliable minute-by-minute checks and 6-hourly summaries.

Prerequisites

Go: Version 1.16 or higher.
Slack Workspace: Access to a Slack workspace where you can create an app and add it to a channel.
Dependencies:
github.com/go-co-op/gocron
github.com/joho/godotenv
github.com/slack-go/slack



Setup
1. Install Dependencies
Run the following command to install required Go packages:
go get github.com/go-co-op/gocron
go get github.com/joho/godotenv
go get github.com/slack-go/slack

2. Create a Slack App

Go to Slack API and click "Create New App" > "From scratch."
Name the app (e.g., NockBalanceBot) and select your workspace.
Under OAuth & Permissions > Bot Token Scopes, add:
chat:write (to send messages).


Install the app to your workspace to generate a Bot User OAuth Token (starts with xoxb-).
Copy the token for use in the .env file.

3. Add the Bot to a Slack Channel

Create or choose a Slack channel (e.g., #nock-balances).
Add the bot:
In the channel, type /invite @NockBalanceBot or go to channel settings > "Integrations" > "Add an App" and select your bot.


Ensure the bot has permission to post messages (granted via chat:write scope).

4. Configure the Environment
Create a .env file in the project root with the following:
SLACK_BOT_TOKEN=xoxb-your-bot-token
SLACK_CHANNEL=#nock-balances
ADDRESSES=3c2fmFYwe1cfMdAoPY3gbiLcLuWbFZzqnkbHu7zLHdJCTVJHXkZSJyBWMBQiQkbmvzH7f5zrLJh8YNn2bBVa1R54zEMk5LL52d1UiwHzHUdbkVMkajvZSo5DRWkjQcPnh6Nq,another_address_here


Replace xoxb-your-bot-token with the Slack bot token.
Replace #nock-balances with your channel name.
Add multiple addresses in ADDRESSES, separated by commas.

5. Run the Program
go run main.go

The program will:

Check balances every minute.
Send Slack alerts for balance changes.
Send a summary every 6 hours.
Store balances in balances.json.

Balance Conversion
Balances are reported in both nick and $NOCK:

1 $NOCK = 2^16 nick = 65,536 nick.
Example: 34492645376 nick = 526.18 $NOCK.

Example Slack Messages
Balance Change Alert
💸 Balance Change Alert
*Address*: `3c2fmFYwe1cfMdAoPY3gbiLcLuWbFZzqnkbHu7zLHdJCTVJHXkZSJyBWMBQiQkbmvzH7f5zrLJh8YNn2bBVa1R54zEMk5LL52d1UiwHzHUdbkVMkajvZSo5DRWkjQcPnh6Nq`
*Old Balance*: 34492645376 nick (526.18 $NOCK)
*New Balance*: 34492645377 nick (526.18 $NOCK)
---
_Updated at 2025-07-17T15:31:00Z_

Balance Summary (Every 6 Hours)
📊 Balance Summary
*Address 1*: `3c2fmFYwe1cfMdAoPY3gbiLcLuWbFZzqnkbHu7zLHdJCTVJHXkZSJyBWMBQiQkbmvzH7f5zrLJh8YNn2bBVa1R54zEMk5LL52d1UiwHzHUdbkVMkajvZSo5DRWkjQcPnh6Nq`
*Balance*: 126701535244 nick (1.93M $NOCK)
*Last Updated*: 2025-07-17T15:31:00Z
---
*Address 2*: `another_address_here`
*Balance*: 123456789 nick (1.88 $NOCK)
*Last Updated*: 2025-07-17T15:30:00Z
---
_Generated at 2025-07-17T15:31:00Z_

Troubleshooting

Bot Not Posting:
Verify SLACK_BOT_TOKEN is correct (starts with xoxb-) and has chat:write scope.
Ensure the bot is added to the channel (/invite @NockBalanceBot).
Check SLACK_CHANNEL matches the channel name (e.g., #nock-balances).
Review logs for errors (e.g., Error sending Slack message: ...).


Network Issues:
Ensure access to https://nockblocks.com/rpc and https://slack.com/api.


Invalid Addresses:
Verify addresses in ADDRESSES are valid and correctly formatted.


Program Stops:
For 24/7 operation, deploy on a server (e.g., AWS, Heroku).



Deployment
For production:

Deploy on a cloud platform (e.g., AWS EC2, Heroku).
Use a process manager like systemd or pm2 to keep the program running.
Monitor logs for errors and ensure network connectivity.

Security

Keep SLACK_BOT_TOKEN secure. Regenerate it in the Slack API portal if compromised.
Store the .env file outside version control (add to .gitignore).

License
MIT License. See LICENSE for details (create if needed).