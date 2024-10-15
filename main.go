package main

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

type AddressActivity struct {
	Address string `json:"address"`
	Score   int    `json:"score"`
}

// for concurrent map access
var (
	activityMap = make(map[string]int)
	mu          sync.Mutex
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		log.Fatal("ETH_RPC_URL not set in the environment")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	http.HandleFunc("GET /top", func(w http.ResponseWriter, r *http.Request) {
		topActiveAddresses(client, w)
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func topActiveAddresses(client *ethclient.Client, w http.ResponseWriter) {
	// Get the latest block
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		http.Error(w, "Failed to fetch the latest block", http.StatusInternalServerError)
		return
	}

	startBlock := blockNumber - 100 // Last 100 blocks

	// Use a WaitGroup to process blocks concurrently
	var wg sync.WaitGroup

	// Limit the number of concurrent workers
	blockChan := make(chan *big.Int, 6) // more that 6 can cause rate limit exceeding

	for i := 0; i < 2; i++ { // more that 2 can cause rate limit exceeding
		wg.Add(1)
		go func() {
			defer wg.Done()
			for blockNum := range blockChan {
				processBlock(client, blockNum)
			}
		}()
	}

	// Send block numbers to the channel
	for i := startBlock; i <= blockNumber; i++ {
		blockChan <- big.NewInt(int64(i))
	}

	close(blockChan)
	wg.Wait()

	// Convert map to a slice and sort by score
	var activities []AddressActivity
	for address, score := range activityMap {
		activities = append(activities, AddressActivity{Address: address, Score: score})
	}

	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Score > activities[j].Score
	})

	// Return top 5 addresses
	if len(activities) > 5 {
		activities = activities[:5]
	}

	json.NewEncoder(w).Encode(activities)
}

// processes token transfers in a block and updates activity scores
func processBlock(client *ethclient.Client, blockNumber *big.Int) {
	query := ethereum.FilterQuery{
		FromBlock: blockNumber,
		ToBlock:   blockNumber,
		//Keccak-256 hash for the ERC-20 Transfer event
		Topics: [][]common.Hash{{common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")}},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Printf("Failed to fetch logs: %v", err)
		return
	}

	for _, vLog := range logs {
		sender := common.HexToAddress(vLog.Topics[1].Hex()).Hex()
		receiver := common.HexToAddress(vLog.Topics[2].Hex()).Hex()

		mu.Lock()
		activityMap[sender]++
		activityMap[receiver]++
		mu.Unlock()
	}
}
