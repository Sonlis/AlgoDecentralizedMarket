<template>
<div>
  <h2>Algorand Decentralized Market Place</h2>
  <div class="Connection"> 
    <my-algo-connection v-on:clicked="saveAddress"></my-algo-connection>
  </div>

  <div class="tabs">
    <b-tabs content-class="mt-3" lazy>
      <b-tab title="Helper functions">
        <create-asset :addrToUse="addrToUse"></create-asset>
      </b-tab> 
      <b-tab title="Sell an asset">
        <sell-asset :addrToUse="addrToUse" v-on:returnSellParameters="createEscrow"></sell-asset>
      </b-tab>
      <b-tab title="Buy an asset">
        <buy-asset v-on:returnBuyParameters="buy"></buy-asset>
      </b-tab>
      <b-tab title="Withdraw from your selling accounts">
        <user-escrows :addrToUse="addrToUse" v-on:returnWithdrawParameters="withdraw"></user-escrows>
      </b-tab>
    </b-tabs>
  </div>
</div>
</template>

<script> 
import MyAlgo from '@randlabs/myalgo-connect';
import algosdk from 'algosdk';
import MyAlgoConnection from './myAlgoConnection.vue';
import CreateAsset from './CreateAsset.vue';
import SellAsset from './SellAsset.vue'
import BuyAsset from './BuyAsset.vue'
import UserEscrows from './UserEscrows.vue';

const algodClient = new algosdk.Algodv2('aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa', 'http://localhost', 4001);
const myAlgoWallet = new MyAlgo();

export default {
    name: 'AlgoDeMarket',
    components: {
      MyAlgoConnection,
        CreateAsset,
        SellAsset,
        BuyAsset,
        UserEscrows
    },
    data(){
      return {
        addrToUse: ''
      }
    },
    methods: {
        saveAddress: function(addr) {
          this.addrToUse = addr
        },
        createEscrow: async function(sellForm) {
        const requestOptions = {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({"assetid": sellForm.assetID, "assetamount": sellForm.assetAmount, "paymentassetid": sellForm.paymentAssetID, "paymentassetamount": sellForm.paymentAssetAmount, "secondpaymentassetid": sellForm.secondPaymentAssetID, "secondpaymentassetamount": sellForm.secondPaymentAssetAmount, "algoamount": sellForm.algoAmount, "creatoraddress": this.addrToUse}),
        }
        console.log(requestOptions.body)
        const response = await fetch("http://localhost:8081/createEscrow", requestOptions);
        let answer = await response.json();
        console.log(answer)
        let stxn1 = answer["firsttx"]
        let tx2 = answer["secondtx"]
        let txn2 = {
          fee: 1000,
          flatFee: true,
          type: 'axfer',
          assetIndex: tx2.XferAsset,
          from: algosdk.encodeAddress(tx2.Sender),
          to: algosdk.encodeAddress(tx2.AssetReceiver),
          amount: tx2.AssetAmount,
          genesisID: tx2.GenesisID,
          firstRound: tx2.FirstValid,
          genesisHash: tx2.GenesisHash,
          lastRound: tx2.LastValid,
          group: tx2.Group
        }
        let stxn2 = await myAlgoWallet.signTransaction(txn2);
        let asciiString = atob(stxn1);
        stxn1 = new Uint8Array([...asciiString].map(char => char.charCodeAt(0)));
        this.tosend = [stxn1, stxn2.blob]
        console.log(this.tosend)
        try {
          await algodClient.sendRawTransaction(this.tosend).do();
        } catch(exception) {
          console.log(exception)
        }
        
        },
        buy: async function(buyForm) {
          if (buyForm.algoAmount === 0){
            this.requestOptions = {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ "tobuy": buyForm.assetID, "topay": buyForm.paymentAssetID, "amount": buyForm.paymentAssetAmount, "address": buyForm.address, "sender": this.addrToUse}),
          }
          } else {
            this.requestOptions = {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ "tobuy": buyForm.assetID, "algoamount": buyForm.algoAmount, "address": buyForm.address, "sender": this.addrToUse}),
          }
          }
          console.log(this.requestOptions.body)
          const response = await fetch("http://localhost:8081/buy", this.requestOptions);
          const buytx = await response.json();
          console.log(buytx)
          let buytx1 = buytx["firsttx"]
          let stxn2 = buytx["secondtx"]
          let buytx3 = buytx["thirdtx"]
          let buytx4 = buytx["forthtx"]
          let txn1 = {
            fee: 1000,
            flatFee: true,
            type: 'axfer',
            assetIndex: buytx1.XferAsset,
            from: algosdk.encodeAddress(buytx1.Sender),
            to: algosdk.encodeAddress(buytx1.AssetReceiver),
            amount: buytx1.AssetAmount,
            genesisID: buytx1.GenesisID,
            firstRound: buytx1.FirstValid,
            genesisHash: buytx1.GenesisHash,
            lastRound: buytx1.LastValid,
            group: buytx1.Group
          }
          if (buyForm.algoAmount === 0) {
            this.txn3 = {
            fee: 1000,
            flatFee: true,
            type: 'axfer',
            assetIndex: buytx3.XferAsset,
            from: algosdk.encodeAddress(buytx3.Sender),
            to: algosdk.encodeAddress(buytx3.AssetReceiver),
            amount: buytx3.AssetAmount,
            genesisID: buytx3.GenesisID,
            firstRound: buytx3.FirstValid,
            genesisHash: buytx3.GenesisHash,
            lastRound: buytx3.LastValid,
            group: buytx3.Group
          }
          } else {
            this.txn3 = {
            fee: 1000,
            flatFee: true,
            type: 'pay',
            from: algosdk.encodeAddress(buytx3.Sender),
            to: algosdk.encodeAddress(buytx3.Receiver),
            amount: buytx3.Amount,
            genesisID: buytx3.GenesisID,
            firstRound: buytx3.FirstValid,
            genesisHash: buytx3.GenesisHash,
            lastRound: buytx3.LastValid,
            group: buytx3.Group
          }
          }
          let txn4 = {
            fee: 1000,
            flatFee: true,
            type: 'pay',
            from: algosdk.encodeAddress(buytx4.Sender),
            to: algosdk.encodeAddress(buytx4.Receiver),
            amount: buytx4.Amount,
            genesisID: buytx4.GenesisID,
            firstRound: buytx4.FirstValid,
            genesisHash: buytx4.GenesisHash,
            lastRound: buytx4.LastValid,
            group: buytx4.Group
          }
          let stxn1 = await myAlgoWallet.signTransaction(txn1);
          let stxn3 = await myAlgoWallet.signTransaction(this.txn3);
          let stxn4 = await myAlgoWallet.signTransaction(txn4);
          let asciiString = atob(stxn2);
          let oui = new Uint8Array([...asciiString].map(char => char.charCodeAt(0)));
          this.tosend = [stxn1.blob, oui, stxn3.blob, stxn4.blob]
          console.log(this.tosend)
          try {
            await algodClient.sendRawTransaction(this.tosend).do();
          } catch(exception) {
            console.log(exception)
          }
        },
        withdraw: async function(withdrawForm) {
          console.log(withdrawForm)
          if (withdrawForm.algosAmount === 0) {
            console.log("in asset transfer")
            this.requestOptions = {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({"assetamount": withdrawForm.assetAmount, "assetid": withdrawForm.assetID, "escrowaddress": withdrawForm.escrowAddress, "to": this.addrToUse}),
          }
            console.log(this.requestOptions.body)
            const response = await fetch("http://localhost:8081/withdraw", this.requestOptions);
            let withdraw = await response.json();
            this.stxn1 = withdraw['firsttx']
            this.tx2 = withdraw['secondtx']
            console.log(withdraw)
          } else {
            this.requestOptions = {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({"algosamount": withdrawForm.algosAmount, "escrowaddress": withdrawForm.escrowAddress, "to": this.addrToUse})
            }
            console.log(this.requestOptions.body)
            const response = await fetch("http://localhost:8081/withdraw", this.requestOptions);
            let withdraw = await response.json();
            this.stxn1 = withdraw['firsttx']
            this.tx2 = withdraw['secondtx']
            console.log(this.tx2)
          }
          let txn2 = {
            fee: 1000,
            flatFee: true,
            type: 'pay',
            from: algosdk.encodeAddress(this.tx2.Sender), 
            to: algosdk.encodeAddress(this.tx2.Receiver),
            amount: this.tx2.Amount,
            genesisID: this.tx2.GenesisID,
            firstRound: this.tx2.FirstValid,
            genesisHash: this.tx2.GenesisHash,
            lastRound: this.tx2.LastValid,
            group: this.tx2.Group
          }
          let stxn2 = await myAlgoWallet.signTransaction(txn2);
          let asciiString = atob(this.stxn1);
          this.stxn1 = new Uint8Array([...asciiString].map(char => char.charCodeAt(0)));
          let tosend = [this.stxn1, stxn2.blob]
          this.download_txns("grouped.txns", tosend)
          try {
            await algodClient.sendRawTransaction(tosend).do();
          } catch(exception) {
            console.log(exception)
          }
          },
        download_txns: function (name, txns)  {
    let b = new Uint8Array(0);
    for(const txn in txns){
        b = this.concatTypedArrays(b, txns[txn])
    }
    var blob = new Blob([b], {type: "application/octet-stream"});

    var link = document.createElement('a');
    link.href = window.URL.createObjectURL(blob);
    link.download = name;
    link.click();
},

concatTypedArrays: function (a, b) { // a, b TypedArray of same type
    var c = new (a.constructor)(a.length + b.length);
    c.set(a, 0);
    c.set(b, a.length);
    return c;
}
    }   
    }        
</script>

<style> 

  body {
    padding-left: 2%;
    padding-right: 2%;
    background-color: rgb(233, 233, 233)
  }
  input {
    width: 25%;
  }
  .error {
    color: red;
  }
  
</style>
