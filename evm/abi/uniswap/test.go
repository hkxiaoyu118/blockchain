package uniswap

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

func Test3() {
	// example https://etherscan.io/tx/0x34654094832b38e6ece5f06b9adcd9870572909b98c54524f9227553097052ad
	path, fees, err := DecodeV3Path("2sF/lY0u5SOiIGIGmUWXwT2DHscAAGSguGmRxiGLNsHRnUounrDONgbrSAAnEF9kqxVE0ocy8KJPRxPCyOwNoInw")
	if err == nil {
		for _, value := range path {
			fmt.Println(value.Hex())
		}
		for _, item := range fees {
			fmt.Println(item)
		}

		data, err := EncodeV3Path(path, fees)
		if err == nil {
			fmt.Println(data)
		} else {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
}

func Test() {
	params := Params{}
	params.Deadline = big.NewInt(1731231231)

	input1 := SwapRouterExactInputParamsV3{}
	input1.Recipient = common.HexToAddress("0x1")
	input1.Payer = common.HexToAddress("0x1")
	input1.AmountIn = big.NewInt(10000)
	input1.AmountOutMinimum = big.NewInt(1)
	command1, _ := GetV3SwapExactInCommand(&input1)
	AppendCommand(&params, command1)

	input12 := SwapRouterExactOutputParamsV3{}
	input12.Recipient = common.HexToAddress("0x1")
	input12.Payer = common.HexToAddress("0x")
	input12.AmountOut = big.NewInt(23333)
	input12.AmountInMaximum = big.NewInt(123123)
	command2, _ := GetV3SwapExactOutCommand(&input12)
	AppendCommand(&params, command2)

	input3 := SwapRouterExactInputParamsV2{}
	input3.Recipient = common.HexToAddress("0x1")
	input3.Payer = common.HexToAddress("0x1")
	input3.AmountIn = big.NewInt(10000)
	input3.AmountOutMinimum = big.NewInt(1)
	command3, _ := GetV2SwapExactInCommand(&input3)
	AppendCommand(&params, command3)

	input14 := SwapRouterExactOutputParamsV2{}
	input14.Recipient = common.HexToAddress("0x1")
	input14.Payer = common.HexToAddress("0x1")
	input14.AmountOut = big.NewInt(23333)
	input14.AmountInMaximum = big.NewInt(123123)
	command4, _ := GetV2SwapExactOutCommand(&input14)
	AppendCommand(&params, command4)

	dataRes, err := Encode(&params)
	if err == nil {
		dataResStr := "0x" + hex.EncodeToString(dataRes)
		fmt.Println(dataResStr)
		Test2(dataResStr)
	} else {
		fmt.Println(err.Error())
	}
}

func Test2(input string) {
	decoded, _ := hexutil.Decode(input)
	cs, _ := Decode(decoded)
	for _, cc := range cs {
		switch RealCommand(cc.Command) {
		case V3_SWAP_EXACT_IN:
			var params SwapRouterExactInputParamsV3
			_ = UnpackV3SwapExactIn(&params, cc.Input)
			b, _ := json.MarshalIndent(params, "", "  ")
			fmt.Printf("%s", string(b))
		case V3_SWAP_EXACT_OUT:
			var params SwapRouterExactOutputParamsV3
			_ = UnpackV3SwapExactOut(&params, cc.Input)
			b, _ := json.MarshalIndent(params, "", "  ")
			fmt.Printf("%s", string(b))
		case V2_SWAP_EXACT_IN:
			var params SwapRouterExactInputParamsV2
			_ = UnpackV2SwapExactIn(&params, cc.Input)
			b, _ := json.MarshalIndent(params, "", "  ")
			fmt.Printf("%s", string(b))
		case V2_SWAP_EXACT_OUT:
			var params SwapRouterExactOutputParamsV2
			_ = UnpackV2SwapExactOut(&params, cc.Input)
			b, _ := json.MarshalIndent(params, "", "  ")
			fmt.Printf("%s", string(b))
		default:
			return
		}
	}
}
