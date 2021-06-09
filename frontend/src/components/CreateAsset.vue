<template>
<div>
    <div class="form-group">
        <label for="assetName">Asset Name</label>
        <input type="assetName" class="form-control" id="assetName" placeholder="Ouis"
            v-model="form.assetName">
        <label for="assetUnitName">Asset Unit Name</label>
        <input type="assetUnitName" class="form-control" id="assetUnitName" placeholder="Oui"
            v-model="form.assetUnitName">
        <label for="assetTotal">Asset Total</label>
        <input type="assetTotal" class="form-control" id="assetTotal" placeholder=1000
            v-model.number="form.assetTotal">
        <label for="assetURL">Asset URL</label>
        <input type="assetURL" class="form-control" id="assetURL" placeholder="http://oui.com"
            v-model="form.assetURL">
        <button v-on:click="createAsset()" ><span>Create asset</span></button>
    </div>
    <div> 
        <h3>You can opt-in an asset here:</h3> 
        <p>Opt-in asset number <input type="number" class="form-control" id="optinnumber" v-model.number="optinnumber"> <button v-on:click="optin(optinnumber)">opt-in</button></p>
    </div>
</div>
</template>

<script> 
import MyAlgo from '@randlabs/myalgo-connect';
import algosdk from 'algosdk';
const myAlgoWallet = new MyAlgo();
const algodClient = new algosdk.Algodv2('aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa', 'http://localhost', 4001);

export default {
    name: "CreateAsset",
    props: ['addrToUse'],
    data() {
        return {
            form: {
                assetName: '',
                assetUnitName: '',
                assetTotal: 0,
                assetURL: ''
            },
            optinnumber: 0,
        }
    },
    methods: {
        createAsset: async function() {
            console.log(this.addrToUse)
            let txn = await algodClient.getTransactionParams().do();
            txn = {
            ...txn,
            type: 'acfg',
            from: this.addrToUse,
            assetName: this.form.assetName,
            assetUnitName: this.form.assetUnitName,
            assetTotal: this.form.assetTotal,
            assetURL: this.form.assetURL,
            assetManager: this.addrToUse,
            assetReserve: this.addrToUse,
    };
        myAlgoWallet.signTransaction(txn);
        },
        optin: async function(optinnumber) {
          let txn = await algodClient.getTransactionParams().do();
          txn = {
          ...txn,
          fee: 1000,
          flatFee: true,
          type: 'axfer',
          assetIndex: optinnumber,
          from: this.addrToUse,
          to:  this.addrToUse,
          amount: 0,
        };
          let signedTxn = await myAlgoWallet.signTransaction(txn);
          console.log(signedTxn.txID);
  
          await algodClient.sendRawTransaction(signedTxn.blob).do();
      }
    }
}
</script>