package main

import (
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
)
//Assets stores a slice of Asset. It also stores the amount of algo of an escrow account.
type Assets struct {
    Assets  []Asset  `json:"assets"`
	Amount  uint64   `json:"amount"`
	Address string   `json:"address"`
}
//Asset represents an AssetId and a list of tradable asset and Algos to buy this Asset.
//It also gives information about the creator of the Escrow.
type Asset struct {
    Amount uint64  			`json:"amount"`
    AssetId uint64 			`json:"asset-id"`
	Address string 			`json:"address"`
	FirstAsset  uint64		`json:"fassetid"`
	FAmount   uint64		`json:"fassetamount"`
	SecondAsset uint64		`json:"sassetid"`
	SAmount uint64			`json:"sassetamount"`
	AlgoAmount uint64   	`json:"algoamount"`
	CreatorAddress string 	`json:"creatoraddress"`
}
//AccountID stores a human readable Algorand address.
type AccountID struct {
	AccountID string  `json:"accountid"`
}

type AlgoClient struct {
	c *algod.Client
}
type CreateAsset struct {
	Signed []uint8 `json:"raw"`
}

type TransactionGroup struct {
	FirstTx []byte           		 `json:"firsttx"`
	SecondTx types.Transaction		 `json:"secondtx"`
}

type Choice struct {
	AssetId  uint64    					`json:"assetid"`
	AssetAmount   uint64    			`json:"assetamount"`
	PaymentAssetId uint64 				`json:"paymentassetid"`
	PaymentAssetAmount uint64 			`json:"paymentassetamount"`
	SecondPaymentAssetId uint64 		`json:"secondpaymentassetid"`
	SecondPaymentAssetAmount uint64 	`json:"secondpaymentassetamount"`
	CreatorAddress  string    			`json:"creatoraddress"`
	AlgoAmount  uint64   				`json:"algoamount"`
}

type Selling struct {
	Asset   uint64
	Address string  
	FirstAsset  uint64
	FAmount   uint64
	SecondAsset uint64
	SAmount uint64
	AlgoAmount uint64
	CreatorAddress string
}

type txGroup struct {
	FirstTx types.Transaction         	`json:"firsttx"`
	SecondTx []byte 		 			`json:"secondtx"`
	ThirdTx types.Transaction			`json:"thirdtx"`
}

type Tx struct {
	ToBuy uint64   `json:"tobuy"`
	ToPay uint64   `json:"topay"`
	Amount uint64  `json:"amount"`
	Sender string  `json:"sender"`
	Address string `json:"address"`
	AlgoAmount uint64 `json:"algoamount"`
}

type Withdraw struct {
	AssetID 	uint64 `json:"assetid"`
	AssetAmount uint64 `json:"assetamount"`
	Algo 	uint64 `json:"algosamount"`
	Address string `json:"escrowaddress"`
	Creator string `json:"to"`
}

var sellings []Selling