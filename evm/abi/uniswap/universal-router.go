package uniswap

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hkxiaoyu118/blockchain/utils"
	"math/big"
	"strings"
)

type Params struct {
	Commands []byte
	Inputs   [][]byte
	Deadline *big.Int
}

var UniversalRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"commands\",\"type\":\"bytes\"},{\"internalType\":\"bytes[]\",\"name\":\"inputs\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

func Decode(data []byte) ([]Command, error) {
	urAbi, err := abi.JSON(strings.NewReader(UniversalRouterMetaData.ABI))
	if err != nil {
		return nil, err
	}
	sigData := data[:4]

	method, err := urAbi.MethodById(sigData)
	if err != nil {
		return nil, err
	}
	argv, err := method.Inputs.Unpack(data[4:])
	if err != nil {
		return nil, err
	}
	if len(argv) != 3 {
		return nil, nil
	}
	var params Params
	err = method.Inputs.Copy(&params, argv)
	if err != nil {
		return nil, err
	}
	if len(params.Commands) != len(params.Inputs) {
		return nil, nil
	}
	var commands []Command
	for i := 0; i < len(params.Commands); i++ {
		commands = append(commands, Command{
			Command: params.Commands[i],
			Input:   params.Inputs[i],
		})
	}
	return commands, nil
}

func Encode(params *Params) ([]byte, error) {
	urAbi, err := abi.JSON(strings.NewReader(UniversalRouterMetaData.ABI))
	if err != nil {
		return nil, err
	}
	method, ok := urAbi.Methods["execute"]
	if !ok {
		return nil, errors.New("no this method")
	}
	data, err := method.Inputs.Pack(params.Commands, params.Inputs, params.Deadline)
	if err != nil {
		return nil, err
	}
	return append(method.ID, data...), nil
}

func AppendCommand(params *Params, command *Command) {
	params.Commands = append(params.Commands, command.Command)
	params.Inputs = append(params.Inputs, command.Input)
}

func EncodeV3Path(paths []common.Address, fees []string) (string, error) {
	var pathString string
	if len(paths) == 0 {
		return pathString, errors.New("path length = 0")
	}

	str := ""
	for i := 0; i < len(paths); i++ {
		str += utils.RemovePrefix0x(paths[i].Hex())
		if i < len(fees) {
			str += utils.RemovePrefix0x(fees[i])
		}
	}

	pathString, err := utils.HexToBase64(str)
	if err != nil {
		return pathString, err
	}
	return pathString, nil
}

// DecodeV3Path https://ethereum.stackexchange.com/questions/144478/uniswap-universal-router-decoding-the-execute-function-parameters
func DecodeV3Path(pathString string) ([]common.Address, []string, error) {
	var path []common.Address
	var fees []string

	pathString, err := utils.Base64ToHex(pathString)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < len(pathString); i += 40 {
		address := utils.AddPrefix0x(utils.StrSubString(pathString, i, i+40))
		path = append(path, common.HexToAddress(address))

		if i+46 <= len(pathString) {
			fees = append(fees, utils.StrSubString(pathString, i+40, i+46))
		}
		i = i + 6
	}
	return path, fees, nil
}

func GetV3SwapExactInCommand(params *SwapRouterExactInputParamsV3) (*Command, error) {
	input, err := PackV3SwapExactIn(params)
	if err != nil {
		return nil, err
	}
	command := RealCommand(V3_SWAP_EXACT_IN)
	return &Command{command, input}, nil
}

func GetV3SwapExactOutCommand(params *SwapRouterExactOutputParamsV3) (*Command, error) {
	input, err := PackV3SwapExactOut(params)
	if err != nil {
		return nil, err
	}
	command := RealCommand(V3_SWAP_EXACT_OUT)
	return &Command{command, input}, nil
}

func GetV2SwapExactInCommand(params *SwapRouterExactInputParamsV2) (*Command, error) {
	input, err := PackV2SwapExactIn(params)
	if err != nil {
		return nil, err
	}
	command := RealCommand(V2_SWAP_EXACT_IN)
	return &Command{command, input}, nil
}

func GetV2SwapExactOutCommand(params *SwapRouterExactOutputParamsV2) (*Command, error) {
	input, err := PackV2SwapExactOut(params)
	if err != nil {
		return nil, err
	}
	command := RealCommand(V2_SWAP_EXACT_OUT)
	return &Command{command, input}, nil
}
