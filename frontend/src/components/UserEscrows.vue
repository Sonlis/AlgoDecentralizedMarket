<template>
    <div>
        <button v-on:click="lookupSellings()">Get your escrow accounts</button>
        <li v-for="asset in escrows" :key="asset['address']"><p>Address: {{asset['address']}}, Algo Amount: {{asset['amount']}}, <p v-for="assets in asset.assets" :key="assets['asset-id']"> {{assets['amount']}} of asset {{assets['asset-id']}}</p></li>
        <p>Withdraw <input type="number" class="form-control" id="algosamount" placeholder=1000 v-model.number="withdrawForm.algosAmount"> Algos from <input class="form-control" id="paymentAssetID" placeholder=1000 v-model="withdrawForm.escrowAddress"><button v-on:click="returnWithdrawParameters()">Withdraw</button></p>
        <p>Withdraw <input type="number" class="form-control" id="withdrawassetamount" placeholder=1000 v-model.number="withdrawForm.assetAmount"> of asset number <input type="number" class="form-control" id="withdrawassetid" placeholder=1000 v-model.number="withdrawForm.assetID"> from <input class="form-control" id="escrowaddress2" placeholder=1000 v-model="withdrawForm.escrowAddress"><button v-on:click="returnWithdrawParameters()">Withdraw</button></p>
    </div>
</template>



<script> 

export default {
    name: 'UserEscrows',
    props: ['addrToUse'],
    data() {
        return{
            escrows: [],
            withdrawForm: {
                escrowAddress: '',
                assetID: 0,
                assetAmount: 0,
                algosAmount: 0,
            },
            
        }
    },
    methods: {
        lookupSellings: async function() {
          const requestOptions = {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({"accountid": this.addrToUse}),
          }
          console.log(requestOptions.body)
          const response = await fetch("http://localhost:8081/lookupSellings", requestOptions);
          this.escrows = await response.json();
          console.log(this.escrows)
        },
        returnWithdrawParameters: function() {
            this.$emit('returnWithdrawParameters', this.withdrawForm)
        }

    }
}

</script>