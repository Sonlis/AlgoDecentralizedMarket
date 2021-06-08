<template>
<div>
  <h1>Algorand Decentralized Market Place</h1>
  <div class="Connection"> 
    <p>First, connect to myAlgo</p>
    <button v-on:click="connectToMyAlgo()" ><span>Connect to myAlgo</span></button>
  </div>
   <div id="create">
     <h3>You can create an asset here:</h3>
  <div class="form-group">
            <label for="assetName">Asset Name</label>
            <input type="assetName" class="form-control" id="assetName" placeholder="Ouis"
                v-model="assetName">
            <label for="assetUnitName">Asset Unit Name</label>
            <input type="assetUnitName" class="form-control" id="assetUnitName" placeholder="Oui"
                v-model="assetUnitName">
            <label for="assetTotal">Asset Total</label>
            <input type="assetTotal" class="form-control" id="assetTotal" placeholder=1000
                v-model.number="assetTotal">
            <label for="assetURL">Asset URL</label>
            <input type="assetURL" class="form-control" id="assetURL" placeholder="Oui"
                v-model="assetURL">
  
  <button v-on:click="createAsset()" ><span>Create asset</span></button>
  </div>
    <h3>You can opt-in an asset here:</h3> 
    <p>Opt-in asset number <input type="number" class="form-control" id="optinnumber" v-model.number="optinnumber"> <button v-on:click="optin(optinnumber)">opt-in</button></p>
  <div class="form-group">
    <h3>Sell an asset</h3>
    <p>Which asset to sell? </p> 
    <div class="form-check" v-for="(asset) in assets.assets" :key="asset['asset-id']">
    <label class="form-check-label" v-bind:for="asset['asset-id']" :name="asset['asset-id']">{{asset["asset-id"]}}</label> 
    <input class="form-check-input" type="radio" v-model.number="assetID" :name="asset['asset-id']" :id="asset['asset-id']" :value="asset['asset-id']">
    </div>
    <label for="assetAmount">Amount of asset to sell</label>
            <input type="assetAmount" class="form-control" id="assetAmount" placeholder=0
                v-model.number="assetAmount">
            <label for="paymentAssetID">First asset ID to accept for trade (payment asset)</label>
            <input type="number" class="form-control" id="paymentAssetID" placeholder=1000
                v-model.number="paymentAssetID">
            <label for="paymentAssetAmount">First payment asset amount</label>
            <input type="number" class="form-control" id="paymentAssetAmount" 
                v-model.number="paymentAssetAmount">
            <label for="secondPaymentAssetID">Second asset ID to accept for trade (paymentasset)</label>
            <input type="number" class="form-control" id="secondPaymentAssetID" 
                v-model.number="secondPaymentAssetID">
            <label for="secondPaymentAssetAmount">Second payment asset Amount</label>
            <input type="number" class="form-control" id="secondPaymentAssetAmount" placeholder=1000
                v-model.number="secondPaymentAssetAmount">
            <label for="algoAmount">MicroAlgo needed to buy the asset</label>
            <input type="number" class="form-control" id="algoAmount" placeholder="1000"
                v-model.number="algoAmount">
            
  <button v-on:click="createEscrow()" ><span>Create Escrow</span></button>
    </div>
    <h3> Which asset to buy ?</h3>
    <div><button v-on:click="getSellings()" ><span>Get assets being sold</span></button>
    <div class="form-check" v-for="(asset) in sellings.assets" :key="asset['asset-id']">
    <p>Asset {{asset['asset-id']}}</p>
    <button :name="asset['fassetid']" :id="asset['fassetid']" v-on:click="buy(asset['fassetid'], asset['fassetamount'], asset['address'], asset['asset-id'])">Buy with {{asset["fassetamount"]}} of asset {{asset["fassetid"]}}</button>
    <button :name="asset['sassetid']" :id="asset['sassetid']" v-on:click="buy(asset['sassetid'], asset['sassetamount'], asset['address'], asset['asset-id'])">Buy with {{asset['sassetamount']}} of asset {{asset["sassetid"]}}</button>
    <button :name="asset['algoamount']" :id="asset['algoamount']" v-on:click="buywithalgos(asset['algoamount'], asset['address'], asset['asset-id'])">Buy with {{asset['algoamount']}} Algos</button>
    </div>
    </div>
    </div>
    <h3>Get your Escrow accounts</h3>
    <button v-on:click="lookupSellings()">Get your escrow accounts</button>
    <li v-for="asset in escrows" :key="asset['address']"><p>Address: {{asset['address']}}, Algo Amount: {{asset['amount']}}, <p v-for="assets in asset.assets" :key="assets['asset-id']"> {{assets['amount']}} of asset {{assets['asset-id']}}</p></li>
    <p>Withdraw <input type="number" class="form-control" id="algosamount" placeholder=1000 v-model.number="algosamount"> Algos from <input class="form-control" id="paymentAssetID" placeholder=1000 v-model="escrowaddress"><button v-on:click="withdrawalgo(algosamount, escrowaddress)">Withdraw</button></p>
    <p>Withdraw <input type="number" class="form-control" id="withdrawassetamount" placeholder=1000 v-model.number="withdrawassetamount"> of asset number <input type="number" class="form-control" id="withdrawassetid" placeholder=1000 v-model.number="withdrawassetid"> from <input class="form-control" id="escrowaddress2" placeholder=1000 v-model="escrowaddress"><button v-on:click="withdrawasset(withdrawassetid, withdrawassetamount, escrowaddress)">Withdraw</button></p>
</div>
</template>

<script> 
import MyAlgo from '@randlabs/myalgo-connect';
import algosdk from 'algosdk';
import axios from 'axios';

const algodClient = new algosdk.Algodv2('aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa', 'http://localhost', 4001);
const myAlgoWallet = new MyAlgo();

export default {
    name: 'Charts',
    data(){
      return {
        form: {
          accountid: ''
        },
        assets: [],
        paymentID: 0,
        sellings: [],
        escrows: [],
        assetName: '',
        assetUnitName: '', 
        assetURL: '',
        assetTotal: '',
        assetID: 0,
        assetAmount: 0,
        paymentAssetID: 0,
        paymentAssetAmount: 0,
        secondPaymentAssetID: 0,
        secondPaymentAssetAmount: 0,
        algoAmount: 0,
        escrowaddress: '',
        withdrawassetid: 0,
        withdrawassetamount: 0,
        algosamount: 0,
        optinnumber: 0

      }
    },
    methods: {
        connectToMyAlgo: async function () {
            try {
                const accounts = await myAlgoWallet.connect();
                console.log(accounts)

                this.addresses = accounts.map(account => account.address);
    
                } catch (err) {
                console.error(err);
            }
          let data = {"accountid": this.addresses[0]}
          console.log(this.addresses)
          axios({method: "POST", url: "http://localhost:8081/lookup", data: data, headers: {"Content-Type": "application/json"}}).then(response => {
            this.assets = response.data;
            console.log(this.assets);
          });
        },
        createAsset: async function() {
             let txn = await algodClient.getTransactionParams().do();
        this.txn = {
          ...txn,
          type: 'acfg',
          from: this.addresses[0],
          assetName: this.assetName,
          assetUnitName: this.assetUnitName,
          assetTotal: this.assetTotal,
          assetURL: this.assetURL,
          assetManager: this.addresses[0],
          assetReserve: this.addresses[0],
    };
        console.log(this.txn)
        let signedTxn = await myAlgoWallet.signTransaction(this.txn);
        let tosend = Array.from(signedTxn.blob);
        let data = {"raw": tosend}
        console.log(this.data)
        axios({method: "POST", url: "http://localhost:8081/createAsset", data: data, headers: {"Content-Type": "application/json"}})
        },

        createEscrow: async function() {
        const requestOptions = {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({"assetid": this.assetID, "assetamount": this.assetAmount, "paymentassetid": this.paymentAssetID, "paymentassetamount": this.paymentAssetAmount, "secondpaymentassetid": this.secondPaymentAssetID, "secondpaymentassetamount": this.secondPaymentAssetAmount, "algoamount": this.algoAmount, "creatoraddress": this.addresses[0]}),
        }
        console.log(requestOptions.body)
        const response = await fetch("http://localhost:8081/createEscrow", requestOptions);
        let answer = await response.json();
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
        getSellings: async function() {
          const response = await fetch("http://localhost:8081/getSellings");
          this.sellings = await response.json();
          console.log(this.sellings)
          },
        buy: async function(paymentasset, paymentassetamount, address, assetid) {
          const requestOptions = {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ "tobuy": assetid, "topay": paymentasset, "amount": paymentassetamount, "address": address, "sender": this.addresses[0]}),
          }
          console.log(requestOptions.body)
          const response = await fetch("http://localhost:8081/buy", requestOptions);
          const buytx = await response.json();
          console.log(buytx)
          let buytx1 = buytx["firsttx"]
          let stxn2 = buytx["secondtx"]
          let buytx3 = buytx["thirdtx"]
          let txn1 = {
            fee: 1000,
            flatFee: true,
            type: 'axfer',
            assetIndex: buytx1.XferAsset,
            from: algosdk.encodeAddress(buytx1.Sender), //btoa(String.fromCharCode.apply(null, new Uint8Array(buytx1.Sender))),
            to: algosdk.encodeAddress(buytx1.AssetReceiver),
            amount: buytx1.AssetAmount,
            genesisID: buytx1.GenesisID,
            firstRound: buytx1.FirstValid,
            genesisHash: buytx1.GenesisHash,
            lastRound: buytx1.LastValid,
            group: buytx1.Group
          }
          console.log(txn1)
          let txn3 = {
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
          console.log(txn3)
          let stxn1 = await myAlgoWallet.signTransaction(txn1);
          let stxn3 = await myAlgoWallet.signTransaction(txn3);
          let asciiString = atob(stxn2);
          let oui = new Uint8Array([...asciiString].map(char => char.charCodeAt(0)));
          this.tosend = [stxn1.blob, oui, stxn3.blob]
          console.log(this.tosend)
          try {
            await algodClient.sendRawTransaction(this.tosend).do();
          } catch(exception) {
            console.log(exception)
          }
        },
        buywithalgos: async function(algoamount, address, assetid) {
          const requestOptions = {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ "tobuy": assetid, "algoamount": algoamount, "address": address, "sender": this.addresses[0]}),
          }
          console.log(requestOptions.body)
          const response = await fetch("http://localhost:8081/buy", requestOptions);
          const buytx = await response.json();
          console.log(buytx)
          let buytx1 = buytx["firsttx"]
          let stxn2 = buytx["secondtx"]
          let buytx3 = buytx["thirdtx"]
          let txn1 = {
            fee: 1000,
            flatFee: true,
            type: 'axfer',
            assetIndex: buytx1.XferAsset,
            from: algosdk.encodeAddress(buytx1.Sender), //btoa(String.fromCharCode.apply(null, new Uint8Array(buytx1.Sender))),
            to: algosdk.encodeAddress(buytx1.AssetReceiver),
            amount: buytx1.AssetAmount,
            genesisID: buytx1.GenesisID,
            firstRound: buytx1.FirstValid,
            genesisHash: buytx1.GenesisHash,
            lastRound: buytx1.LastValid,
            group: buytx1.Group
          }
          console.log(txn1)
          let txn3 = {
            fee: 1000,
            flatFee: true,
            type: 'pay',
            from: algosdk.encodeAddress(buytx3.Sender),
            to: algosdk.encodeAddress(buytx3.Receiver),
            amount: algoamount,
            genesisID: buytx3.GenesisID,
            firstRound: buytx3.FirstValid,
            genesisHash: buytx3.GenesisHash,
            lastRound: buytx3.LastValid,
            group: buytx3.Group
          }
          console.log(txn3)
          let stxn1 = await myAlgoWallet.signTransaction(txn1);
          let stxn3 = await myAlgoWallet.signTransaction(txn3);
          let asciiString = atob(stxn2);
          stxn2 = new Uint8Array([...asciiString].map(char => char.charCodeAt(0)));
          this.tosend = [stxn1.blob, stxn2, stxn3.blob]
          console.log(this.tosend)
          try {
            await algodClient.sendRawTransaction(this.tosend).do();
          } catch(exception) {
            console.log(exception)
          }
        },
        lookupSellings: async function() {
          const requestOptions = {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({"accountid": this.addresses[0]}),
          }
          console.log(requestOptions.body)
          const response = await fetch("http://localhost:8081/lookupSellings", requestOptions);
          this.escrows = await response.json();
          console.log(this.escrows)
        },
        withdrawalgo: async function(algosamount, escrowaddress) {
        const requestOptions = {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({"algosamount": algosamount, "escrowaddress": escrowaddress, "to": this.addresses[0]}),
          }
          console.log(requestOptions.body)
          const response = await fetch("http://localhost:8081/withdrawAlgos", requestOptions);
          let algos = await response.json();
          let stxn1 = algos['firsttx']
          let params = algos['secondtx']
          let txn2 = {
            fee: 1000,
            flatFee: true,
            type: 'pay',
            from: algosdk.encodeAddress(params.Sender), //btoa(String.fromCharCode.apply(null, new Uint8Array(buytx1.Sender))),
            to: algosdk.encodeAddress(params.Receiver),
            amount: params.Amount,
            genesisID: params.GenesisID,
            firstRound: params.FirstValid,
            genesisHash: params.GenesisHash,
            lastRound: params.LastValid,
            group: params.Group
          }
          let stxn2 = await myAlgoWallet.signTransaction(txn2);
          let asciiString = atob(stxn1);
          stxn1 = new Uint8Array([...asciiString].map(char => char.charCodeAt(0)));
          this.tosend = [stxn1, stxn2.blob]
          this.download_txns("grouped.txns", this.tosend)
          console.log(this.tosend)
          
          try {
            await algodClient.sendRawTransaction(this.tosend).do();
          } catch(exception) {
            console.log(exception)
          }
        },
        withdrawasset: async function(assetid, assetamount, escrowaddress) {
          const requestOptions = {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({"assetamount": assetamount, "assetid": assetid, "escrowaddress": escrowaddress, "to": this.addresses[0]}),
          }
          console.log(requestOptions.body)
          const response = await fetch("http://localhost:8081/withdrawAssets", requestOptions);
          let assets = await response.json();
          console.log(assets)
          let stxn1 = assets['firsttx']
          let params = assets['secondtx']
          let txn2 = {
            fee: 1000,
            flatFee: true,
            type: 'pay',
            from: algosdk.encodeAddress(params.Sender), 
            to: algosdk.encodeAddress(params.Receiver),
            amount: params.Amount,
            genesisID: params.GenesisID,
            firstRound: params.FirstValid,
            genesisHash: params.GenesisHash,
            lastRound: params.LastValid,
            group: params.Group
          }
          let stxn2 = await myAlgoWallet.signTransaction(txn2);
          let asciiString = atob(stxn1);
          stxn1 = new Uint8Array([...asciiString].map(char => char.charCodeAt(0)));
          let tosend = [stxn1, stxn2.blob]
          try {
            await algodClient.sendRawTransaction(tosend).do();
          } catch(exception) {
            console.log(exception)
          }
      },
        optin: async function(optinnumber) {
          let txn = await algodClient.getTransactionParams().do();
      
          txn = {
          ...txn,
          fee: 1000,
          flatFee: true,
          type: 'axfer',
          assetIndex: optinnumber,
          from: this.addresses[0],
          to:  this.addresses[0],
          amount: 0,
        };
  
          let signedTxn = await myAlgoWallet.signTransaction(txn);
          console.log(signedTxn.txID);
  
          await algodClient.sendRawTransaction(signedTxn.blob).do();
      }
      
    }
    }
        
        
</script>