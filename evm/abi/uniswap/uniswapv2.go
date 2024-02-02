package uniswap

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
)

var swapRouterMetaDataV2 = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"payer\",\"type\":\"address\"}],\"name\":\"v2SwapExactInput\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMaximum\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"payer\",\"type\":\"address\"}],\"name\":\"v2SwapExactOutput\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

type SwapRouterExactInputParamsV2 struct {
	Recipient        common.Address
	AmountIn         *big.Int
	AmountOutMinimum *big.Int
	Path             []common.Address
	Payer            common.Address
}

type SwapRouterExactOutputParamsV2 struct {
	Recipient       common.Address
	AmountOut       *big.Int
	AmountInMaximum *big.Int
	Path            []common.Address
	Payer           common.Address
}

func UnpackV2SwapExactIn(params *SwapRouterExactInputParamsV2, input []byte) error {
	a, err := abi.JSON(strings.NewReader(swapRouterMetaDataV2.ABI))
	if err != nil {
		return err
	}
	method, ok := a.Methods["v2SwapExactInput"]
	if !ok {
		return errors.New("no this method")
	}
	argv, err := method.Inputs.Unpack(input)
	if err != nil {
		return err
	}
	err = method.Inputs.Copy(params, argv)
	if err != nil {
		return err
	}
	return nil
}

func UnpackV2SwapExactOut(params *SwapRouterExactOutputParamsV2, input []byte) error {
	a, err := abi.JSON(strings.NewReader(swapRouterMetaDataV2.ABI))
	if err != nil {
		return err
	}
	method, ok := a.Methods["v2SwapExactOutput"]
	if !ok {
		return errors.New("no this method")
	}
	argv, err := method.Inputs.Unpack(input)
	if err != nil {
		return err
	}
	err = method.Inputs.Copy(params, argv)
	if err != nil {
		return err
	}
	return nil
}

func PackV2SwapExactIn(params *SwapRouterExactInputParamsV2) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(swapRouterMetaDataV2.ABI))
	if err != nil {
		return nil, err
	}
	method, ok := a.Methods["v2SwapExactInput"]
	if !ok {
		return nil, errors.New("no this method")
	}
	output, err := method.Inputs.Pack(params.Recipient, params.AmountIn, params.AmountOutMinimum, params.Path, params.Payer)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func PackV2SwapExactOut(params *SwapRouterExactOutputParamsV2) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(swapRouterMetaDataV2.ABI))
	if err != nil {
		return nil, err
	}
	method, ok := a.Methods["v2SwapExactOutput"]
	if !ok {
		return nil, errors.New("no this method")
	}
	output, err := method.Inputs.Pack(params.Recipient, params.AmountOut, params.AmountInMaximum, params.Path, params.Payer)
	if err != nil {
		return nil, err
	}
	return output, nil
}
