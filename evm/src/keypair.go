package src

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type KeyPair struct {
	PriKey  *ecdsa.PrivateKey
	Address common.Address
}

func GetKeyPair(secretKey string) (*KeyPair, error) {
	privateKey, err := crypto.HexToECDSA(secretKey)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("生成私钥地址对发生错误:%s", err.Error()))
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New(fmt.Sprintf("生成私钥地址对发生错误:%s", err.Error()))
	}
	walletAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &KeyPair{
		PriKey:  privateKey,
		Address: walletAddress,
	}, nil
}
