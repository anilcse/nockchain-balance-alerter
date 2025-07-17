package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

// Config holds the application configuration
type Config struct {
	SlackBotToken    string   `json:"slackBotToken"`
	SlackChannel     string   `json:"slackChannel"`
	TelegramBotToken string   `json:"telegramBotToken"`
	TelegramChatID   string   `json:"telegramChatID"`
	Addresses        []string `json:"addresses"`
}

// BalanceData stores the balance information for an address
type BalanceData struct {
	Address        string `json:"address"`
	CurrentBalance int64  `json:"currentBalance"`
	LastUpdated    int64  `json:"lastUpdated"`
}

// RPCRequest represents the JSON-RPC request structure
type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      string        `json:"id"`
}

// RPCResponse represents the JSON-RPC response structure
type RPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		Address        string `json:"address"`
		CurrentBalance int64  `json:"currentBalance"`
	} `json:"result"`
	ID string `json:"id"`
}

// State holds the current state of balances
type State struct {
	Balances []BalanceData `json:"balances"`
}

const (
	rpcURL          = "https://nockblocks.com/rpc"
	balanceFile     = "balances.json"
	checkInterval   = 1 * time.Minute
	summaryInterval = 6 * time.Hour
	nickPerNock     = 65536 // 2^16 nick per $NOCK
)

// loadConfig loads configuration from environment variables
func loadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables directly")
	}

	config := Config{
		SlackBotToken:    os.Getenv("SLACK_BOT_TOKEN"),
		SlackChannel:     os.Getenv("SLACK_CHANNEL"),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
		Addresses:        []string{},
	}

	addresses := os.Getenv("ADDRESSES")
	if addresses != "" {
		config.Addresses = strings.Split(addresses, ",")
	}

	if (config.SlackBotToken == "" || config.SlackChannel == "") && (config.TelegramBotToken == "" || config.TelegramChatID == "") {
		return config, fmt.Errorf("either SLACK_BOT_TOKEN and SLACK_CHANNEL or TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID must be set")
	}

	return config, nil
}

// loadState loads the previous balances from file
func loadState() (State, error) {
	var state State
	data, err := os.ReadFile(balanceFile)
	if err != nil {
		if os.IsNotExist(err) {
			return State{Balances: []BalanceData{}}, nil
		}
		return state, err
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return state, err
	}
	return state, nil
}

// saveState saves the current balances to file
func saveState(state State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(balanceFile, data, 0644)
}

// getBalance queries the balance for a given address
func getBalance(address string) (int64, error) {
	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  "getTransactionsByAddress",
		Params: []interface{}{
			map[string]interface{}{
				"address": address,
				"limit":   20,
				"offset":  0,
			},
		},
		ID: fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	body, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var rpcResp RPCResponse
	if err := json.Unmarshal(responseBody, &rpcResp); err != nil {
		return 0, err
	}

	return rpcResp.Result.CurrentBalance, nil
}

// convertToNock converts nick to $NOCK
func convertToNock(nick int64) float64 {
	return float64(nick) / float64(nickPerNock)
}

// formatBalance formats the balance in both nick and $NOCK
func formatBalance(nick int64) string {
	nock := convertToNock(nick)
	return fmt.Sprintf("%d nick (%.2f $NOCK)", nick, nock)
}

// sendSlackMessage sends a formatted message to a Slack channel using block kit
func sendSlackMessage(botToken, channel string, blocks []slack.Block) error {
	if botToken == "" || channel == "" {
		return nil // Skip if Slack is not configured
	}
	api := slack.New(botToken)
	_, _, err := api.PostMessage(
		channel,
		slack.MsgOptionBlocks(blocks...),
		slack.MsgOptionAsUser(true),
	)
	return err
}

// sendTelegramMessage sends a formatted message to a Telegram chat
func sendTelegramMessage(botToken, chatID, message string) error {
	if botToken == "" || chatID == "" {
		return nil // Skip if Telegram is not configured
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       message,
		"parse_mode": "MarkdownV2",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// createBalanceChangeBlocks creates Slack blocks for a balance change alert
func createBalanceChangeBlocks(address, oldBalance, newBalance string) []slack.Block {
	return []slack.Block{
		slack.NewHeaderBlock(
			slack.NewTextBlockObject("plain_text", "💸 Balance Change Alert", true, false),
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Address*: `%s`", address), false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Old Balance*: %s", oldBalance), false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*New Balance*: %s", newBalance), false, false),
			nil,
			nil,
		),
		slack.NewDividerBlock(),
		slack.NewContextBlock(
			"",
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("_Updated at %s_", time.Now().Format(time.RFC3339)), false, false),
		),
	}
}

// createSummaryBlocks creates Slack blocks for the balance summary
func createSummaryBlocks(balances []BalanceData) []slack.Block {
	blocks := []slack.Block{
		slack.NewHeaderBlock(
			slack.NewTextBlockObject("plain_text", "📊 Balance Summary", true, false),
		),
	}

	for i, balance := range balances {
		blocks = append(blocks,
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Address %d*: `%s`", i+1, balance.Address), false, false),
				nil,
				nil,
			),
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Balance*: %s", formatBalance(balance.CurrentBalance)), false, false),
				nil,
				nil,
			),
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Last Updated*: %s", time.Unix(balance.LastUpdated, 0).Format(time.RFC3339)), false, false),
				nil,
				nil,
			),
			slack.NewDividerBlock(),
		)
	}

	blocks = append(blocks,
		slack.NewContextBlock(
			"",
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("_Generated at %s_", time.Now().Format(time.RFC3339)), false, false),
		),
	)

	return blocks
}

// createTelegramBalanceChangeMessage creates a Telegram markdown message for a balance change
func createTelegramBalanceChangeMessage(address, oldBalance, newBalance string) string {
	// Escape special characters for Telegram MarkdownV2
	escapedAddress := strings.ReplaceAll(address, "_", "\\_")
	return fmt.Sprintf(
		"💸 *Balance Change Alert*\n\n"+
			"*Address*: `%s`\n"+
			"*Old Balance*: %s\n"+
			"*New Balance*: %s\n"+
			"──────────\n"+
			"_Updated at %s_",
		escapedAddress,
		oldBalance,
		newBalance,
		time.Now().Format(time.RFC3339),
	)
}

// createTelegramSummaryMessage creates a Telegram markdown message for the balance summary
func createTelegramSummaryMessage(balances []BalanceData) string {
	message := "📊 *Balance Summary*\n\n"
	for i, balance := range balances {
		// Escape special characters for Telegram MarkdownV2
		escapedAddress := strings.ReplaceAll(balance.Address, "_", "\\_")
		message += fmt.Sprintf(
			"*Address %d*: `%s`\n"+
				"*Balance*: %s\n"+
				"*Last Updated*: %s\n"+
				"──────────\n",
			i+1,
			escapedAddress,
			formatBalance(balance.CurrentBalance),
			time.Unix(balance.LastUpdated, 0).Format(time.RFC3339),
		)
	}
	message += fmt.Sprintf("_Generated at %s_", time.Now().Format(time.RFC3339))
	return message
}

// checkBalances checks all addresses for balance changes
func checkBalances(config Config, state *State) {
	for _, address := range config.Addresses {
		newBalance, err := getBalance(address)
		if err != nil {
			log.Printf("Error checking balance for %s: %v", address, err)
			continue
		}

		var oldBalance int64
		var balanceIndex = -1
		for i, b := range state.Balances {
			if b.Address == address {
				oldBalance = b.CurrentBalance
				balanceIndex = i
				break
			}
		}

		if balanceIndex == -1 {
			// New address
			state.Balances = append(state.Balances, BalanceData{
				Address:        address,
				CurrentBalance: newBalance,
				LastUpdated:    time.Now().Unix(),
			})
			// Slack notification
			blocks := createBalanceChangeBlocks(
				address,
				"Initial balance",
				formatBalance(newBalance),
			)
			if err := sendSlackMessage(config.SlackBotToken, config.SlackChannel, blocks); err != nil {
				log.Printf("Error sending Slack message: %v", err)
			}
			// Telegram notification
			message := createTelegramBalanceChangeMessage(
				address,
				"Initial balance",
				formatBalance(newBalance),
			)
			if err := sendTelegramMessage(config.TelegramBotToken, config.TelegramChatID, message); err != nil {
				log.Printf("Error sending Telegram message: %v", err)
			}
		} else if newBalance != oldBalance {
			// Balance changed
			state.Balances[balanceIndex].CurrentBalance = newBalance
			state.Balances[balanceIndex].LastUpdated = time.Now().Unix()
			// Slack notification
			blocks := createBalanceChangeBlocks(
				address,
				formatBalance(oldBalance),
				formatBalance(newBalance),
			)
			if err := sendSlackMessage(config.SlackBotToken, config.SlackChannel, blocks); err != nil {
				log.Printf("Error sending Slack message: %v", err)
			}
			// Telegram notification
			message := createTelegramBalanceChangeMessage(
				address,
				formatBalance(oldBalance),
				formatBalance(newBalance),
			)
			if err := sendTelegramMessage(config.TelegramBotToken, config.TelegramChatID, message); err != nil {
				log.Printf("Error sending Telegram message: %v", err)
			}
		}
	}

	if err := saveState(*state); err != nil {
		log.Printf("Error saving state: %v", err)
	}
}

// sendSummary sends a summary of all balances
func sendSummary(config Config, state State) {
	// Slack notification
	blocks := createSummaryBlocks(state.Balances)
	if err := sendSlackMessage(config.SlackBotToken, config.SlackChannel, blocks); err != nil {
		log.Printf("Error sending Slack summary: %v", err)
	}
	// Telegram notification
	message := createTelegramSummaryMessage(state.Balances)
	if err := sendTelegramMessage(config.TelegramBotToken, config.TelegramChatID, message); err != nil {
		log.Printf("Error sending Telegram summary: %v", err)
	}
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	state, err := loadState()
	if err != nil {
		log.Fatalf("Error loading state: %v", err)
	}

	scheduler := gocron.NewScheduler(time.UTC)

	// Schedule balance check every minute
	_, err = scheduler.Every(checkInterval).Do(func() {
		checkBalances(config, &state)
	})
	if err != nil {
		log.Fatalf("Error scheduling balance check: %v", err)
	}

	// Schedule summary every 6 hours
	_, err = scheduler.Every(summaryInterval).Do(func() {
		sendSummary(config, state)
	})
	if err != nil {
		log.Fatalf("Error scheduling summary: %v", err)
	}

	scheduler.StartAsync()
	log.Println("Cron job started. Monitoring addresses...")

	// Keep the program running
	select {}
}
