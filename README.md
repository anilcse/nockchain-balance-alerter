# Nock Balance Monitor

A Go program that monitors Nockblocks blockchain addresses, converts balances from nick to $NOCK (1 $NOCK = 65,536 nick), and sends notifications to Slack and/or Telegram. It checks balances every minute, alerts on changes, and sends summaries every 6 hours. Balances are stored in `balances.json`.

## Features
- Queries balances via `https://nockblocks.com/rpc`.
- Sends formatted alerts to Slack (block kit) and/or Telegram (MarkdownV2).
- Converts balances: 1 $NOCK = 2^16 nick.
- Supports multiple addresses.
- Stores balances locally.

## Prerequisites
- Go 1.16+
- Slack workspace and/or Telegram account
- Dependencies: `go-co-op/gocron`, `joho/godotenv`, `slack-go/slack`

## Setup

1. **Install Dependencies**:
   ```bash
   go get github.com/go-co-op/gocron github.com/joho/godotenv github.com/slack-go/slack
   ```

2. **Configure Notifications** (at least one required):
  For detailed guide on how to setup slack/telegram bots, refer to [GUIDES](./GUIDES.md)
   - **Slack**:
     - Create an app at [api.slack.com/apps](https://api.slack.com/apps).
     - Add `chat:write` scope, install to workspace, get `xoxb-` token.
     - Add bot to channel: `/invite @BotName`.
   - **Telegram**:
     - Create bot via `@BotFather` in Telegram, get token.
     - Add bot to a group, get chat ID with `@GetIDsBot`.
     - Optionally disable privacy mode: `/setprivacy` > "Disable".

3. **Create `.env`**:
   ```env
   SLACK_BOT_TOKEN=xoxb-your-bot-token
   SLACK_CHANNEL=#channel
   TELEGRAM_BOT_TOKEN=your-telegram-bot-token
   TELEGRAM_CHAT_ID=your-chat-id
   ADDRESSES=3L1PmyRwjyZQ5EQcn4iXECB4v7pyLNAnaU5JCex7NzcJNbFpd3hz5znMYVA33QAHrVc72XeTi62GHqLJqQoJ5w3e4dDDrEQSW7ShSnAvhA7p9RLKXXh2fi7WbKJWJzgmAUMw
   ```
   - Provide at least Slack or Telegram credentials.
   - Add multiple addresses (comma-separated).

4. **Run**:
   ```bash
   go run main.go
   ```

## Example Notification
**Balance Change (Slack/Telegram)**:
```
💸 Balance Change Alert
Address: 3L1P...AUMw
Old: 34492645376 nick (526.18 $NOCK)
New: 34492645377 nick (526.18 $NOCK)
---
Updated: 2025-07-17T15:31:00Z
```

## Troubleshooting
- **No Notifications**:
  - Check `SLACK_BOT_TOKEN` (`xoxb-`), `SLACK_CHANNEL`, `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID`.
  - Ensure bot is in Slack channel or Telegram group.
  - Verify Telegram privacy mode is disabled.
- **Network**: Ensure access to `nockblocks.com`, `slack.com`, `api.telegram.org`.
- **Addresses**: Validate `ADDRESSES` format.

## Security
- Keep tokens secure, regenerate if compromised.
- Add `.env` to `.gitignore`.

## License
MIT License.