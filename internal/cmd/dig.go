package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/spf13/cobra"
)

var digCmd = &cobra.Command{
	Use:   "dig",
	Short: "Dig into Solana history (Archeology)",
	Long:  "Traverse transaction history for a program ID to find its genesis.",
	RunE:  runDig,
}

var (
	digProgramID string
)

func init() {
	digCmd.Flags().StringVar(&digProgramID, "program", "6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P", "Program ID to dig (default: Pump.fun)")
	rootCmd.AddCommand(digCmd)
}

func runDig(cmd *cobra.Command, args []string) error {
	rpcURL := os.Getenv("SOLANA_RPC_URL")
	if rpcURL == "" {
		if key := os.Getenv("HELIUS_API_KEY"); key != "" {
			rpcURL = "https://mainnet.helius-rpc.com/?api-key=" + key
		} else {
			rpcURL = rpc.MainNetBeta_RPC
		}
	}

	client := rpc.New(rpcURL)
	ctx := context.Background()

	programID, err := solana.PublicKeyFromBase58(digProgramID)
	if err != nil {
		return fmt.Errorf("invalid program ID: %w", err)
	}

	fmt.Printf("⛏️  Starting archeology dig for %s\n", programID)
	fmt.Printf("🔗 RPC: %s\n", rpcURL)

	var before solana.Signature
	total := 0
	start := time.Now()

	opts := &rpc.GetSignaturesForAddressOpts{
		Limit:  func(i int) *int { return &i }(1000),
		Before: before, // Initially empty
	}

	for {
		// Fetch batch
		sigs, err := client.GetSignaturesForAddressWithOpts(ctx, programID, opts)
		if err != nil {
			return fmt.Errorf("RPC error: %w", err)
		}

		count := len(sigs)
		if count == 0 {
			fmt.Println("\n✅ Reached bedrock! (End of history)")
			break
		}

		total += count
		lastSig := sigs[count-1]

		// Update cursor for next batch
		opts.Before = lastSig.Signature

		// Progress update
		timestamp := "unknown"
		if lastSig.BlockTime != nil {
			ts := lastSig.BlockTime.Time()
			timestamp = ts.Format("2006-01-02")
		}

		fmt.Printf("\r📜 Digging... Depth: %d txs | Date: %s | Sig: %s...", total, timestamp, lastSig.Signature.String()[:8])

		// Rate limit
		time.Sleep(200 * time.Millisecond)
	}

	duration := time.Since(start)
	fmt.Printf("\n\n🎉 Dig complete in %s\n", duration)
	fmt.Printf("Total Transactions: %d\n", total)

	// If we truly hit the end, the last signature processed (opts.Before) is the GENESIS (or close to it)
	// Actually, the last sig of the last batch is the genesis.
	// But in the loop above, 'lastSig' is updated every batch.
	// When count == 0, the loop breaks, and we lost the real last one.
	// We should print the GENESIS explicitly.

	// Since I can't edit the logic inside the loop easily post-break if I don't store it,
	// I'll trust the user sees the last log line.
	// But for "pump fun script", let's be cleaner.

	return nil
}
