package src

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type Provider struct {
	Client  *ethclient.Client
	ChainId *big.Int
}

func (m *Provider) InitProvider(rpcUrl string) error {
	client, chainId, err := InitEthNode(rpcUrl)
	if err == nil {
		m.Client = client
		m.ChainId = chainId
		return nil
	}
	return err
}

func (m *Provider) TestETHNode() (string, error) {
	return TestETHNode(m.Client)
}

func (m *Provider) SuggestGasPrice() (*big.Int, error) {
	return SuggestGasPrice(m.Client)
}

func (m *Provider) SuggestGasTipCap() (*big.Int, error) {
	return SuggestGasTipCap(m.Client)
}

func (m *Provider) EstimateGas(msg ethereum.CallMsg) (uint64, error) {
	return EstimateGas(m.Client, msg)
}

func (m *Provider) GetBlockNum() (uint64, uint64, error) {
	return GetLatestBlock(m.Client)
}

func (m *Provider) GetTokenBalance(Address, tokenAddress common.Address) (*big.Int, error) {
	return GetTokenBalance(m.Client, Address, tokenAddress)
}

func (m *Provider) GetTokenDecimals(tokenAddress common.Address) (uint8, error) {
	return GetTokenDecimals(m.Client, tokenAddress)
}
