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
	SlackBotToken string   `json:"slackBotToken"`
	SlackChannel  string   `json:"slackChannel"`
	Addresses     []string `json:"addresses"`
}

// BalanceData stores the balance information for an address
type BalanceData struct {
	Address       string `json:"address"`
	CurrentBalance int64  `json:"currentBalance"`
	LastUpdated   int64  `json:"lastUpdated"`
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
	rpcURL           = "https://nockblocks.com/rpc"
	balanceFile      = "balances.json"
	checkInterval    = 1 * time.Minute
	summaryInterval  = 6 * time.Hour
	nickPerNock      = 65536 // 2^16 nick per $NOCK
)

// loadConfig loads configuration from environment variables
func loadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables directly")
	}

	config := Config{
		SlackBotToken: os.Getenv("SLACK_BOT_TOKEN"),
		SlackChannel:  os.Getenv("SLACK_CHANNEL"),
		Addresses:     []string{},
	}

	addresses := os.Getenv("ADDRESSES")
	if addresses != "" {
		config.Addresses = strings.Split(addresses, ",")
	}

	if config.SlackBotToken == "" || config.SlackChannel == "" {
		return config, fmt.Errorf("SLACK_BOT_TOKEN and SLACK_CHANNEL must be set")
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

// sendSlackMessage sends a message to a Slack channel
func sendSlackMessage(botToken, channel, message string) error {
	api := slack.New(botToken)
	_, _, err := api.PostMessage(
		channel,
		slack.MsgOptionText(message, false),
		slack.MsgOptionAsUser(true),
	)
	return err
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
				Address:       address,
				CurrentBalance: newBalance,
				LastUpdated:   time.Now().Unix(),
			})
			message := fmt.Sprintf("New address monitored: %s\nBalance: %s", address, formatBalance(newBalance))
			if err := sendSlackMessage(config.SlackBotToken, config.SlackChannel, message); err != nil {
				log.Printf("Error sending Slack message: %v", err)
			}
		} else if newBalance != oldBalance {
			// Balance changed
			state.Balances[balanceIndex].CurrentBalance = newBalance
			state.Balances[balanceIndex].LastUpdated = time.Now().Unix()
			message := fmt.Sprintf("Balance change for %s\nOld: %s\nNew: %s", address, formatBalance(oldBalance), formatBalance(newBalance))
			if err := sendSlackMessage(config.SlackBotToken, config.SlackChannel, message); err != nil {
				log.Printf("Error sending Slack message: %v", err)
			}
		}
	}

	if err := saveState(*state); err != nil {
		log.Printf("Error saving state: %v", err)
	}
}

// sendSummary sends a summary of all balances
func sendSummary(config Config, state State) {
	message := "Balance Summary:\n\n"
	for _, balance := range state.Balances {
		message += fmt.Sprintf("Address: %s\nBalance: %s\nLast Updated: %s\n\n",
			balance.Address,
			formatBalance(balance.CurrentBalance),
			time.Unix(balance.LastUpdated, 0).Format(time.RFC3339))
	}

	if err := sendSlackMessage(config.SlackBotToken, config.SlackChannel, message); err != nil {
		log.Printf("Error sending summary: %v", err)
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
