<template>
    <div>
        <div class="button">
            <button v-on:click="getAssets()" ><span>Sell an Asset</span></button>
        </div>
        <template v-if="pressed">
            <h4> Which asset to sell ?</h4>
            <div class="chooseAsset" v-for="(asset) in assets.assets" :key="asset['asset-id']">
                <label v-if="asset['amount']" v-bind:for="asset['asset-id']" :name="asset['asset-id']">asset {{asset["asset-id"]}}, held: {{asset["amount"]}}</label> 
                <input v-if="asset['amount']" type="radio" v-model.number="sellForm.assetID" :name="asset['asset-id']" :id="asset['asset-id']" :value="asset['asset-id']">
            </div>
            <label for="assetAmount">Amount of asset to sell</label>
            <input type="assetAmount"  id="assetAmount" placeholder=0
                v-model.number="sellForm.assetAmount">

            <label for="paymentAssetID">First asset ID to accept for trade (payment asset)</label>
            <input type="number"  id="paymentAssetID" placeholder=1000
                v-model.number="sellForm.paymentAssetID">

            <label for="paymentAssetAmount">First payment asset amount</label>
            <input type="number"  id="paymentAssetAmount" 
                v-model.number="sellForm.paymentAssetAmount">

            <label for="secondPaymentAssetID">Second asset ID to accept for trade (payment asset)</label>
            <input type="number"  id="secondPaymentAssetID" 
                v-model.number="sellForm.secondPaymentAssetID">

            <label for="secondPaymentAssetAmount">Second payment asset Amount</label>
            <input type="number"  id="secondPaymentAssetAmount" placeholder=1000
                v-model.number="sellForm.secondPaymentAssetAmount">

            <label for="algoAmount">MicroAlgo needed to buy the asset</label>
            <input type="number"  id="algoAmount" placeholder="1000"
                v-model.number="sellForm.algoAmount">
            <button v-on:click="returnSellParameters()">Sell the asset</button>
        </template>
    </div>
</template>


<script>


export default {
    name: 'SellAsset',
    props: ['addrToUse'],
    data(){
        return{
            assets: [''],
            pressed: false,
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
            if (this.pressed === false) {
                this.pressed = !this.pressed
            }
        },
        returnSellParameters: function() {
            this.$emit('returnSellParameters', this.sellForm)
        }
    }
}


</script>

<style scoped>
.chooseAsset {
    width: 50%;
    float: left;
}
</style>
