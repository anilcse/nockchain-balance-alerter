# Nock Balance Monitor

A Go program that monitors blockchain address balances on the Nockblocks network, converts balances from nick to $NOCK (2^16 nick = 65,536 nick per $NOCK), and sends formatted notifications to a Slack channel. It checks balances every minute, sends alerts on changes, and provides a summary every 6 hours. Balances are stored locally in `balances.json` for comparison.

## Features
- **Balance Monitoring**: Queries balances for multiple addresses every minute via the Nockblocks RPC API (`https://nockblocks.com/rpc`).
- **Slack Notifications**: Sends beautifully formatted alerts for balance changes and periodic summaries using Slack's block kit.
- **Balance Conversion**: Converts balances from nick to $NOCK (1 $NOCK = 65,536 nick).
- **Local Storage**: Persists balances in `balances.json` to track changes across runs.
- **Multi-Address Support**: Monitors multiple addresses configured in the `.env` file.
- **Scheduling**: Uses `gocron` for reliable minute-by-minute checks and 6-hourly summaries.

## Prerequisites
- **Go**: Version 1.16 or higher.
- **Slack Workspace**: Access to a Slack workspace where you can create an app and add it to a channel.
- **Dependencies**:
  - `github.com/go-co-op/gocron`
  - `github.com/joho/godotenv`
  - `github.com/slack-go/slack`

## Setup

### 1. Install Dependencies
Run the following command to install required Go packages:
```bash
go get github.com/go-co-op/gocron
go get github.com/joho/godotenv
go get github.com/slack-go/slack
```

### 2. Create a Slack App
1. Go to [Slack API](https://api.slack.com/apps) and click "Create New App" > "From scratch."
2. Name the app (e.g., `NockBalanceBot`) and select your workspace.
3. Under **OAuth & Permissions** > **Bot Token Scopes**, add:
   - `chat:write` (to send messages).
4. Install the app to your workspace to generate a **Bot User OAuth Token** (starts with `xoxb-`).
5. Copy the token for use in the `.env` file.

### 3. Add the Bot to a Slack Channel
1. Create or choose a Slack channel (e.g., `#nock-balances`).
2. Add the bot:
   - In the channel, type `/invite @NockBalanceBot` or go to channel settings > "Integrations" > "Add an App" and select your bot.
3. Ensure the bot has permission to post messages (granted via `chat:write` scope).

### 4. Configure the Environment
Create a `.env` file in the project root with the following:
```env
SLACK_BOT_TOKEN=xoxb-your-bot-token
SLACK_CHANNEL=#nock-balances
ADDRESSES=2jFjhxTEaBw7oyCdjLNTJYQoNbBEPBPmT92rYgD33UqA3RjCSVCDrp2qnGGoiaKnCC9ERn8x3iGGco2geXRVRwRq4wTGDHSodfL4HtJ8vxuGTnbKzVbVtFVVPP6jyQswJR9r,another_address_here
```
- Replace `xoxb-your-bot-token` with the Slack bot token.
- Replace `#nock-balances` with your channel name.
- Add multiple addresses in `ADDRESSES`, separated by commas.

### 5. Run the Program
```bash
go run main.go
```
The program will:
- Check balances every minute.
- Send Slack alerts for balance changes.
- Send a summary every 6 hours.
- Store balances in `balances.json`.

## Balance Conversion
Balances are reported in both nick and $NOCK:
- **1 $NOCK = 2^16 nick = 65,536 nick**.
- Example: 34492645376 nick = 526.18 $NOCK.

## Example Slack Messages

### Balance Change Alert
```
💸 Balance Change Alert
*Address*: `2jFjhxTEaBw7oyCdjLNTJYQoNbBEPBPmT92rYgD33UqA3RjCSVCDrp2qnGGoiaKnCC9ERn8x3iGGco2geXRVRwRq4wTGDHSodfL4HtJ8vxuGTnbKzVbVtFVVPP6jyQswJR9r`
*Old Balance*: 34492645376 nick (526.18 $NOCK)
*New Balance*: 34492645377 nick (526.18 $NOCK)
---
_Updated at 2025-07-17T15:31:00Z_
```

### Balance Summary (Every 6 Hours)
```
📊 Balance Summary
*Address 1*: `2jFjhxTEaBw7oyCdjLNTJYQoNbBEPBPmT92rYgD33UqA3RjCSVCDrp2qnGGoiaKnCC9ERn8x3iGGco2geXRVRwRq4wTGDHSodfL4HtJ8vxuGTnbKzVbVtFVVPP6jyQswJR9r`
*Balance*: 34492645376 nick (526.18 $NOCK)
*Last Updated*: 2025-07-17T15:31:00Z
---
*Address 2*: `another_address_here`
*Balance*: 123456789 nick (1.88 $NOCK)
*Last Updated*: 2025-07-17T15:30:00Z
---
_Generated at 2025-07-17T15:31:00Z_
```

## Troubleshooting
- **Bot Not Posting**:
  - Verify `SLACK_BOT_TOKEN` is correct (starts with `xoxb-`) and has `chat:write` scope.
  - Ensure the bot is added to the channel (`/invite @NockBalanceBot`).
  - Check `SLACK_CHANNEL` matches the channel name (e.g., `#nock-balances`).
  - Review logs for errors (e.g., `Error sending Slack message: ...`).
- **Network Issues**:
  - Ensure access to `https://nockblocks.com/rpc` and `https://slack.com/api`.
- **Invalid Addresses**:
  - Verify addresses in `ADDRESSES` are valid and correctly formatted.
- **Program Stops**:
  - For 24/7 operation, deploy on a server (e.g., AWS, Heroku).

## Deployment
For production:
1. Deploy on a cloud platform (e.g., AWS EC2, Heroku).
2. Use a process manager like `systemd` or `pm2` to keep the program running.
3. Monitor logs for errors and ensure network connectivity.

## Security
- Keep `SLACK_BOT_TOKEN` secure. Regenerate it in the Slack API portal if compromised.
- Store the `.env` file outside version control (add to `.gitignore`).

## License
MIT License. See [LICENSE](LICENSE) for details (create if needed).
