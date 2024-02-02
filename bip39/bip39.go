package bip39

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"strconv"
)

type WalletInfo struct {
	Address    string `json:"address"`    //钱包地址
	PrivateKey string `json:"privateKey"` //钱包私钥
}

func MnemonicToKeys(mnemonic, password string, num int, showTag int) string {
	var privateKeys string
	seed := bip39.NewSeed(mnemonic, password)
	wallet, _ := hdwallet.NewFromSeed(seed)
	for i := 0; i < num; i++ {
		path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", i))
		account, _ := wallet.Derive(path, false)
		privateKey, _ := wallet.PrivateKeyHex(account)
		walletAddress := account.Address.Hex()
		keyItem := ""
		if showTag == 0 {
			keyItem = privateKey
		} else if showTag == 1 {
			keyItem = privateKey + "\t" + walletAddress + "\t账号:" + strconv.Itoa(i)
		}
		privateKeys = privateKeys + keyItem + "\n"
	}
	return privateKeys
}

func GenRandMnemonic(wordsNum int) string {
	bitSize := 128
	if wordsNum == 12 {
		bitSize = 128
	} else if wordsNum == 24 {
		bitSize = 256
	}
	entropy, _ := bip39.NewEntropy(bitSize)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return mnemonic
}

// 通过钱包的私钥获取钱包的信息
func GetWalletAddressByPriKey(privateKeyStr string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("私钥生成公钥失败")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address, nil
}
