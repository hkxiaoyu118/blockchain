package raydium

import (
	"context"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

func Slippage(client *rpc.Client, tokenIn string, tokenOut string, amountIn uint64, slippage uint64) (uint64, error) {
	inRes, err := client.GetAccountInfo(
		context.TODO(),
		solana.MustPublicKeyFromBase58(tokenIn),
	)
	if err != nil {
		return 0, err
	}

	outRes, err := client.GetAccountInfo(
		context.TODO(),
		solana.MustPublicKeyFromBase58(tokenOut),
	)
	if err != nil {
		return 0, err
	}

	var poolCoinBalance token.Account
	err = bin.NewBinDecoder(inRes.Value.Data.GetBinary()).Decode(&poolCoinBalance)
	if err != nil {
		return 0, err
	}

	var poolPcBalance token.Account
	err = bin.NewBinDecoder(outRes.Value.Data.GetBinary()).Decode(&poolPcBalance)
	if err != nil {
		return 0, err
	}

	denominator := poolCoinBalance.Amount + amountIn
	minimumOutAmount := poolPcBalance.Amount * amountIn / denominator
	// slippage
	minimumOutAmount = minimumOutAmount * (100 - slippage) / 100
	return minimumOutAmount, nil
}
