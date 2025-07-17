# Nock Balance Monitor: Configuration Guides

This document provides detailed instructions for configuring Slack and Telegram bots for the Nock Balance Monitor, which sends balance change alerts and summaries for Nockblocks blockchain addresses. Follow these steps to set up notifications for Slack, Telegram, or both, as supported by the program.

## Prerequisites
- A Slack workspace where you have permission to create apps and add bots to channels.
- A Telegram account and access to a group where you can add bots.
- The Nock Balance Monitor program (`main.go`) and `.env` file set up as per the [README](README.md).

## Slack Bot Configuration

### Step 1: Create a Slack App
1. **Visit the Slack API Portal**:
   - Go to [api.slack.com/apps](https://api.slack.com/apps).
   - Click **Create New App** > **From scratch**.
2. **Name the App**:
   - Enter a name, e.g., `NockBalanceBot`.
   - Select your Slack workspace from the dropdown.
   - Click **Create App**.
3. **Configure Bot Permissions**:
   - Navigate to **OAuth & Permissions** in the left sidebar.
   - Under **Bot Token Scopes**, click **Add an OAuth Scope**.
   - Add the `chat:write` scope to allow the bot to send messages.
   - Optionally, add `channels:read` if you need to list or verify channels (not required for basic functionality).
4. **Install the App**:
   - Go to **OAuth & Permissions** > **Install to Workspace**.
   - Authorize the app for your workspace.
   - After installation, copy the **Bot User OAuth Token** (starts with `xoxb-`, e.g., `xoxb-1234567890-abcdefg`).
   - Save this token securely for the `.env` file.

### Step 2: Add the Bot to a Slack Channel
1. **Choose or Create a Channel**:
   - In Slack, create a new channel (e.g., `#nock-balances`) or select an existing one.
   - Ensure you have permission to add apps to the channel.
2. **Invite the Bot**:
   - In the channel, type `/invite @NockBalanceBot` (replace with your bot’s name) and press Enter.
   - Alternatively:
     - Click the channel name, select **Integrations** > **Add an App**.
     - Search for your bot (e.g., `NockBalanceBot`) and add it.
3. **Verify Permissions**:
   - Ensure the bot has permission to post messages (granted by the `chat:write` scope).
   - If the channel is private, confirm the bot is added as a member.
4. **Note the Channel Name**:
   - Use the channel name with a `#` prefix (e.g., `#nock-balances`) for the `SLACK_CHANNEL` variable in the `.env` file.

### Step 3: Update `.env`
Add the Slack credentials to your `.env` file:
```env
SLACK_BOT_TOKEN=xoxb-your-bot-token
SLACK_CHANNEL=#nock-balances
```

### Step 4: Test the Bot
- Run the program (`go run main.go`).
- Verify that balance change alerts and summaries appear in the Slack channel with formatted blocks (e.g., headers, dividers, and emojis like 💸 and 📊).

## Telegram Bot Configuration

### Step 1: Create a Telegram Bot
1. **Open Telegram**:
   - Use the Telegram app or web version ([web.telegram.org](https://web.telegram.org)).
2. **Start BotFather**:
   - Search for `@BotFather` (Telegram’s official bot for creating bots).
   - Send `/start` to begin.
3. **Create a New Bot**:
   - Send `/newbot`.
   - Provide a **name** (e.g., `NockBalanceBot`) and **username** (e.g., `@NockBalanceBot`, must end with `Bot` or `bot`).
   - Copy the **bot token** provided by `@BotFather` (e.g., `123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11`).
   - Save the token securely for the `.env` file.

### Step 2: Add the Bot to a Group and Get Chat ID
1. **Create or Choose a Group**:
   - Create a new group in Telegram or select an existing one where you have admin privileges.
2. **Add the Bot**:
   - Open the group, click the group name, and select **Add Members**.
   - Search for your bot’s username (e.g., `@NockBalanceBot`) and add it.
3. **Grant Permissions**:
   - Make the bot an admin to ensure it can send messages:
     - Go to group settings, select **Administrators**, and add the bot.
     - Enable permissions like **Send Messages**.
4. **Get the Chat ID**:
   - Send a message in the group (e.g., `/start` or any text).
   - Use `@GetIDsBot` in the group to get the chat ID (e.g., `-123456789`).
   - Alternatively, send a message to the bot directly and check its response or logs (requires code modification to capture `chat_id`).
5. **Disable Privacy Mode (Optional)**:
   - By default, Telegram bots in groups only see messages mentioning them (e.g., `@NockBalanceBot`).
   - To allow the bot to see all messages (if needed):
     - In `@BotFather`, send `/setprivacy`, select your bot, and choose **Disable**.

### Step 3: Update `.env`
Add the Telegram credentials to your `.env` file:
```env
TELEGRAM_BOT_TOKEN=your-telegram-bot-token
TELEGRAM_CHAT_ID=your-telegram-chat-id
```

### Step 4: Test the Bot
- Run the program (`go run main.go`).
- Verify that balance change alerts and summaries appear in the Telegram group with MarkdownV2 formatting (e.g., bold text, code blocks, and emojis like 💸 and 📊).

## Combined Configuration
- **Dual Notifications**: The program sends notifications to both Slack and Telegram if both are configured in `.env`. At least one must be set up.
- **Example `.env`**:
  ```env
  SLACK_BOT_TOKEN=xoxb-your-bot-token
  SLACK_CHANNEL=#nock-balances
  TELEGRAM_BOT_TOKEN=your-telegram-bot-token
  TELEGRAM_CHAT_ID=your-telegram-chat-id
  ADDRESSES=3L1PmyRwjyZQ5EQcn4iXECB4v7pyLNAnaU5JCex7NzcJNbFpd3hz5znMYVA33QAHrVc72XeTi62GHqLJqQoJ5w3e4dDDrEQSW7ShSnAvhA7p9RLKXXh2fi7WbKJWJzgmAUMw
  ```
- **Run**:
  ```bash
  go run main.go
  ```

## Troubleshooting
- **Slack Issues**:
  - **No Messages**: Verify `SLACK_BOT_TOKEN` starts with `xoxb-` and has `chat:write`. Ensure bot is in the channel (`/invite @NockBalanceBot`).
  - **Wrong Channel**: Confirm `SLACK_CHANNEL` matches the channel name (e.g., `#nock-balances`).
- **Telegram Issues**:
  - **No Messages**: Check `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID`. Ensure bot is in the group and privacy mode is disabled (`/setprivacy` > "Disable").
  - **Formatting Issues**: Verify special characters (e.g., `_`) are escaped in messages.
- **Network**: Ensure access to `https://nockblocks.com/rpc`, `https://slack.com/api`, and `https://api.telegram.org`.
- **Logs**: Check console logs for errors (e.g., `Error sending Slack message` or `Error sending Telegram message`).

## Security Notes
- Keep `SLACK_BOT_TOKEN` and `TELEGRAM_BOT_TOKEN` secure.
- Regenerate tokens if compromised:
  - Slack: In [api.slack.com/apps](https://api.slack.com/apps), go to **OAuth & Permissions** > **Regenerate Token**.
  - Telegram: In `@BotFather`, send `/token` to regenerate.
- Add `.env` to `.gitignore`.