package evm

import (
	"blockchain/evm/src"
	"math/big"
)

type Account struct {
	KeyPair  *src.KeyPair
	Provider *src.Provider
}

func (m *Account) Init(provider *src.Provider, secretKey string) error {
	m.Provider = provider
	return m.initKeypair(secretKey)
}

func (m *Account) initKeypair(secretKey string) error {
	keyPair, err := src.GetKeyPair(secretKey)
	if err == nil {
		m.KeyPair = keyPair
	}
	return err
}

func (m *Account) GetETHBalance() (*big.Int, error) {
	return src.GetEthBalance(m.Provider.Client, m.KeyPair.Address)
}

func (m *Account) GetNonce() (uint64, error) {
	return src.GetNonce(m.Provider.Client, m.KeyPair.Address)
}

func (m *Account) Approve(tokenAddress, spenderAddress, gasPrice, gasLimit, gasFeeCap, gasTipCap string, eip1559 bool) (string, error) {
	return src.Approve(
		m.Provider.Client,
		m.KeyPair,
		m.Provider.ChainId,
		tokenAddress,
		spenderAddress,
		gasPrice,
		gasLimit,
		gasFeeCap,
		gasTipCap,
		eip1559,
	)
}

func (m *Account) TransferToken(tokenAddress, toAddress, gasPrice, gasFeeCap, gasTipCap, gasLimit string, amount *big.Int, eip1559 bool) (string, error) {
	return src.TransferToken(
		m.Provider.Client,
		m.KeyPair,
		m.Provider.ChainId,
		tokenAddress,
		toAddress,
		gasPrice,
		gasLimit,
		gasFeeCap,
		gasTipCap,
		amount,
		eip1559,
	)
}

func (m *Account) SendRawTransaction(value, to, gasPrice, gasFeeCap, gasTipCap, data string, nonce, gasLimit uint64, eip1559 bool) (string, error) {
	if eip1559 == true {
		return src.SendRawTransactionEIP1559(
			m.Provider.Client,
			m.KeyPair,
			m.Provider.ChainId,
			nonce,
			gasLimit,
			src.ToUnit(value, 18),
			src.ToUnit(gasFeeCap, 9),
			src.ToUnit(gasTipCap, 9),
			to,
			data,
		)
	} else {
		return src.SendRawTransaction(
			m.Provider.Client,
			m.KeyPair,
			m.Provider.ChainId,
			nonce,
			gasLimit,
			src.ToUnit(value, 18),
			src.ToUnit(gasPrice, 9),
			to,
			data,
		)
	}
}
