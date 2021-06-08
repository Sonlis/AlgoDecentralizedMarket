package main 

import (
	"context"
    "log"
    "fmt"
    "github.com/algorand/go-algorand-sdk/client/v2/algod"
	"encoding/base64"
	"encoding/json"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
	"net/http"
	"crypto/ed25519"
)

func createEscrow(w http.ResponseWriter, r *http.Request) {
	var p Choice
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		log.Println("Error decoding body:", err)
	}
	log.Printf("%#v\n", p)
	account1 := getAddress(mnemonic1)
	sk1, err := mnemonic.ToPrivateKey(mnemonic1)

	var algoClient AlgoClient
    algoClient.c, err = algod.MakeClient(algodAddress, algodToken)
    if err != nil {
        fmt.Printf("failed to make algod client: %s\n", err) 
    }
	lsig, addr := algoClient.generateTeal(p)

	genID, genHash, minFee, firstValidRound, lastValidRound := algoClient.getParams()
	var amount uint64 = 1000000
	tx1, err := transaction.MakePaymentTxnWithFlatFee(account1, addr, minFee, amount, firstValidRound, lastValidRound, nil, "", genID, genHash)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
	}
	log.Printf("...tx1: from %s to %s for %v microAlgos\n", account1, addr, amount)
	
	log.Println("Signing transactions...")
	_, stx1, err := crypto.SignTransaction(sk1, tx1)
	if err != nil {
		log.Printf("Failed to sign transaction: %s\n", err)
	}
	
	pendingTxID, err := algoClient.c.SendRawTransaction(stx1).Do(context.Background())
	if err != nil {
		log.Printf("failed to send transaction: %s\n", err)
	}
	err = algoClient.opTin(p.PaymentAssetId, addr, account1, lsig, sk1)
	err = algoClient.opTin(p.SecondPaymentAssetId, addr, account1, lsig, sk1)
	err = algoClient.waitForConfirmation(pendingTxID)
	txGroup, err := fundEscrow(p, addr, &algoClient, w)

	if err != nil {
		log.Printf("Error funding the escrow: %v", err)
	}
	var toSave Selling
	toSave.Asset = p.AssetId
	toSave.Address = addr 
	toSave.FirstAsset = p.PaymentAssetId
	toSave.FAmount = p.PaymentAssetAmount
	toSave.SecondAsset = p.SecondPaymentAssetId
	toSave.SAmount = p.SecondPaymentAssetAmount
	toSave.AlgoAmount = p.AlgoAmount
	toSave.CreatorAddress = p.CreatorAddress
	sellings = append(sellings, toSave)
	if err := json.NewEncoder(w).Encode(txGroup); err != nil {
        log.Printf("Error sending back response: %s", err)
    }
}

func (c *AlgoClient) opTin(index uint64, addr string, account string, lsig types.LogicSig, sk1 ed25519.PrivateKey) (err error) {
	genID, genHash, minFee, firstValidRound, lastValidRound := c.getParams()
	genHash64 := base64.StdEncoding.EncodeToString(genHash)
	txn1, err := transaction.MakeAssetAcceptanceTxn(addr, 0, firstValidRound, lastValidRound, nil, genID, genHash64, index)
	if err != nil {
		fmt.Printf("Failed to send transaction MakeAssetAcceptanceTxn: %s\n", err)
	}
	txn2, err := transaction.MakePaymentTxnWithFlatFee(account, addr, minFee, 1000, firstValidRound, lastValidRound, nil, "", genID, genHash)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
	}
	gid, err := crypto.ComputeGroupID([]types.Transaction{txn1, txn2})
	txn1.Group = gid
	txn2.Group = gid
	_, stx1, err := crypto.SignLogicsigTransaction(lsig, txn1)
	_, stx2, err := crypto.SignTransaction(sk1, txn2)
	var signedGroup []byte
	signedGroup = append(signedGroup, stx1...)
	signedGroup = append(signedGroup, stx2...)
	_, err = c.c.SendRawTransaction(signedGroup).Do(context.Background())
	if err != nil {
		log.Printf("failed to send transaction: %s\n", err)
	}
	return err

}
func fundEscrow(p Choice, addr string, c *AlgoClient, w http.ResponseWriter) (txGroup TransactionGroup, err error) {
	genID, genHash, _, firstValidRound, lastValidRound := c.getParams()
	genHash64 := base64.StdEncoding.EncodeToString(genHash)
	lsig, addr := c.generateTeal(p)
	txn1, err := transaction.MakeAssetAcceptanceTxn(addr, 0, firstValidRound, lastValidRound, nil, genID, genHash64, p.AssetId)
	if err != nil {
		fmt.Printf("Failed to send transaction MakeAssetAcceptanceTxn: %s\n", err)
	}
	txn2, err := transaction.MakeAssetTransferTxn(p.CreatorAddress, addr, "", p.AssetAmount, 0, firstValidRound, lastValidRound, nil, genID, genHash64, p.AssetId)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %s\n", err)
	}
	gid, err := crypto.ComputeGroupID([]types.Transaction{txn1, txn2})
	txn1.Group = gid
	txn2.Group = gid
	fmt.Println("...computed groupId: ", gid)
	log.Println(gid)
	_, stx1, err := crypto.SignLogicsigTransaction(lsig, txn1)

	txGroup.FirstTx = stx1 
	txGroup.SecondTx = txn2
	return txGroup, err
}
