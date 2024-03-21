package src

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func AccountSubscribe(rpcUrl string, pubkey string) {
	client, err := ws.Connect(context.Background(), rpcUrl)
	if err != nil {
		panic(err)
	}
	program := solana.MustPublicKeyFromBase58(pubkey) // serum

	sub, err := client.AccountSubscribe(
		program,
		"",
	)
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			panic(err)
		}
		spew.Dump(got)
	}
}

func LogsSubscribe(rpcUrl string, pubkey string) {
	client, err := ws.Connect(context.Background(), rpcUrl)
	if err != nil {
		panic(err)
	}
	program := solana.MustPublicKeyFromBase58(pubkey) // serum

	// Subscribe to log events that mention the provided pubkey:
	sub, err := client.LogsSubscribeMentions(
		program,
		rpc.CommitmentRecent,
	)
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			panic(err)
		}
		spew.Dump(got)
	}
}

func ProgramSubscribe(rpcUrl string, pubkey string) {
	client, err := ws.Connect(context.Background(), rpcUrl)
	if err != nil {
		panic(err)
	}
	program := solana.MustPublicKeyFromBase58(pubkey) // token

	sub, err := client.ProgramSubscribeWithOpts(
		program,
		rpc.CommitmentRecent,
		solana.EncodingBase64Zstd,
		nil,
	)
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			panic(err)
		}
		spew.Dump(got)

		decodedBinary := got.Value.Account.Data.GetBinary()
		if decodedBinary != nil {
			// spew.Dump(decodedBinary)
		}

		// or if you requested solana.EncodingJSONParsed and it is supported:
		rawJSON := got.Value.Account.Data.GetRawJSON()
		if rawJSON != nil {
			// spew.Dump(rawJSON)
		}
	}
}

func RootSubscribe(rpcUrl string) {
	client, err := ws.Connect(context.Background(), rpcUrl)
	if err != nil {
		panic(err)
	}

	sub, err := client.RootSubscribe()
	if err != nil {
		panic(err)
	}

	for {
		got, err := sub.Recv()
		if err != nil {
			panic(err)
		}
		spew.Dump(got)
	}
}

func SignatureSubscribe(rpcUrl string, hash string) {
	client, err := ws.Connect(context.Background(), rpcUrl)
	if err != nil {
		panic(err)
	}

	txSig := solana.MustSignatureFromBase58(hash)

	sub, err := client.SignatureSubscribe(
		txSig,
		"",
	)
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			panic(err)
		}
		spew.Dump(got)
	}
}

func SlotSubscribe(rpcUrl string) {
	client, err := ws.Connect(context.Background(), rpcUrl)
	if err != nil {
		panic(err)
	}

	sub, err := client.SlotSubscribe()
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			panic(err)
		}
		spew.Dump(got)
	}
}

func VoteSubscribe(rpcUrl string) {
	client, err := ws.Connect(context.Background(), rpcUrl)
	if err != nil {
		panic(err)
	}

	// NOTE: this subscription must be enabled by the node you're connecting to.
	// This subscription is disabled by default.
	sub, err := client.VoteSubscribe()
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			panic(err)
		}
		spew.Dump(got)
	}
}
