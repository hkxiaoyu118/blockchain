package src

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"
)

func GenMethodId(method string) string {
	transferFnSignature := []byte(method)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	return hexutil.Encode(methodID)
}

func GenPaddingString(str string) string {
	d, err := hex.DecodeString(str)
	if err == nil {
		paddedAmount := common.LeftPadBytes(d, 32)
		return hexutil.Encode(paddedAmount)
	}
	return ""
}
