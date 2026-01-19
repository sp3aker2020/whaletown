package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/speaker20/whaletown/internal/agents/common"
	"github.com/speaker20/whaletown/internal/agents/copytrade"
)

func main() {
	// Load .env
	godotenv.Load()

	tokenMint := "Beatbd1WM7MfhDk9oHQeBNe1Uii5nKqZskURsZHupump"
	if len(os.Args) > 1 {
		tokenMint = os.Args[1]
	}

	fmt.Printf("🧪 Testing Fast Lane Buy: %s\n", tokenMint)

	config := common.DefaultConfig()

	executor, err := copytrade.NewExecutor(config)
	if err != nil {
		fmt.Printf("❌ Executor init failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🚀 Executing buy...")
	sig, err := executor.ExecuteCopyBuy(tokenMint)
	if err != nil {
		fmt.Printf("❌ Buy failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ SUCCESS! TX: https://solscan.io/tx/%s\n", sig)
}
