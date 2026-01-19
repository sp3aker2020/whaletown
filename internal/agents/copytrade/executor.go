package copytrade

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/speaker20/whaletown/internal/agents/common"
)

// Executor handles trade execution via Jupiter.
type Executor struct {
	config     *common.Config
	privateKey solana.PrivateKey
	rpcClient  *rpc.Client
	httpClient *http.Client
}

// NewExecutor creates a new trade executor.
func NewExecutor(config *common.Config) (*Executor, error) {
	if config.SolanaPrivateKey == "" {
		return nil, fmt.Errorf("SOLANA_PRIVATE_KEY is not set")
	}

	// Parse private key (Base58)
	privKey, err := solana.PrivateKeyFromBase58(config.SolanaPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	rpcURL := config.SolanaRPCURL
	if rpcURL == "" {
		rpcURL = rpc.MainNetBeta_RPC
	}

	return &Executor{
		config:     config,
		privateKey: privKey,
		rpcClient:  rpc.New(rpcURL),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// ExecuteCopyBuy executes a buy order for the specified token.
// It uses a fixed amount of SOL (~$1 @ $200/SOL = 0.005 SOL)
func (e *Executor) ExecuteCopyBuy(tokenMint string) (string, error) {
	fmt.Printf("🚀 FAST LANE: Executing buy for %s\n", tokenMint)

	// 0.005 SOL = 5000000 lamports
	amount := 5000000
	quote, err := e.getJupiterQuote("So11111111111111111111111111111111111111112", tokenMint, amount)
	if err != nil {
		return "", fmt.Errorf("jupiter quote failed: %w", err)
	}

	// 2. Get Swap Transaction
	swapTx, err := e.getJupiterSwapTx(quote)
	if err != nil {
		return "", fmt.Errorf("jupiter swap build failed: %w", err)
	}

	// 3. Sign and Send
	sig, err := e.signAndSend(swapTx)
	if err != nil {
		return "", fmt.Errorf("sign/send failed: %w", err)
	}

	return sig.String(), nil
}

// ExecutionResult holds the result of a copy buy execution.
type ExecutionResult struct {
	TokenMint string
	TxHash    string
}

// ProcessSignal analyzes a transaction signature and executes a copy trade if applicable.
// Returns the execution result if successful.
func (e *Executor) ProcessSignal(signature solana.Signature) (*ExecutionResult, error) {
	ctx := context.Background()

	// Fetch transaction
	tx, err := e.rpcClient.GetTransaction(ctx, signature, &rpc.GetTransactionOpts{
		Commitment:                     rpc.CommitmentConfirmed,
		MaxSupportedTransactionVersion: func(v uint64) *uint64 { return &v }(0),
	})
	if err != nil {
		return nil, fmt.Errorf("fetch tx failed: %w", err)
	}

	if tx == nil || tx.Meta == nil {
		return nil, fmt.Errorf("tx meta missing")
	}

	// Analyze balance changes to find what was bought
	for _, balance := range tx.Meta.PostTokenBalances {
		preAmount := int64(0)

		for _, pre := range tx.Meta.PreTokenBalances {
			if pre.AccountIndex == balance.AccountIndex {
				if val, err := strconv.ParseInt(pre.UiTokenAmount.Amount, 10, 64); err == nil {
					preAmount = val
				}
				break
			}
		}

		postAmount, _ := strconv.ParseInt(balance.UiTokenAmount.Amount, 10, 64)

		if postAmount > preAmount {
			// Ignore WSOL
			if balance.Mint.String() == "So11111111111111111111111111111111111111112" {
				continue
			}

			mint := balance.Mint.String()
			fmt.Printf("🎯 Signal Identified: Whale bought %s\n", mint)

			txSig, err := e.ExecuteCopyBuy(mint)
			if err != nil {
				return nil, fmt.Errorf("copy buy execution failed: %w", err)
			}
			fmt.Printf("✅ Copy Trade Executed! Sig: %s\n", txSig)
			return &ExecutionResult{TokenMint: mint, TxHash: txSig}, nil
		}
	}

	return nil, fmt.Errorf("no buy signal detected in tx")
}

// Internal Jupiter helpers

type jupiterQuoteResponse struct {
	Data []interface{} `json:"data"` // Simplified
	// We just need the raw JSON to pass back to swap
	Raw json.RawMessage
}

func (e *Executor) getJupiterQuote(inputMint, outputMint string, amount int) (json.RawMessage, error) {
	url := fmt.Sprintf("https://quote-api.jup.ag/v6/quote?inputMint=%s&outputMint=%s&amount=%d&slippageBps=50",
		inputMint, outputMint, amount)

	resp, err := e.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	// Calculate reading the body efficiently
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(body), nil
}

func (e *Executor) getJupiterSwapTx(quoteResponse json.RawMessage) (string, error) {
	reqBody := map[string]interface{}{
		"quoteResponse":    quoteResponse,
		"userPublicKey":    e.privateKey.PublicKey().String(),
		"wrapAndUnwrapSol": true,
	}

	jsonBody, _ := json.Marshal(reqBody)
	resp, err := e.httpClient.Post("https://quote-api.jup.ag/v6/swap", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		SwapTransaction string `json:"swapTransaction"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.SwapTransaction, nil
}

func (e *Executor) signAndSend(base64Tx string) (solana.Signature, error) {
	// Decode transaction
	txBytes, err := base64.StdEncoding.DecodeString(base64Tx)
	if err != nil {
		return solana.Signature{}, err
	}

	tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(txBytes))
	if err != nil {
		return solana.Signature{}, fmt.Errorf("decoding tx: %w", err)
	}

	// Sign
	// Note: Jupiter transactions are Versioned Transactions usually.
	// solana-go handles them slightly differently.
	// If it's a legacy tx, this works. If versioned, we need specific handling.
	// Jupiter V6 returns Versioned Transactions (base64).

	// Using generic handling for raw bytes might be safer if library supports it,
	// but solana-go requires object model to sign.

	// Let's try to sign with the private key
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(e.privateKey.PublicKey()) {
				return &e.privateKey
			}
			return nil
		},
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("signing error: %w", err)
	}

	// Send
	sig, err := e.rpcClient.SendTransaction(context.Background(), tx)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("rpc send error: %w", err)
	}

	return sig, nil
}
