package raydium

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	gas "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/hkxiaoyu118/blockchain/solana/config"
	"github.com/hkxiaoyu118/blockchain/solana/src"
)

type RaydiumSwap struct {
	clientRPC *rpc.Client
	account   solana.PrivateKey
}

func NewRaydiumSwap(clientRPC *rpc.Client, account solana.PrivateKey) *RaydiumSwap {
	return &RaydiumSwap{
		clientRPC: clientRPC,
		account:   account,
	}
}

func (s *RaydiumSwap) EasySwap(
	ctx context.Context,
	targetPool *Pool,
	amountIn uint64,
	amountOutMin uint64,
	fromToken string,
	fromAccount solana.PublicKey,
	toToken string,
	toAccount solana.PublicKey,
	unitPrice uint64,
	unitLimit uint32,
) (*solana.Signature, error) {
	return s.Swap(ctx, &RaydiumPoolConfig{
		AmmId:                 targetPool.ID,
		AmmAuthority:          targetPool.Authority,
		AmmOpenOrders:         targetPool.OpenOrders,
		AmmTargetOrders:       targetPool.TargetOrders,
		AmmQuantities:         config.NativeSOL,
		PoolCoinTokenAccount:  targetPool.BaseVault,
		PoolPcTokenAccount:    targetPool.QuoteVault,
		SerumProgramId:        targetPool.MarketProgramId,
		SerumMarket:           targetPool.MarketId,
		SerumBids:             targetPool.MarketBids,
		SerumAsks:             targetPool.MarketAsks,
		SerumEventQueue:       targetPool.MarketEventQueue,
		SerumCoinVaultAccount: targetPool.MarketBaseVault,
		SerumPcVaultAccount:   targetPool.MarketQuoteVault,
		SerumVaultSigner:      targetPool.MarketAuthority,
	}, amountIn, amountOutMin, fromToken, fromAccount, toToken, toAccount, unitPrice, unitLimit)
}

func (s *RaydiumSwap) Swap(
	ctx context.Context,
	pool *RaydiumPoolConfig,
	amount uint64,
	amountOutMin uint64,
	fromToken string,
	fromAccount solana.PublicKey,
	toToken string,
	toAccount solana.PublicKey,
	unitPrice uint64,
	unitLimit uint32,

) (*solana.Signature, error) {
	if amountOutMin <= 0 {
		return nil, errors.New("min swap output amount must be grater then zero, try to swap a bigger amount")
	}

	var instList []solana.Instruction
	signers := []solana.PrivateKey{s.account}

	if unitPrice != 0 {
		unitPriceInst, err := gas.NewSetComputeUnitPriceInstruction(unitPrice).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		instList = append(instList, unitPriceInst)
	}

	if unitLimit != 0 {
		unitLimit, err := gas.NewSetComputeUnitLimitInstruction(unitLimit).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		instList = append(instList, unitLimit)
	}

	tempAccount := solana.NewWallet()
	needWrapSOL := fromToken == config.NativeSOL || toToken == config.NativeSOL
	if needWrapSOL {
		rentCost, err := s.clientRPC.GetMinimumBalanceForRentExemption(
			ctx,
			config.TokenAccountSize,
			rpc.CommitmentConfirmed,
		)
		if err != nil {
			return nil, err
		}
		accountLamports := rentCost
		if fromToken == config.NativeSOL {
			// If is from a SOL account, transfer the amount
			accountLamports += amount
		}
		createInst, err := system.NewCreateAccountInstruction(
			accountLamports,
			config.TokenAccountSize,
			solana.TokenProgramID,
			s.account.PublicKey(),
			tempAccount.PublicKey(),
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		instList = append(instList, createInst)
		initInst, err := token.NewInitializeAccountInstruction(
			tempAccount.PublicKey(),
			solana.MustPublicKeyFromBase58(config.WrappedSOL),
			s.account.PublicKey(),
			solana.SysVarRentPubkey,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		instList = append(instList, initInst)
		signers = append(signers, tempAccount.PrivateKey)
		// Use this new temp account as from or to
		if fromToken == config.NativeSOL {
			fromAccount = tempAccount.PublicKey()
		}
		if toToken == config.NativeSOL {
			toAccount = tempAccount.PublicKey()
		}
	}

	instList = append(instList, NewRaydiumSwapInstruction(
		amount,
		amountOutMin,
		solana.TokenProgramID,
		solana.MustPublicKeyFromBase58(pool.AmmId),
		solana.MustPublicKeyFromBase58(pool.AmmAuthority),
		solana.MustPublicKeyFromBase58(pool.AmmOpenOrders),
		solana.MustPublicKeyFromBase58(pool.AmmTargetOrders),
		solana.MustPublicKeyFromBase58(pool.PoolCoinTokenAccount),
		solana.MustPublicKeyFromBase58(pool.PoolPcTokenAccount),
		solana.MustPublicKeyFromBase58(pool.SerumProgramId),
		solana.MustPublicKeyFromBase58(pool.SerumMarket),
		solana.MustPublicKeyFromBase58(pool.SerumBids),
		solana.MustPublicKeyFromBase58(pool.SerumAsks),
		solana.MustPublicKeyFromBase58(pool.SerumEventQueue),
		solana.MustPublicKeyFromBase58(pool.SerumCoinVaultAccount),
		solana.MustPublicKeyFromBase58(pool.SerumPcVaultAccount),
		solana.MustPublicKeyFromBase58(pool.SerumVaultSigner),
		fromAccount,
		toAccount,
		s.account.PublicKey(),
	))

	if needWrapSOL {
		closeInst, err := token.NewCloseAccountInstruction(
			tempAccount.PublicKey(),
			s.account.PublicKey(),
			s.account.PublicKey(),
			[]solana.PublicKey{},
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		instList = append(instList, closeInst)
	}

	sig, err := src.ExecuteInstructions(ctx, s.clientRPC, signers, instList...)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

/** Instructions  **/

type RaySwapInstruction struct {
	bin.BaseVariant
	InAmount                uint64
	MinimumOutAmount        uint64
	solana.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (inst *RaySwapInstruction) ProgramID() solana.PublicKey {
	return solana.MustPublicKeyFromBase58(RaydiumLiquidityPoolProgramIDV4)
}

func (inst *RaySwapInstruction) Accounts() (out []*solana.AccountMeta) {
	return inst.Impl.(solana.AccountsGettable).GetAccounts()
}

func (inst *RaySwapInstruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := bin.NewBorshEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *RaySwapInstruction) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	// Swap instruction is number 9
	err = encoder.WriteUint8(9)
	if err != nil {
		return err
	}
	err = encoder.WriteUint64(inst.InAmount, binary.LittleEndian)
	if err != nil {
		return err
	}
	err = encoder.WriteUint64(inst.MinimumOutAmount, binary.LittleEndian)
	if err != nil {
		return err
	}
	return nil
}

func NewRaydiumSwapInstruction(
	// Parameters:
	inAmount uint64,
	minimumOutAmount uint64,
	// Accounts:
	tokenProgram solana.PublicKey,
	ammId solana.PublicKey,
	ammAuthority solana.PublicKey,
	ammOpenOrders solana.PublicKey,
	ammTargetOrders solana.PublicKey,
	poolCoinTokenAccount solana.PublicKey,
	poolPcTokenAccount solana.PublicKey,
	serumProgramId solana.PublicKey,
	serumMarket solana.PublicKey,
	serumBids solana.PublicKey,
	serumAsks solana.PublicKey,
	serumEventQueue solana.PublicKey,
	serumCoinVaultAccount solana.PublicKey,
	serumPcVaultAccount solana.PublicKey,
	serumVaultSigner solana.PublicKey,
	userSourceTokenAccount solana.PublicKey,
	userDestTokenAccount solana.PublicKey,
	userOwner solana.PublicKey,
) *RaySwapInstruction {

	inst := RaySwapInstruction{
		InAmount:         inAmount,
		MinimumOutAmount: minimumOutAmount,
		AccountMetaSlice: make(solana.AccountMetaSlice, 18),
	}
	inst.BaseVariant = bin.BaseVariant{
		Impl: inst,
	}

	inst.AccountMetaSlice[0] = solana.Meta(tokenProgram)
	inst.AccountMetaSlice[1] = solana.Meta(ammId).WRITE()
	inst.AccountMetaSlice[2] = solana.Meta(ammAuthority)
	inst.AccountMetaSlice[3] = solana.Meta(ammOpenOrders).WRITE()
	inst.AccountMetaSlice[4] = solana.Meta(ammTargetOrders).WRITE()
	inst.AccountMetaSlice[5] = solana.Meta(poolCoinTokenAccount).WRITE()
	inst.AccountMetaSlice[6] = solana.Meta(poolPcTokenAccount).WRITE()
	inst.AccountMetaSlice[7] = solana.Meta(serumProgramId)
	inst.AccountMetaSlice[8] = solana.Meta(serumMarket).WRITE()
	inst.AccountMetaSlice[9] = solana.Meta(serumBids).WRITE()
	inst.AccountMetaSlice[10] = solana.Meta(serumAsks).WRITE()
	inst.AccountMetaSlice[11] = solana.Meta(serumEventQueue).WRITE()
	inst.AccountMetaSlice[12] = solana.Meta(serumCoinVaultAccount).WRITE()
	inst.AccountMetaSlice[13] = solana.Meta(serumPcVaultAccount).WRITE()
	inst.AccountMetaSlice[14] = solana.Meta(serumVaultSigner)
	inst.AccountMetaSlice[15] = solana.Meta(userSourceTokenAccount).WRITE()
	inst.AccountMetaSlice[16] = solana.Meta(userDestTokenAccount).WRITE()
	inst.AccountMetaSlice[17] = solana.Meta(userOwner).SIGNER()

	return &inst
}
