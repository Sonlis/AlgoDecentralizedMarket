package main 

import (
    "context"
    "log"
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

func activateEscrow(w http.ResponseWriter, r *http.Request) {
    var p Choice
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&p)
    if err != nil {
        log.Printf("Error decoding body: %v\n", err)
        http.Error(w, err.Error(), 400)
        return
    }
    var algoClient AlgoClient
    algoClient.c, err = algod.MakeClient(algodAddress, algodToken)
    if err != nil {
        log.Printf("failed to make algod client: %v\n\n", err) 
        http.Error(w, err.Error(), 500)
        return
    }
    _, addr, err := algoClient.generateTeal(p)
	if err != nil {
		log.Printf("Error compiling teal: %v\n", err)
		http.Error(w, err.Error(), 500)
        return
	}

    genID, genHash, minFee, firstValidRound, lastValidRound, err := algoClient.getParams()
	if err != nil {
		log.Printf("Error getting suggested parameters: %v\n", err)
		http.Error(w, err.Error(), 500)
        return
	}

    var amount uint64 = 400000
    tx1, err := transaction.MakePaymentTxnWithFlatFee(p.CreatorAddress, addr, minFee, amount, firstValidRound, lastValidRound, nil, "", genID, genHash)
    if err != nil {
        log.Printf("Error creating transaction: %v\n\n", err)
        http.Error(w, err.Error(), 400)
        return
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
    if err := json.NewEncoder(w).Encode(tx1); err != nil {
        log.Printf("Error sending back response: %s\n", err)
    }
}


func (c *AlgoClient) opTin(index uint64, addr string, account string, lsig types.LogicSig, sk1 ed25519.PrivateKey) (err error) {
    genID, genHash, minFee, firstValidRound, lastValidRound, err := c.getParams()
	if err != nil {
		return err
	}
    genHash64 := base64.StdEncoding.EncodeToString(genHash)
    txn1, err := transaction.MakeAssetAcceptanceTxn(addr, 0, firstValidRound, lastValidRound, nil, genID, genHash64, index)
    if err != nil {
        return err
    }
    txn2, err := transaction.MakePaymentTxnWithFlatFee(account, addr, minFee, 1000, firstValidRound, lastValidRound, nil, "", genID, genHash)
    if err != nil {
        return err
    }
    gid, err := crypto.ComputeGroupID([]types.Transaction{txn1, txn2})
    if err != nil {
        return err
    }
    txn1.Group = gid
    txn2.Group = gid

    _, stx1, err := crypto.SignLogicsigTransaction(lsig, txn1)
    if err != nil {
        return err
    }
    _, stx2, err := crypto.SignTransaction(sk1, txn2)
    if err != nil {
        return err
    }
    var signedGroup []byte
    signedGroup = append(signedGroup, stx1...)
    signedGroup = append(signedGroup, stx2...)
    _, err = c.c.SendRawTransaction(signedGroup).Do(context.Background())
    if err != nil {
        return err
    }
    return nil
}

func fundEscrow(w http.ResponseWriter, r *http.Request) {
    var p Choice 
    var txGroup TransactionGroup
    var algoClient AlgoClient

    account1 := getAddress(mnemonic1)
	sk1, err := mnemonic.ToPrivateKey(mnemonic1)
    decoder := json.NewDecoder(r.Body)
    err = decoder.Decode(&p)
    if err != nil {
        log.Printf("Error decoding body: %v\n", err)
        http.Error(w, err.Error(), 400)
        return
    }
    algoClient.c, err = algod.MakeClient(algodAddress, algodToken)
    if err != nil {
        log.Printf("failed to make algod client: %v\n\n", err) 
        http.Error(w, err.Error(), 500)
        return
    }

    genID, genHash, minFee, firstValidRound, lastValidRound, err := algoClient.getParams()
	if err != nil {
		log.Printf("Error getting suggested parameters: %v\n", err)
		http.Error(w, err.Error(), 500)
        return
	}
    genHash64 := base64.StdEncoding.EncodeToString(genHash)
    lsig, addr, err := algoClient.generateTeal(p)
	if err != nil {
		return
	}
    txn1, err := transaction.MakeAssetAcceptanceTxn(addr, 0, firstValidRound, lastValidRound, nil, genID, genHash64, p.AssetId)
    if err != nil {
        log.Printf("Error Making acceptance txn for %s asset number %d: %v\n", addr, p.AssetId, err)
        http.Error(w, err.Error(), 400)
        return
    }
    txn2, err := transaction.MakeAssetTransferTxnWithFlatFee(p.CreatorAddress, addr, "", p.AssetAmount, minFee, firstValidRound, lastValidRound, nil, genID, genHash64, p.AssetId)
    if err != nil {
        log.Printf("Error making asset transfer of asset %d from %s to %s: %v\n", p.AssetId, p.CreatorAddress, addr, err)
        http.Error(w, err.Error(), 400)
        return
    }
    gid, err := crypto.ComputeGroupID([]types.Transaction{txn1, txn2})
    if err != nil {
        log.Printf("Error computing GroupID: %v\n", err)
        http.Error(w, err.Error(), 400)
        return
    }
    txn1.Group = gid
    txn2.Group = gid
    _, stx1, err := crypto.SignLogicsigTransaction(lsig, txn1)
    if err != nil {
        log.Printf("Error signing with logic: %v\n", err)
        http.Error(w, err.Error(), 400)
        return
    }
    txGroup.FirstTx = stx1 
    txGroup.SecondTx = txn2

    if p.PaymentAssetId != 0 {
        err = algoClient.opTin(p.PaymentAssetId, addr, account1, lsig, sk1)
        if err != nil {
            log.Printf("Failed to opt-in asset %d: %v\n", p.PaymentAssetId, err)
            http.Error(w, err.Error(), 400)
            return
    }
    }
    if p.SecondPaymentAssetId != 0 {
        err = algoClient.opTin(p.SecondPaymentAssetId, addr, account1, lsig, sk1)
        if err != nil {
            log.Printf("Failed to opt-in asset %d: %v\n", p.SecondPaymentAssetId, err)
            http.Error(w, err.Error(), 400)
            return
        }
    }
    
    if err := json.NewEncoder(w).Encode(txGroup); err != nil {
        log.Printf("Error sending back response: %s", err)
    }
}
