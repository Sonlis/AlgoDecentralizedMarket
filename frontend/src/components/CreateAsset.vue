<template>
<div>
    <div>
        <h4>Create an asset</h4>
        <label for="assetName">Asset Name</label>
        <input type="assetName"  id="assetName" placeholder="Ouis"
            v-model="form.assetName">
        <div class="error" v-if="$v.form.assetName.$invalid && submitStatus === 'ERROR'">Please enter an asset name</div>
        <label for="assetUnitName">Asset Unit Name</label>
        <input type="assetUnitName"  id="assetUnitName" placeholder="Oui"
            v-model="form.assetUnitName">
        <div class="error" v-if="$v.form.assetUnitName.$invalid && submitStatus === 'ERROR'">Please enter an asset unit name</div>
        <label for="assetTotal">Asset Total</label>
        <input type="assetTotal"  id="assetTotal" placeholder=1000
            v-model.number="form.assetTotal">
        <div class="error" v-if="$v.form.assetTotal.$invalid && submitStatus === 'ERROR'">Please enter an asset total</div>
        <label for="assetURL">Asset URL</label>
        <input type="assetURL"  id="assetURL" placeholder="http://oui.com"
            v-model="form.assetURL">
        <div class="error" v-if="$v.form.assetURL.$invalid && submitStatus === 'ERROR'">Please enter an asset URL</div>
            <br />
        <button v-on:click="createAsset()" ><span>Create asset</span></button>
    </div>
    <br />
    <div> 
        <h4>You can opt-in an asset here:</h4> 
        <p>Opt-in asset number <input type="number" id="optinnumber" v-model.number="optinnumber"> <button v-on:click="optin(optinnumber)">opt-in</button></p>
    </div>
</div>
</template>

<script> 
import { required, minValue } from 'vuelidate/lib/validators'
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
            submitStatus: null
        }
    },
    validations: {
        form: {
            assetName: { required },
            assetUnitName: { required },
            assetTotal: {required, minValue: minValue(1) },
            assetURL: { required },
        },
    },
    methods: {
        createAsset: async function() {
            this.$v.$touch()
            if (this.$v.$invalid) {
                this.submitStatus = 'ERROR'
            } else {
                this.submitStatus = 'PENDING'
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
            let signedTxn = myAlgoWallet.signTransaction(txn);
            await algodClient.sendRawTransaction(signedTxn.blob).do();
            this.submitStatus = 'OK'
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
          from: this.addrToUse,
          to:  this.addrToUse,
          amount: 0,
        };
          let signedTxn = await myAlgoWallet.signTransaction(txn);

          await algodClient.sendRawTransaction(signedTxn.blob).do();
      }
    }
}
</script>

<style scoped>
#form-control {
    width: 15%;
}
</style>