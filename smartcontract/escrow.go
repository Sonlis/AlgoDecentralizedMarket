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

func createEscrow(w http.ResponseWriter, r *http.Request) {
    var p Choice
    var pointer *uint64
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&p)
    if err != nil {
        log.Printf("Error decoding body: %v\n", err)
        http.Error(w, err.Error(), 400)
        return
    }
    account1 := getAddress(mnemonic1)
    sk1, err := mnemonic.ToPrivateKey(mnemonic1)

    var algoClient AlgoClient
    algoClient.c, err = algod.MakeClient(algodAddress, algodToken)
    if err != nil {
        log.Printf("failed to make algod client: %v\n\n", err) 
        http.Error(w, err.Error(), 500)
        return
    }
    lsig, addr, err := algoClient.generateTeal(p)
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
    var amount uint64 = 1000000
    tx1, err := transaction.MakePaymentTxnWithFlatFee(account1, addr, minFee, amount, firstValidRound, lastValidRound, nil, "", genID, genHash)
    if err != nil {
        log.Printf("Error creating transaction: %v\n\n", err)
        http.Error(w, err.Error(), 400)
        return
    }

    _, stx1, err := crypto.SignTransaction(sk1, tx1)
    if err != nil {
        log.Printf("Failed to sign transaction: %s\n", err)
        http.Error(w, err.Error(), 500)
        return
    }
    
    pendingTxID, err := algoClient.c.SendRawTransaction(stx1).Do(context.Background())
    if err != nil {
        log.Printf("Failed to send transaction: %s\n", err)
        http.Error(w, err.Error(), 500)
        return
    }
    err = algoClient.opTin(p.PaymentAssetId, addr, account1, lsig, sk1)

    if err != nil {
        log.Printf("Failed to opt-in asset %d: %v\n", p.PaymentAssetId, err)
        http.Error(w, err.Error(), 400)
        return
    }
    pointer = &p.SecondPaymentAssetId
    if pointer != nil {
        err = algoClient.opTin(p.SecondPaymentAssetId, addr, account1, lsig, sk1)
        if err != nil {
            log.Printf("Failed to opt-in asset %d: %v\n", p.SecondPaymentAssetId, err)
            http.Error(w, err.Error(), 400)
            return
        }
    }
    
    err = algoClient.waitForConfirmation(pendingTxID)
    if err != nil {
        log.Printf("Failed to Send algos to %s: %v\n", addr, err)
        http.Error(w, err.Error(), 500)
        return
    }

    txGroup, err := fundEscrow(p, addr, &algoClient, w)
    if err != nil {
        log.Printf("Error funding the account %s: %v\n", addr, err)
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
    if err := json.NewEncoder(w).Encode(txGroup); err != nil {
        log.Printf("Error sending back response: %s", err)
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
func fundEscrow(p Choice, addr string, c *AlgoClient, w http.ResponseWriter) (txGroup TransactionGroup, err error) {
    genID, genHash, minFee, firstValidRound, lastValidRound, err := c.getParams()
	if err != nil {
		log.Printf("Error getting suggested parameters: %v\n", err)
		http.Error(w, err.Error(), 500)
        return
	}
    genHash64 := base64.StdEncoding.EncodeToString(genHash)
    lsig, addr, err := c.generateTeal(p)
	if err != nil {
		return
	}
    txn1, err := transaction.MakeAssetAcceptanceTxn(addr, 0, firstValidRound, lastValidRound, nil, genID, genHash64, p.AssetId)
    if err != nil {
        return 
    }
    txn2, err := transaction.MakeAssetTransferTxnWithFlatFee(p.CreatorAddress, addr, "", p.AssetAmount, minFee, firstValidRound, lastValidRound, nil, genID, genHash64, p.AssetId)
    if err != nil {
        return
    }
    gid, err := crypto.ComputeGroupID([]types.Transaction{txn1, txn2})
    if err != nil {
        return 
    }
    txn1.Group = gid
    txn2.Group = gid
    _, stx1, err := crypto.SignLogicsigTransaction(lsig, txn1)
    if err != nil {
        return
    }
    txGroup.FirstTx = stx1 
    txGroup.SecondTx = txn2
    return txGroup, nil
}
