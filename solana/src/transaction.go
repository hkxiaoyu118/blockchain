package src

import (
	"context"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func BuildTransacion(ctx context.Context, clientRPC *rpc.Client, signers []solana.PrivateKey, instrs ...solana.Instruction) (*solana.Transaction, error) {
	recent, err := clientRPC.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	tx, err := solana.NewTransaction(
		instrs,
		recent.Value.Blockhash,
		solana.TransactionPayer(signers[0].PublicKey()),
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			for _, payer := range signers {
				if payer.PublicKey().Equals(key) {
					return &payer
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func ExecuteInstructions(
	ctx context.Context,
	clientRPC *rpc.Client,
	signers []solana.PrivateKey,
	instrs ...solana.Instruction,
) (*solana.Signature, error) {

	tx, err := BuildTransacion(ctx, clientRPC, signers, instrs...)
	if err != nil {
		return nil, err
	}

	sig, err := clientRPC.SendTransactionWithOpts(
		ctx,
		tx,
		rpc.TransactionOpts{
			SkipPreflight:       false,
			PreflightCommitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		return nil, err
	}

	return &sig, nil
}

func ExecuteInstructionsAndWaitConfirm(
	ctx context.Context,
	clientRPC *rpc.Client,
	RPCWs string,
	signers []solana.PrivateKey,
	instrs ...solana.Instruction,
) (*solana.Signature, error) {

	tx, err := BuildTransacion(ctx, clientRPC, signers, instrs...)
	if err != nil {
		return nil, err
	}

	clientWS, err := ws.Connect(ctx, RPCWs)
	if err != nil {
		return nil, err
	}

	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		clientRPC,
		clientWS,
		tx,
	)
	if err != nil {
		return nil, err
	}

	return &sig, nil
}
