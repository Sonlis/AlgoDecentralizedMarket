<template>
    <div>
        <button v-on:click="getAssets()" ><span>Sell an Asset</span></button>
        <div class="form-check" v-for="(asset) in assets.assets" :key="asset['asset-id']">
            <label class="form-check-label" v-bind:for="asset['asset-id']" :name="asset['asset-id']">{{asset["asset-id"]}}</label> 
            <input class="form-check-input" type="radio" v-model.number="sellForm.assetID" :name="asset['asset-id']" :id="asset['asset-id']" :value="asset['asset-id']">
        </div>
        <label for="assetAmount">Amount of asset to sell</label>
        <input type="assetAmount" class="form-control" id="assetAmount" placeholder=0
            v-model.number="sellForm.assetAmount">

        <label for="paymentAssetID">First asset ID to accept for trade (payment asset)</label>
        <input type="number" class="form-control" id="paymentAssetID" placeholder=1000
            v-model.number="sellForm.paymentAssetID">

        <label for="paymentAssetAmount">First payment asset amount</label>
        <input type="number" class="form-control" id="paymentAssetAmount" 
            v-model.number="sellForm.paymentAssetAmount">

        <label for="secondPaymentAssetID">Second asset ID to accept for trade (paymentasset)</label>
        <input type="number" class="form-control" id="secondPaymentAssetID" 
            v-model.number="sellForm.secondPaymentAssetID">

        <label for="secondPaymentAssetAmount">Second payment asset Amount</label>
        <input type="number" class="form-control" id="secondPaymentAssetAmount" placeholder=1000
            v-model.number="sellForm.secondPaymentAssetAmount">

        <label for="algoAmount">MicroAlgo needed to buy the asset</label>
        <input type="number" class="form-control" id="algoAmount" placeholder="1000"
            v-model.number="sellForm.algoAmount">
        <button v-on:click="returnSellParameters()">Sell the asset</button>
    </div>
</template>


<script>


export default {
    name: 'SellAsset',
    props: ['addrToUse'],
    data(){
        return{
            assets: [''],
            sellForm: {
                assetID: 0,
                assetAmount: 0,
                paymentAssetID: 0,
                paymentAssetAmount: 0,
                secondPaymentAssetID: 0,
                secondPaymentAssetAmount: 0,
                algoAmount: 0
            }
        }
    },
    methods: {
        getAssets: async function () {
            const requestOptions = {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({"accountid": this.addrToUse}),
        }
            const response = await fetch("http://localhost:8081/lookup", requestOptions);
            this.assets = await response.json();
        },
        returnSellParameters: function() {
            this.$emit('returnSellParameters', this.sellForm)
        }
    }
}


</script>
