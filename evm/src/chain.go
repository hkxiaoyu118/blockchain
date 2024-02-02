package src

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func IsSupportEIP1559(client *ethclient.Client) (bool, error) {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return false, err
	}
	return header.BaseFee != nil, nil
}

func GetHeaderByHash(client *ethclient.Client, hash string) (*types.Header, error) {
	return client.HeaderByHash(context.Background(), common.HexToHash(hash))
}

func GetHeaderByNumber(client *ethclient.Client, number *big.Int) (*types.Header, error) {
	return client.HeaderByNumber(context.Background(), number)
}

func GetBlockByHash(client *ethclient.Client, hash string) (*types.Block, error) {
	return client.BlockByHash(context.Background(), common.HexToHash(hash))
}

func GetBlockByNumber(client *ethclient.Client, number *big.Int) (*types.Block, error) {
	return client.BlockByNumber(context.Background(), number)
}

func GetBlockTxCount(client *ethclient.Client, hash string) (uint, error) {
	return client.TransactionCount(context.Background(), common.HexToHash(hash))
}
