package main

import (
	"context"
    "log"
	"fmt"
	"crypto/ed25519"
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
	"io/ioutil"
)

const algodToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const algodAddress = "http://localhost:4001"
const indexerAddress = "http://localhost:8980"
const indexerToken = ""

func (c *AlgoClient) waitForConfirmation(txID string) (err error) {
	status, err := c.c.Status().Do(context.Background())
	if err != nil {
		log.Printf("error getting algod status: %s\n", err)
		
	}
	lastRound := status.LastRound
	for {
		pt, _, err := c.c.PendingTransactionInformation(txID).Do(context.Background())
		if err != nil {
			fmt.Printf("error getting pending transaction: %s\n", err)
		}
		if pt.ConfirmedRound > 0 {
			fmt.Printf("Transaction "+txID+" confirmed in round %d\n", pt.ConfirmedRound)
			return err
		}
		fmt.Printf("...waiting for confirmation\n")
		lastRound++
		status, err = c.c.StatusAfterBlock(lastRound).Do(context.Background())
	}
}

func getAddress(mn string) string {
	sk, err := mnemonic.ToPrivateKey(mn)
	if err != nil {
		log.Printf("error recovering account: %s\n", err)
		return ""
	}
	
	pk := sk.Public()
	var a types.Address
	cpk := pk.(ed25519.PublicKey)
	copy(a[:], cpk[:])
	address := a.String()
	return address
}

func (c *AlgoClient) getParams() (genID string, genHash []byte, minFee, firstValidRound, lastValidRound uint64, err error)  {
	txParams, err := c.c.SuggestedParams().Do(context.Background())
	if err != nil {
		return 
	}
	minFee = 1000
	genID = txParams.GenesisID
	genHash = txParams.GenesisHash
	firstValidRound = uint64(txParams.FirstRoundValid)
	lastValidRound = uint64(txParams.LastRoundValid)
	return 
}

func (c *AlgoClient) generateTeal(p Choice) (lsig types.LogicSig, addr string, err error) {
	var (
		sk ed25519.PrivateKey
    	ma crypto.MultisigAccount
		teal []byte
	)
	account1 := getAddress(mnemonic1)
	if p.SecondPaymentAssetId == 0 {
		teal = []byte(fmt.Sprintf("#pragma version 3\ngtxn 0 AssetCloseTo\nglobal ZeroAddress\n==\nassert \ngtxn 1 AssetCloseTo\nglobal ZeroAddress\n==\nassert \ngtxn 0 RekeyTo\nglobal ZeroAddress\n==\nassert \ngtxn 1 RekeyTo\nglobal ZeroAddress\n==\nassert\ngtxn 0 AssetCloseTo\nglobal ZeroAddress\n==\nassert\ngtxn 1 AssetCloseTo\nglobal ZeroAddress\n==\nassert\nglobal GroupSize\nint 2\n==\nbnz optIn\nglobal GroupSize\nint 4\n==\nassert\ngtxn 2 AssetCloseTo\nglobal ZeroAddress\n==\nassert \ngtxn 2 AssetCloseTo\nglobal ZeroAddress\n==\nassert \ngtxn 2 RekeyTo\nglobal ZeroAddress\n==\nassert\ngtxn 3 AssetCloseTo\nglobal ZeroAddress\n==\nassert \ngtxn 3 AssetCloseTo\nglobal ZeroAddress\n==\nassert \ngtxn 3 RekeyTo\nglobal ZeroAddress\n==\ngtxn 0 Fee\nint 2000\n<=\nassert\ngtxn 1 Fee \nint 2000\n<=\nassert\ngtxn 2 Fee \nint 2000\n<=\nassert\ngtxn 3 Fee \nint 2000\n<=\nassert\ngtxn 3 TypeEnum \nint 1\n==\nassert\ngtxn 1 Fee \ngtxn 3 Amount\n==\nassert\ngtxn 3 Receiver \ngtxn 1 Sender \n==\nassert\ngtxn 0 TypeEnum\nint 4\n==\ngtxn 1 TypeEnum\nint 4\n==\n&&\ngtxn 2 TypeEnum\nint 4\n==\nbz algoPayment\ngtxn 0 AssetAmount\nint 0\n==\n&&\ngtxn 1 AssetAmount\nint 1\n==\n&&\ngtxn 0 Sender\ngtxn 0 AssetReceiver\n==\n&&\ngtxn 0 XferAsset\nint %d\n==\n&&\ngtxn 1 XferAsset\nint %d\n==\n&&\ngtxn 2 XferAsset \nint %d\n==\n&&\ngtxn 2 AssetAmount\nint %d\n==\n&&\nreturn\nalgoPayment:\ngtxn 2 TypeEnum\nint 1\n==\n&&\ngtxn 2 Amount\nint %d\n==\n&&\nreturn\noptIn:\nglobal GroupSize\nint 2\n==\ngtxn 1 TypeEnum \nint 4\n==\nbz secondOptin\ngtxn 0 TypeEnum\nint 4\n==\n&&\ngtxn 0 XferAsset\nint %d\n==\n&&\ngtxn 1 XferAsset \nint %d\n==\n&&\ngtxn 0 AssetAmount\nint 0\n==\n&&\ngtxn 0 Sender\ngtxn 0 AssetReceiver\n==\n&&\ngtxn 1 Sender \naddr %s\n==\n&&\nreturn \nsecondOptin:\ngtxn 0 TypeEnum\nint 4\n==\nbz withdraw\ngtxn 0 AssetAmount\nint 0\n==\nbz withdraw\ngtxn 0 XferAsset\nint %d\n==\ngtxn 1 TypeEnum \nint 1\n==\n&&\ngtxn 0 Sender\ngtxn 0 AssetReceiver\n==\n&&\ngtxn 1 Sender \naddr %s\n==\n&&\nreturn\nwithdraw:\ngtxn 1 Amount \ngtxn 0 Fee \n==\ngtxn 0 Receiver\naddr %s\n==\nbz assetWithdraw\ngtxn 1 Sender\naddr %s\n==\n&&\nreturn\nassetWithdraw:\ngtxn 0 AssetReceiver\naddr %s\n==\ngtxn 1 Sender \naddr %s\n==\n&&\nreturn", p.AssetId, p.AssetId, p.PaymentAssetId, p.PaymentAssetAmount, p.AlgoAmount, p.AssetId, p.AssetId, p.CreatorAddress, p.PaymentAssetId, account1, p.CreatorAddress, p.CreatorAddress, p.CreatorAddress, p.CreatorAddress))
	} else {
		teal = []byte(fmt.Sprintf("#pragma version 3\ngtxn 0 AssetCloseTo\nglobal ZeroAddress\n==\nassert \ngtxn 1 AssetCloseTo\nglobal ZeroAddress\n==\nassert \ngtxn 0 RekeyTo\nglobal ZeroAddress\n==\nassert \ngtxn 1 RekeyTo\nglobal ZeroAddress\n==\nassert\ngtxn 0 AssetCloseTo\nglobal ZeroAddress\n==\nassert\ngtxn 1 AssetCloseTo\nglobal ZeroAddress\n==\nassert\ngtxn 0 Fee\nint 2000\n<=\nassert\ngtxn 1 Fee \nint 2000\n<=\nassert\nglobal GroupSize\nint 2\n==\nbnz optIn\nglobal GroupSize\nint 4\n==\nassert\ngtxn 2 AssetCloseTo\nglobal ZeroAddress\n==\nassert \ngtxn 2 AssetCloseTo\nglobal ZeroAddress\n==\nassert \ngtxn 2 RekeyTo\nglobal ZeroAddress\n==\nassert\ngtxn 3 AssetCloseTo\nglobal ZeroAddress\n==\nassert \ngtxn 3 AssetCloseTo\nglobal ZeroAddress\n==\nassert \ngtxn 3 RekeyTo\nglobal ZeroAddress\n==\nassert\ngtxn 2 Fee \nint 2000\n<=\nassert\ngtxn 3 Fee \nint 2000\n<=\nassert\ngtxn 3 TypeEnum \nint 1\n==\nassert\ngtxn 1 Fee\ngtxn 3 Amount\n==\nassert\ngtxn 3 Receiver \ngtxn 1 Sender \n==\nassert\ngtxn 0 TypeEnum\nint 4\n==\ngtxn 1 TypeEnum\nint 4\n==\n&&\ngtxn 2 TypeEnum\nint 4\n==\nbz algoPayment\ngtxn 0 AssetAmount\nint 0\n==\n&&\ngtxn 1 AssetAmount\nint 1\n==\n&&\ngtxn 0 Sender\ngtxn 0 AssetReceiver\n==\n&&\ngtxn 0 XferAsset\nint %d\n==\n&&\ngtxn 1 XferAsset\nint %d\n==\n&&\ngtxn 2 XferAsset \nint %d\n==\nbz secondAsset\ngtxn 2 AssetAmount\nint %d\n==\n&&\nreturn\nsecondAsset:\ngtxn 2 XferAsset\nint %d\n==\nassert\ngtxn 2 AssetAmount\nint %d\n==\n&&\nreturn\nalgoPayment:\ngtxn 2 TypeEnum\nint 1\n==\n&&\ngtxn 2 Amount\nint %d\n==\n&&\nreturn\noptIn:\nglobal GroupSize\nint 2\n==\ngtxn 1 TypeEnum \nint 4\n==\nbz secondOptin\ngtxn 0 TypeEnum\nint 4\n==\n&&\ngtxn 0 XferAsset\nint %d\n==\n&&\ngtxn 1 XferAsset \nint %d\n==\n&&\ngtxn 0 AssetAmount\nint 0\n==\n&&\ngtxn 0 Sender\ngtxn 0 AssetReceiver\n==\n&&\ngtxn 1 Sender \naddr %s\n==\n&&\nreturn \nsecondOptin:\ngtxn 0 XferAsset\nint %d\n==\nbz thirdOptin\ngtxn 1 TypeEnum \nint 1\n==\ngtxn 0 AssetAmount\nint 0\n==\nbz withdraw\ngtxn 0 Sender\ngtxn 0 AssetReceiver\n==\n&&\ngtxn 1 Sender \naddr %s\n==\n&&\nreturn\nthirdOptin:\ngtxn 0 XferAsset\nint %d\n==\nbz withdraw\ngtxn 1 TypeEnum \nint 1\n==\ngtxn 0 AssetAmount\nint 0\n==\nbz withdraw\ngtxn 0 Sender\ngtxn 0 AssetReceiver\n==\n&&\ngtxn 1 Sender \naddr %s\n==\n&&\nreturn\nwithdraw:\ngtxn 1 Amount \ngtxn 0 Fee \n==\ngtxn 0 Receiver\naddr %s\n==\nbz assetWithdraw\ngtxn 1 Sender\naddr %s\n==\n&&\nreturn\nassetWithdraw:\ngtxn 0 AssetReceiver\naddr %s\n==\n&&\ngtxn 1 Sender \naddr %s\n==\n&&\nreturn", p.AssetId, p.AssetId, p.PaymentAssetId, p.PaymentAssetAmount, p.SecondPaymentAssetId, p.SecondPaymentAssetAmount, p.AlgoAmount, p.AssetId, p.AssetId, p.CreatorAddress, p.PaymentAssetId, account1, p.SecondPaymentAssetId, p.CreatorAddress, p.CreatorAddress, p.CreatorAddress, p.CreatorAddress, p.CreatorAddress ))
	}

	response, err := c.c.TealCompile(teal).Do(context.Background())
	if err != nil {
		return
	}
	program, err :=  base64.StdEncoding.DecodeString(response.Result)	
    lsig, err = crypto.MakeLogicSig(program, nil, sk, ma)
	if err != nil {
		return 
	}
	addr = crypto.LogicSigAddress(lsig).String()
	err = ioutil.WriteFile(addr, teal, 0644)
	if err != nil {
		return
	}
	return
}

func readTeal(addr string) (b []byte, err error) {
	b, err = ioutil.ReadFile(addr)
    if err != nil {
        log.Printf("Error reading teal %s", err)
    }
	return
}