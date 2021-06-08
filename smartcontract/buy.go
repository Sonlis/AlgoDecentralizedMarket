package main 

import (
	"net/http"
	"encoding/json"
	"encoding/base64"
	"log"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
	"context"
	"crypto/ed25519"
)

func buy(w http.ResponseWriter, r *http.Request) {
	var (
		sk ed25519.PrivateKey
   		ma crypto.MultisigAccount
		tx Tx 
		algoClient AlgoClient
		err error
		txGroup txGroup
	)

    algoClient.c, err = algod.MakeClient(algodAddress, algodToken)
    if err != nil {
        log.Printf("failed to make algod client: %s\n", err)
    }
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&tx)
	if err != nil {
		log.Println(err)
	}
	genID, genHash, minFee, firstValidRound, lastValidRound := algoClient.getParams()
	genHash64 := base64.StdEncoding.EncodeToString(genHash)
	txn1, err := transaction.MakeAssetAcceptanceTxn(tx.Sender, 1, firstValidRound, lastValidRound, nil, genID, genHash64, tx.ToBuy)

	txn2, err := transaction.MakeAssetTransferTxnWithFlatFee(tx.Address, tx.Sender, "", 1, minFee, firstValidRound, lastValidRound, nil,
		genID, genHash64, tx.ToBuy)
	if err != nil {
		log.Printf("Failed to send transaction MakeAssetTransfer Txn: %s\n", err)
	}
	var txn3 types.Transaction
	if tx.AlgoAmount != 0 {
		txn3, err = transaction.MakePaymentTxnWithFlatFee(tx.Sender, tx.Address, minFee, tx.AlgoAmount, firstValidRound, lastValidRound, nil, "", genID, genHash)
		if err != nil {
			log.Printf("Error creating transaction: %s\n", err)
		}
	} else {
		txn3, err = transaction.MakeAssetTransferTxnWithFlatFee(tx.Sender, tx.Address, "", tx.Amount, minFee, firstValidRound, lastValidRound, nil,
		genID, genHash64, tx.ToPay)
		if err != nil {
			log.Printf("Failed to send transaction MakeAssetTransfer Txn: %s\n", err)
		}
	}
	teal := readTeal(tx.Address)
	response, err := algoClient.c.TealCompile(teal).Do(context.Background())
	if err != nil {
		log.Println("Error compiling:", err)
	}
	
	program, err :=  base64.StdEncoding.DecodeString(response.Result)	
    lsig, err := crypto.MakeLogicSig(program, nil, sk, ma)

	gid, err := crypto.ComputeGroupID([]types.Transaction{txn1, txn2, txn3})
	txn1.Group = gid
	txn2.Group = gid
	txn3.Group = gid
	_, stx2, err := crypto.SignLogicsigTransaction(lsig, txn2)
	txGroup.FirstTx = txn1
	txGroup.SecondTx = stx2
	txGroup.ThirdTx = txn3
	if err := json.NewEncoder(w).Encode(txGroup); err != nil {
        log.Printf("Error sending back response: %s", err)
    }
}

func withdrawAlgos(w http.ResponseWriter, r *http.Request) {
	var (
		algoClient AlgoClient
		withdraw Withdraw
		txGroup TransactionGroup
		err error
		sk ed25519.PrivateKey
    	ma crypto.MultisigAccount
	)
	algoClient.c, err = algod.MakeClient(algodAddress, algodToken)
    if err != nil {
        log.Printf("failed to make algod client: %s\n", err)
    }
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&withdraw)
	if err != nil {
		log.Println(err)
	}
	genID, genHash, minFee, firstValidRound, lastValidRound := algoClient.getParams()
	txn1, err := transaction.MakePaymentTxnWithFlatFee(withdraw.Address, withdraw.Creator, minFee, withdraw.Algo, firstValidRound, lastValidRound, nil, "", genID, genHash)

	txn2, err := transaction.MakePaymentTxnWithFlatFee(withdraw.Creator, withdraw.Address, minFee, 1, firstValidRound, lastValidRound, nil, "", genID, genHash)
	if err != nil {
		log.Printf("Failed to send transaction MakeAssetTransfer Txn: %s\n", err)
	}
	teal := readTeal(withdraw.Address)
	response, err := algoClient.c.TealCompile(teal).Do(context.Background())
	if err != nil {
		log.Println("Error compiling:", err)
	}

	program, err :=  base64.StdEncoding.DecodeString(response.Result)	
    lsig, err := crypto.MakeLogicSig(program, nil, sk, ma)

	gid, err := crypto.ComputeGroupID([]types.Transaction{txn1, txn2})
	txn1.Group = gid
	txn2.Group = gid
	_, stx1, err := crypto.SignLogicsigTransaction(lsig, txn1)
	txGroup.FirstTx = stx1
	txGroup.SecondTx = txn2
	if err := json.NewEncoder(w).Encode(txGroup); err != nil {
        log.Printf("Error sending back response: %s", err)
    }
}

func withdrawAssets(w http.ResponseWriter, r *http.Request) {
	var (
		algoClient AlgoClient
		withdraw WithdrawAsset
		txGroup TransactionGroup
		err error
		sk ed25519.PrivateKey
    	ma crypto.MultisigAccount
	) 
	algoClient.c, err = algod.MakeClient(algodAddress, algodToken)
    if err != nil {
        log.Printf("failed to make algod client: %s\n", err)
    }
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&withdraw)
	if err != nil {
		log.Println(err)
	}
	genID, genHash, minFee, firstValidRound, lastValidRound := algoClient.getParams()
	genHash64 := base64.StdEncoding.EncodeToString(genHash)
	txn1, err := transaction.MakeAssetTransferTxnWithFlatFee(withdraw.Address, withdraw.Creator, "", withdraw.AssetAmount, minFee, firstValidRound, lastValidRound, nil,
		genID, genHash64, withdraw.AssetID)
	if err != nil {
		log.Printf("Failed to make transaction MakeAssetTransfer: %s", err)
	}

	txn2, err := transaction.MakePaymentTxnWithFlatFee(withdraw.Creator, withdraw.Address, minFee, 1, firstValidRound, lastValidRound, nil, "", genID, genHash)
	if err != nil {
		log.Printf("Failed to send transaction MakePayment Txn: %s\n", err)
	}
	teal := readTeal(withdraw.Address)
	response, err := algoClient.c.TealCompile(teal).Do(context.Background())
	if err != nil {
		log.Println("Error compiling:", err)
	}
	program, err :=  base64.StdEncoding.DecodeString(response.Result)	
    lsig, err := crypto.MakeLogicSig(program, nil, sk, ma)
	gid, err := crypto.ComputeGroupID([]types.Transaction{txn1, txn2})
	txn1.Group = gid
	txn2.Group = gid
	_, stx1, err := crypto.SignLogicsigTransaction(lsig, txn1)
	txGroup.FirstTx = stx1
	txGroup.SecondTx = txn2
	if err := json.NewEncoder(w).Encode(txGroup); err != nil {
        log.Printf("Error sending back response: %s", err)
    }
}


