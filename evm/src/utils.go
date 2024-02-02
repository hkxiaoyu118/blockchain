package src

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math"
	"math/big"
	"reflect"
)

func IsValidAddress(address string) bool {
	return common.IsHexAddress(address)
}

func IsZeroAddress(iAddress interface{}) bool {
	var address common.Address
	switch v := iAddress.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	zeroAddressBytes := common.FromHex("0x0000000000000000000000000000000000000000")
	addressBytes := address.Bytes()
	return reflect.DeepEqual(addressBytes, zeroAddressBytes)
}

func FromUnit(iValue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := iValue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

func ToUnit(iAmount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iAmount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

func IntWithDecimal(v uint64, decimal int) *big.Int {
	pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	return new(big.Int).Mul(big.NewInt(int64(v)), pow)
}

func IntDivDecimal(v *big.Int, decimal int) *big.Int {
	pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	return new(big.Int).Div(v, pow)
}

func FloatStringToBigInt(amount string, decimals int) *big.Int {
	fAmount, _ := new(big.Float).SetString(amount)
	fi, _ := new(big.Float).Mul(fAmount, big.NewFloat(math.Pow10(decimals))).Int(nil)
	return fi
}
