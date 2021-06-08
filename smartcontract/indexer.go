package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
)



func lookupAssets(w http.ResponseWriter, r *http.Request) {
	indexerClient, err := indexer.MakeClient(indexerAddress, indexerToken)
	if err != nil {
		return
	}
	var accountID AccountID
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&accountID)
	if err != nil {
		log.Println("Error decoding:", err)
	}
	_, result, err := indexerClient.LookupAccountByID(accountID.AccountID).Do(context.Background())
    var assets Assets
    assetslist, err := json.Marshal(result)
	if err != nil {
		log.Println("Error marshaling:", err)
	}
    err = json.Unmarshal(assetslist, &assets)
	if err != nil {
		log.Println("Error Unmarshaling:", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(assets); err != nil {
        log.Printf("Error sending back response: %s", err)
    }
}

func lookupEscrowAssets(w http.ResponseWriter, r *http.Request) {
	indexerClient, err := indexer.MakeClient(indexerAddress, indexerToken)
	if err != nil {
		log.Printf("Error creating indexer")
	}
	var toReturn Assets
	for i := range(sellings) {
		_, result, err := indexerClient.LookupAccountByID(sellings[i].Address).Do(context.Background())
		var assets Assets
    	assetslist, err := json.Marshal(result)
		if err != nil {
			log.Println("Error marshaling:", err)
		}
    	err = json.Unmarshal(assetslist, &assets)
		if err != nil {
			log.Println("Error Unmarshaling:", err)
		}	
		for j := range(assets.Assets) {
			if assets.Assets[j].AssetId == sellings[i].Asset {
				assets.Assets[j].Address = sellings[i].Address
				assets.Assets[j].FirstAsset = sellings[i].FirstAsset
				assets.Assets[j].FAmount = sellings[i].FAmount
				assets.Assets[j].SecondAsset = sellings[i].SecondAsset
				assets.Assets[j].SAmount = sellings[i].SAmount
				assets.Assets[j].AlgoAmount = sellings[i].AlgoAmount
				toReturn.Assets = append(toReturn.Assets, assets.Assets[j])
			}
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(toReturn); err != nil {
        log.Printf("Error sending back response: %s", err)
    }
}

func lookupSellings(w http.ResponseWriter, r *http.Request) {
	var (
		toReturn []Assets
		toCheck []Selling
	 	accountID AccountID
	)
	indexerClient, err := indexer.MakeClient(indexerAddress, indexerToken)
	if err != nil {
		log.Printf("Error creating indexer: %s", err)
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&accountID)
	if err != nil {
		log.Println("Error decoding:", err)
	}
	for i := range(sellings) {
		if sellings[i].CreatorAddress == accountID.AccountID {
			toCheck = append(toCheck, sellings[i])
		}
	}
	for i := range(toCheck) {
		_, result, err := indexerClient.LookupAccountByID(toCheck[i].Address).Do(context.Background())
		if err != nil {
			log.Printf("Error checking escrow account %s: %s", toCheck[i].Address, err)
		}
		var assets Assets
    	assetslist, err := json.Marshal(result)
		if err != nil {
			log.Println("Error marshaling:", err)
		}
    	err = json.Unmarshal(assetslist, &assets)
		if err != nil {
			log.Println("Error Unmarshaling:", err)
		}
		toReturn = append(toReturn, assets)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(toReturn); err != nil {
        log.Printf("Error sending back response: %s", err)
    }
}
