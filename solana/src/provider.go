package src

import (
	"context"
	"github.com/AlekSi/pointer"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"math/big"
)

type Provider struct {
	Client *rpc.Client
}

func (m *Provider) InitProvider(rpcUrl string) error {
	client := rpc.New(rpcUrl)
	_, err := client.GetRecentBlockhash(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
	if err == nil {
		m.Client = client
	}
	return err
}

func (m *Provider) GetHealth() (bool, error) {
	out, err := m.Client.GetHealth(
		context.TODO(),
	)
	if err != nil {
		return false, err
	}
	return out == rpc.HealthOk, err
}

func (m *Provider) GetVersion() (*rpc.GetVersionResult, error) {
	return m.Client.GetVersion(
		context.TODO(),
	)
}

func (m *Provider) GetSlot() (uint64, error) {
	return m.Client.GetSlot(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
}

func (m *Provider) GetSlotLeader() (solana.PublicKey, error) {
	return m.Client.GetSlotLeader(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
}

func (m *Provider) GetRecentBlockHash() (string, error) {
	blockHash := ""
	res, err := m.Client.GetRecentBlockhash(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
	if err == nil {
		blockHash = res.Value.Blockhash.String()
	}
	return blockHash, err
}

func (m *Provider) GetBlock() (*rpc.GetBlockResult, error) {
	res, err := m.Client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	return m.Client.GetBlock(context.TODO(), res.Context.Slot)
}

func (m *Provider) GetBlockCommitment() (*rpc.GetBlockCommitmentResult, error) {
	example, err := m.Client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}
	return m.Client.GetBlockCommitment(
		context.TODO(),
		example.Context.Slot,
	)
}

func (m *Provider) GetLatestBlockHash() (string, error) {
	blockHash := ""
	example, err := m.Client.GetLatestBlockhash(
		context.Background(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return blockHash, err
	}
	blockHash = example.Value.Blockhash.String()
	return blockHash, nil
}

func (m *Provider) GetBlockHeight() (uint64, error) {
	return m.Client.GetBlockHeight(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
}

func (m *Provider) IsBlockHashValid() bool {
	blockHash := solana.MustHashFromBase58("J7rBdM6AecPDEZp8aPq5iPSNKVkU5Q76F3oAV4eW5wsW")
	out, err := m.Client.IsBlockhashValid(
		context.TODO(),
		blockHash,
		rpc.CommitmentFinalized,
	)
	if err == nil {
		return out.Value
	}
	return false
}

func (m *Provider) GetAccountInfo(pubkey string) (token.Mint, error) {
	pubKey := solana.MustPublicKeyFromBase58(pubkey)
	res, err := m.Client.GetAccountInfo(
		context.TODO(),
		pubKey,
	)
	var mint token.Mint
	if err == nil {
		err = bin.NewBinDecoder(res.GetBinary()).Decode(&mint)
	}
	return mint, err
}

func (m *Provider) GetBlockTime() (*solana.UnixTimeSeconds, error) {
	example, err := m.Client.GetRecentBlockhash(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return nil, err
	}
	return m.Client.GetBlockTime(
		context.TODO(),
		example.Context.Slot,
	)
}

func (m *Provider) GetConfirmedBlock() (*rpc.GetBlockResult, error) {
	example, err := m.Client.GetRecentBlockhash(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return nil, err
	}

	return m.Client.GetBlockWithOpts(
		context.TODO(),
		example.Context.Slot,
		&rpc.GetBlockOpts{
			Encoding:           solana.EncodingBase64,
			Commitment:         rpc.CommitmentFinalized,
			TransactionDetails: rpc.TransactionDetailsSignatures,
			Rewards:            pointer.ToBool(false),
		},
	)
}

func (m *Provider) GetBalance(pubkey string) (*big.Float, error) {
	pubKey := solana.MustPublicKeyFromBase58(pubkey)
	out, err := m.Client.GetBalance(
		context.TODO(),
		pubKey,
		rpc.CommitmentFinalized,
	)
	var balance *big.Float
	if err == nil {
		balance = ToSol(out.Value, solana.LAMPORTS_PER_SOL)
	}
	return balance, err
}

func (m *Provider) GetSignaturesForAddress(pubkey string) ([]*rpc.TransactionSignature, error) {
	pubKey := solana.MustPublicKeyFromBase58(pubkey)
	return m.Client.GetSignaturesForAddress(
		context.TODO(),
		pubKey,
	)
}

func (m *Provider) GetFees() (uint64, error) {
	res, err := m.Client.GetFees(
		context.TODO(),
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		return 0, err
	}
	return res.Value.FeeCalculator.LamportsPerSignature, err
}

func (m *Provider) GetFeeForMessage(rawMsg string) (uint64, error) {
	res, err := m.Client.GetFeeForMessage(
		context.Background(),
		rawMsg,
		rpc.CommitmentProcessed,
	)
	if err != nil {
		return 0, err
	}
	return *res.Value, err
}

func (m *Provider) GetRecentPrioritizationFees(pubkey string) (uint64, error) {
	out, err := m.Client.GetRecentPrioritizationFees(
		context.TODO(),
		[]solana.PublicKey{
			solana.MustPublicKeyFromBase58(pubkey),
		},
	)
	if err != nil {
		return 0, err
	}
	return out[0].PrioritizationFee, nil
}

func (m *Provider) GetProgramAccounts(pubkey string) ([]*rpc.KeyedAccount, error) {
	out, err := m.Client.GetProgramAccounts(
		context.TODO(),
		solana.MustPublicKeyFromBase58(pubkey),
	)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (m *Provider) GetTokenAccountBalance(pubkey string) (*rpc.UiTokenAmount, error) {
	pubKey := solana.MustPublicKeyFromBase58(pubkey)
	out, err := m.Client.GetTokenAccountBalance(
		context.TODO(),
		pubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return nil, err
	}
	return out.Value, nil
}

func (m *Provider) GetTokenAccountsByDelegate(pubkey string, mintId string) ([]*rpc.TokenAccount, error) {
	pubKey := solana.MustPublicKeyFromBase58(pubkey)
	tokenMintId := solana.MustPublicKeyFromBase58(mintId)
	out, err := m.Client.GetTokenAccountsByDelegate(
		context.TODO(),
		pubKey,
		&rpc.GetTokenAccountsConfig{
			Mint: &tokenMintId,
		},
		nil,
	)

	if err != nil {
		return nil, err
	}
	return out.Value, nil
}

func (m *Provider) GetTokenAccountsByOwner(pubkey string) ([]token.Account, error) {
	pubKey := solana.MustPublicKeyFromBase58(pubkey)
	out, err := m.Client.GetTokenAccountsByOwner(
		context.TODO(),
		pubKey,
		&rpc.GetTokenAccountsConfig{
			Mint: solana.WrappedSol.ToPointer(),
		},
		&rpc.GetTokenAccountsOpts{
			Encoding: solana.EncodingBase64Zstd,
		},
	)
	if err != nil {
		return nil, err
	}
	tokenAccounts := make([]token.Account, 0)
	for _, rawAccount := range out.Value {
		var tokAcc token.Account

		data := rawAccount.Account.Data.GetBinary()
		dec := bin.NewBinDecoder(data)
		err := dec.Decode(&tokAcc)
		if err != nil {
			return nil, err
		}
		tokenAccounts = append(tokenAccounts, tokAcc)
	}
	return tokenAccounts, nil
}

func (m *Provider) GetTokenLargestAccounts(pubkey string) ([]*rpc.TokenLargestAccountsResult, error) {
	pubKey := solana.MustPublicKeyFromBase58(pubkey) // serum token
	out, err := m.Client.GetTokenLargestAccounts(
		context.TODO(),
		pubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return nil, err
	}
	return out.Value, nil
}

func (m *Provider) GetTokenSupply() (*rpc.UiTokenAmount, error) {
	pubKey := solana.MustPublicKeyFromBase58("SRMuApVNdxXokk5GT7XD5cUUgXMBCoAz2LHeuAoKWRt") // serum token
	out, err := m.Client.GetTokenSupply(
		context.TODO(),
		pubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return nil, err
	}
	return out.Value, nil
}

func (m *Provider) GetTransaction(hash string) (*solana.Transaction, error) {
	txSig := solana.MustSignatureFromBase58(hash)
	out, err := m.Client.GetTransaction(
		context.TODO(),
		txSig,
		&rpc.GetTransactionOpts{
			Encoding: solana.EncodingBase64,
		},
	)
	if err != nil {
		return nil, err
	}
	return solana.TransactionFromDecoder(bin.NewBinDecoder(out.Transaction.GetBinary()))
}

func (m *Provider) GetTransactionCount() (uint64, error) {
	return m.Client.GetTransactionCount(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
}
