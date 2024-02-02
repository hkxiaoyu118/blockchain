package src

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gochain/evm/erc"
	"gochain/utils"
	"golang.org/x/crypto/sha3"
	"math/big"
	"strconv"
	"time"
)

func InitEthNode(rpcUrl string) (*ethclient.Client, *big.Int, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err == nil {
		chainId, err := client.NetworkID(context.Background())
		if err == nil {
			return client, chainId, nil
		}
	}
	return nil, nil, err
}

func TestETHNode(client *ethclient.Client) (string, error) {
	var response string
	beginTime := utils.GetUnixTimeMs()
	header, err := client.HeaderByNumber(context.Background(), nil)
	endTime := utils.GetUnixTimeMs()

	if err == nil {
		response = fmt.Sprintf("节点测试成功! 当前最新区块:%s 节点延迟:%d ms", header.Number.String(), endTime-beginTime)
	}
	return response, err
}

func GetNonce(client *ethclient.Client, address common.Address) (uint64, error) {
	return client.PendingNonceAt(context.Background(), address)
}

func GetTokenBalance(client *ethclient.Client, walletAddress, tokenAddress common.Address) (*big.Int, error) {
	instance, err := erc.NewToken(tokenAddress, client)
	if err == nil {
		bal, err := instance.BalanceOf(&bind.CallOpts{}, walletAddress)
		if err == nil {
			return bal, nil
		}
	}
	return nil, err
}

func GetTokenDecimals(client *ethclient.Client, tokenAddress common.Address) (uint8, error) {
	instance, err := erc.NewToken(tokenAddress, client)
	if err != nil {
		return 0, err
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return 0, err
	} else {
		return decimals, nil
	}
}

func GetEthBalance(client *ethclient.Client, walletAddress common.Address) (*big.Int, error) {
	return client.BalanceAt(context.Background(), walletAddress, nil)
}

func SuggestGasPrice(client *ethclient.Client) (*big.Int, error) {
	return client.SuggestGasPrice(context.Background())
}

func SuggestGasTipCap(client *ethclient.Client) (*big.Int, error) {
	return client.SuggestGasTipCap(context.Background())
}

func EstimateGas(client *ethclient.Client, msg ethereum.CallMsg) (uint64, error) {
	return client.EstimateGas(context.Background(), msg)
}

func GetLatestBlock(client *ethclient.Client) (uint64, uint64, error) {
	var err error
	blockNum, err := client.BlockNumber(context.Background())
	if err == nil {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
		if err == nil {
			timestamp := block.Time()
			return blockNum, timestamp, nil
		}
	}
	return 0, 0, err
}

func QueryTransactions(client *ethclient.Client, txId string) bool {
	txHash := common.HexToHash(txId)
	for i := 0; i < 30; i++ {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			if receipt.Status == 1 {
				return true
			} else {
				return false
			}
		}
		time.Sleep(time.Second)
	}
	return false
}

func Approve(client *ethclient.Client, keypair *KeyPair, chainId *big.Int, tokenAddress, spenderAddress, gasPrice, gasLimit, gasFeeCap, gasTipCap string, eip1559 bool) (string, error) {
	var txID string
	approveFnSignature := []byte("approve(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(approveFnSignature)
	methodID := hash.Sum(nil)[:4]

	spender := common.HexToAddress(spenderAddress)
	paddedSpender := common.LeftPadBytes(spender.Bytes(), 32)
	paddedNumber := common.LeftPadBytes(common.Hex2Bytes("ffffffffffffffffffffffffffffffffffffffff"), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedSpender...)
	data = append(data, paddedNumber...)

	nonce, _ := GetNonce(client, keypair.Address)
	gasLimitBigInt, _ := strconv.ParseUint(gasLimit, 10, 64)

	if eip1559 == true {
		txId, err := SendRawTransactionEIP1559(
			client,
			keypair,
			chainId,
			nonce,
			gasLimitBigInt,
			big.NewInt(0),
			ToUnit(gasFeeCap, 9),
			ToUnit(gasTipCap, 9),
			tokenAddress,
			hex.Dump(data),
		)
		if err == nil {
			txID = txId
		}
		return txID, err
	} else {
		txId, err := SendRawTransaction(
			client,
			keypair,
			chainId,
			nonce,
			gasLimitBigInt,
			big.NewInt(0),
			ToUnit(gasPrice, 9),
			tokenAddress,
			hex.Dump(data),
		)
		if err == nil {
			txID = txId
		}
		return txID, err
	}
}

func TransferToken(client *ethclient.Client, keypair *KeyPair, chainId *big.Int, token, to, gasPrice, gasLimit, gasFeeCap, gasTipCap string, amount *big.Int, eip1559 bool) (string, error) {
	var txID string
	approveFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(approveFnSignature)
	methodID := hash.Sum(nil)[:4]

	toAddress := common.HexToAddress(to)
	paddedSpender := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedNumber := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedSpender...)
	data = append(data, paddedNumber...)

	nonce, _ := GetNonce(client, keypair.Address)
	gasLimitBigInt, _ := strconv.ParseUint(gasLimit, 10, 64)

	if eip1559 == true {
		txId, err := SendRawTransactionEIP1559(
			client,
			keypair,
			chainId,
			nonce,
			gasLimitBigInt,
			big.NewInt(0),
			ToUnit(gasFeeCap, 9),
			ToUnit(gasTipCap, 9),
			token,
			hex.Dump(data),
		)
		if err == nil {
			txID = txId
		}
		return txID, err
	} else {
		gasPriceBigInt := ToUnit(gasPrice, 9)
		txId, err := SendRawTransaction(
			client,
			keypair,
			chainId,
			nonce,
			gasLimitBigInt,
			big.NewInt(0),
			gasPriceBigInt,
			token,
			hex.Dump(data),
		)
		if err == nil {
			txID = txId
		}
		return txID, err
	}
}
