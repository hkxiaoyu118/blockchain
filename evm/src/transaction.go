package src

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

// SendRawTransactionEIP1559 EIP1559
func SendRawTransactionEIP1559(client *ethclient.Client, keypair *KeyPair, chainID *big.Int, nonce, gasLimit uint64, value, gasFeeCap, gasTipCap *big.Int, to, dataStr string) (string, error) {
	var txId string
	data, _ := hex.DecodeString(dataStr)
	toAddress := common.HexToAddress(to)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &toAddress,
		Value:     value,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gasLimit,
		Data:      data,
	})

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), keypair.PriKey)
	if err == nil {
		txId = signedTx.Hash().Hex()
		err = client.SendTransaction(context.Background(), signedTx)
	}
	return txId, err
}

// SendRawTransaction 传统交易
func SendRawTransaction(client *ethclient.Client, keypair *KeyPair, chainID *big.Int, nonce, gasLimit uint64, value, gasPrice *big.Int, to, dataStr string) (string, error) {
	var txId string
	data, _ := hex.DecodeString(dataStr)
	toAddress := common.HexToAddress(to)

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), keypair.PriKey)
	if err == nil {
		txId = signedTx.Hash().Hex()
		err = client.SendTransaction(context.Background(), signedTx)
	}
	return txId, err
}
