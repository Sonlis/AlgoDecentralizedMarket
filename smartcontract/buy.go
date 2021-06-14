
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
        log.Printf("Failed to make algod client: %s\n", err)
        http.Error(w, err.Error(), 500)
        return
    }
    decoder := json.NewDecoder(r.Body)
    err = decoder.Decode(&tx)
    if err != nil {
        log.Printf("Failed to decode Buy request body: %v\n", err)
        http.Error(w, err.Error(), 400)
        return
    }
    genID, genHash, minFee, firstValidRound, lastValidRound, err := algoClient.getParams()
	if err != nil {
		log.Printf("Error retrieving network parameters: %v\n", err)
		http.Error(w, err.Error(), 500)
        return
	}
    genHash64 := base64.StdEncoding.EncodeToString(genHash)
    txn1, err := transaction.MakeAssetAcceptanceTxn(tx.Sender, 1, firstValidRound, lastValidRound, nil, genID, genHash64, tx.ToBuy)
    if err != nil {
        log.Printf("Error making opt-in transaction of account %s for asset %d: %v\n", tx.Sender, tx.ToBuy, err)
        http.Error(w, err.Error(), 400)
        return
    }
    txn2, err := transaction.MakeAssetTransferTxnWithFlatFee(tx.Address, tx.Sender, "", 1, minFee, firstValidRound, lastValidRound, nil,
        genID, genHash64, tx.ToBuy)
    if err != nil {
        log.Printf("Failed making asset transfer for asset %d from %s to %s: %v\n\n", tx.ToBuy, tx.Address, tx.Sender, err)
        http.Error(w, err.Error(), 400)
        return
    }
    var txn3 types.Transaction
    if tx.AlgoAmount != 0 {
        txn3, err = transaction.MakePaymentTxnWithFlatFee(tx.Sender, tx.Address, minFee, tx.AlgoAmount, firstValidRound, lastValidRound, nil, "", genID, genHash)
        if err != nil {
            log.Printf("Error making payment transaction of %d Algos from %s to %s: %v\n\n", tx.AlgoAmount, tx.Sender, tx.Address, err)
            http.Error(w, err.Error(), 400)
            return
        }
    } else {
        txn3, err = transaction.MakeAssetTransferTxnWithFlatFee(tx.Sender, tx.Address, "", tx.Amount, minFee, firstValidRound, lastValidRound, nil,
        genID, genHash64, tx.ToPay)
        if err != nil {
            log.Printf("Failed making asset transfer of asset %d from %s to %s: %v\n\n", tx.ToPay, tx.Sender, tx.Address, err)
            http.Error(w, err.Error(), 400)
            return
        }
    }
    txn4, err := transaction.MakePaymentTxnWithFlatFee(tx.Sender, tx.Address, minFee, minFee, firstValidRound, lastValidRound, nil, "", genID, genHash)
    if err != nil {
        log.Printf("Failed making payment transaction of amount %d from %s to %s: %v\n\n", minFee, tx.Sender, tx.Address, err)
            http.Error(w, err.Error(), 400)
            return
    }
    teal, err := readTeal(tx.Address)
	if err != nil {
		log.Printf("Error reading teal: %v\n", err)
		http.Error(w, err.Error(), 500)
        return
	}
    response, err := algoClient.c.TealCompile(teal).Do(context.Background())
    if err != nil {
        log.Printf("Error compiling teal: %v\n", err)
        http.Error(w, err.Error(), 400)
        return
    }
    
    program, err :=  base64.StdEncoding.DecodeString(response.Result)   
    lsig, err := crypto.MakeLogicSig(program, nil, sk, ma)
    if err != nil {
        log.Printf("Error signing teal: %v\n", err)
        http.Error(w, err.Error(), 400)
        return
    }

    gid, err := crypto.ComputeGroupID([]types.Transaction{txn1, txn2, txn3, txn4})
    if err != nil {
        log.Printf("Error calculating the group ID: %v\n", err)
        http.Error(w, err.Error(), 400)
        return
    }
    txn1.Group = gid
    txn2.Group = gid
    txn3.Group = gid
    txn4.Group = gid
    _, stx2, err := crypto.SignLogicsigTransaction(lsig, txn2)
    if err != nil {
        log.Printf("Error signing with logic: %v\n", err)
        http.Error(w, err.Error(), 400)
        return
    }
    txGroup.FirstTx = txn1
    txGroup.SecondTx = stx2
    txGroup.ThirdTx = txn3
    txGroup.ForthTx = txn4

    if err := json.NewEncoder(w).Encode(txGroup); err != nil {
        log.Printf("Error sending back response: %v\n", err)
    }
}

func withdraw(w http.ResponseWriter, r *http.Request) {
    var (
        algoClient AlgoClient
        withdraw Withdraw
        txGroup TransactionGroup
        err error
        sk ed25519.PrivateKey
        ma crypto.MultisigAccount
		txn1 types.Transaction
    )
    algoClient.c, err = algod.MakeClient(algodAddress, algodToken)
    if err != nil {
        log.Printf("Failed to make algod client: %s\n", err)
		http.Error(w, err.Error(), 500)
        return
    }
    decoder := json.NewDecoder(r.Body)
    err = decoder.Decode(&withdraw)
    if err != nil {
        log.Printf("Error decoding request's body: %v\n", err)
		http.Error(w, err.Error(), 400)
        return
    }
    genID, genHash, minFee, firstValidRound, lastValidRound, err := algoClient.getParams()
	if err != nil {
		log.Printf("Error getting suggested parameters: %v\n", err)
		http.Error(w, err.Error(), 500)
        return
	}
	genHash64 := base64.StdEncoding.EncodeToString(genHash)

	if (withdraw.Algo != 0){
		txn1, err = transaction.MakePaymentTxnWithFlatFee(withdraw.Address, withdraw.Creator, minFee, withdraw.Algo, firstValidRound, lastValidRound, nil, "", genID, genHash)
		if err != nil {
			log.Printf("Error making payment transaction of %d Algos from %s to %s: %v\n", withdraw.Algo, withdraw.Address, withdraw.Creator, err)
			http.Error(w, err.Error(), 400)
        	return
	}
	} else {
		txn1, err = transaction.MakeAssetTransferTxnWithFlatFee(withdraw.Address, withdraw.Creator, "", withdraw.AssetAmount, minFee, firstValidRound, lastValidRound, nil,
        genID, genHash64, withdraw.AssetID)
    	if err != nil {
        	log.Printf("Failed to transfer asset %d from %s to %s: %v\n", withdraw.AssetID, withdraw.Address, withdraw.Creator, err)
			http.Error(w, err.Error(), 400)
			return
    }
	}
    txn2, err := transaction.MakePaymentTxnWithFlatFee(withdraw.Creator, withdraw.Address, minFee, minFee, firstValidRound, lastValidRound, nil, "", genID, genHash)
    if err != nil {
        log.Printf("Error making payment transaction of %d Algo from %s to %s: %v\n\n", minFee, withdraw.Creator, withdraw.Address,err)
		http.Error(w, err.Error(), 400)
        return
    }
    teal, err := readTeal(withdraw.Address)
	if err != nil {
		log.Printf("Error reading teal: %v\n", err)
		http.Error(w, err.Error(), 500)
        return
	}
    response, err := algoClient.c.TealCompile(teal).Do(context.Background())
    if err != nil {
        log.Printf("Error compiling teal: %v\n", err)
		http.Error(w, err.Error(), 500)
        return
    }

    program, err :=  base64.StdEncoding.DecodeString(response.Result)   
    lsig, err := crypto.MakeLogicSig(program, nil, sk, ma)
	if err != nil {
		log.Printf("Error signing with logic: %v\n", err)
		http.Error(w, err.Error(), 400)
        return
	}

    gid, err := crypto.ComputeGroupID([]types.Transaction{txn1, txn2})
	if err != nil {
		log.Printf("Error calculating group ID: %v\n", err)
		http.Error(w, err.Error(), 500)
        return
	}
    txn1.Group = gid
    txn2.Group = gid
    _, stx1, err := crypto.SignLogicsigTransaction(lsig, txn1)
	if err != nil {
		log.Printf("Error signing with logic: %v\n\n", err)
		http.Error(w, err.Error(), 500)
        return
	}
    txGroup.FirstTx = stx1
    txGroup.SecondTx = txn2
    if err := json.NewEncoder(w).Encode(txGroup); err != nil {
        log.Printf("Error sending back response: %s", err)
    }
}
