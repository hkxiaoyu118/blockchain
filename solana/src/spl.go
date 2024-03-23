package src

import (
	"context"
	"github.com/gagliardetto/solana-go"
	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

// CreateAssociatedTokenAddress 创建ATA账户
func CreateAssociatedTokenAddress(client *rpc.Client, wallet solana.PrivateKey, mint solana.PublicKey) (*solana.Signature, error) {
	var instList []solana.Instruction
	signers := []solana.PrivateKey{wallet}

	newAtaAccountInst, err := ata.NewCreateInstruction(wallet.PublicKey(), wallet.PublicKey(), mint).ValidateAndBuild()
	if err != nil {
		return nil, err
	}
	instList = append(instList, newAtaAccountInst)

	signers = append(signers)

	return ExecuteInstructions(context.TODO(), client, signers, instList...)
}

// CloseAssociatedTokenAddress 关闭ATA账户
func CloseAssociatedTokenAddress(client *rpc.Client, wallet solana.PrivateKey, mint solana.PublicKey) (*solana.Signature, error) {
	var instList []solana.Instruction
	signers := []solana.PrivateKey{wallet}

	ataAccount, _, err := solana.FindAssociatedTokenAddress(
		wallet.PublicKey(),
		mint,
	)
	if err != nil {
		return nil, err
	}

	closeInst, err := token.NewCloseAccountInstruction(
		ataAccount,
		wallet.PublicKey(),
		wallet.PublicKey(),
		[]solana.PublicKey{},
	).ValidateAndBuild()

	instList = append(instList, closeInst)

	signers = append(signers)

	return ExecuteInstructions(context.TODO(), client, signers, instList...)
}
