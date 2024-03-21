package src

import (
	"github.com/gagliardetto/solana-go"
)

type KeyPair struct {
	PriKey *solana.PrivateKey
	PubKey *solana.PublicKey
}

func GetKeyPair(secretKey string) (*KeyPair, error) {
	privateKey, err := solana.PrivateKeyFromBase58(secretKey)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.PublicKey()
	keyPair := KeyPair{}
	keyPair.PriKey = &privateKey
	keyPair.PubKey = &publicKey
	return &keyPair, nil
}
