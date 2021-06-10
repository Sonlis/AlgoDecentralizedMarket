<template>
    <div>
        <button v-on:click="getSellings()" ><span>Get assets being sold</span></button>
        <div class="sellings" v-for="(asset) in sellings.assets" :key="asset['asset-id']">
            <p>Asset {{asset['asset-id']}}</p>
                <button :name="asset['fassetid']" :id="asset['fassetid']" v-on:click="returnBuyParameters(asset['fassetid'], asset['fassetamount'], asset['address'], asset['asset-id'], 0)">Buy with {{asset["fassetamount"]}} of asset {{asset["fassetid"]}}</button>
                <button :name="asset['sassetid']" :id="asset['sassetid']" v-on:click="returnBuyParameters(asset['sassetid'], asset['sassetamount'], asset['address'], asset['asset-id'], 0)">Buy with {{asset['sassetamount']}} of asset {{asset["sassetid"]}}</button>
                <button :name="asset['algoamount']" :id="asset['algoamount']" v-on:click="returnBuyParameters(0, 0, asset['address'],  asset['asset-id'], asset['algoamount'] )">Buy with {{asset['algoamount']}} Algos</button>
        </div>
    </div>
</template>


<script> 

export default {
    name: 'BuyAsset',
    data() {
        return {
            sellings: [],
            buyForm: {
                paymentAssetID: 0,
                paymentAssetAmount: 0,
                algoAmount: 0,
                assetid: 0,
                address: ''
            }
        }
    },
    methods: {
        getSellings: async function() {
          const response = await fetch("http://localhost:8081/getSellings");
          this.sellings = await response.json();
          console.log(this.sellings)
          },
          returnBuyParameters: function(paymentAssetID, paymentAssetAmount, address, assetID,algoAmount) {
              this.buyForm.assetID = assetID;
              this.buyForm.paymentAssetID = paymentAssetID;
              this.buyForm.paymentAssetAmount = paymentAssetAmount;
              this.buyForm.address = address
              this.buyForm.algoAmount = algoAmount
              this.$emit('returnBuyParameters', this.buyForm);
          }
    }
}

</script>

<style scoped>

.sellings {
    padding-top: 1%;
    margin-top: 1%;
}

</style>