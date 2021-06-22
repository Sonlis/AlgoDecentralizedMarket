<template>
    <div>
        <button v-on:click="lookupSellings()">Get your escrow accounts</button>
        <div v-if="pressed">
            <div class="escrows" v-for="asset in escrows" :key="asset['address']" >
                <p>Address: {{asset['address']}}</p>
                <p>Algo Amount: {{asset['amount'] / 1000000 }}</p> <p>Withdraw <input type="number" id="algosamount" v-model.number="withdrawForm.algosAmount"> Algos <button v-on:click="returnWithdrawParameters(asset['address'], 0)">Withdraw</button></p>
                <div v-for="assets in asset.assets" :key="assets['asset-id']"><p v-if="assets['amount'] != 0"> {{assets['amount']}} of asset {{assets['asset-id']}}. Withdraw <input type="number" id="withdrawassetamount" placeholder=1000 v-model.number="withdrawForm.assetAmount"> of asset {{ assets['asset-id'] }} <button v-on:click="returnWithdrawParameters(asset['address'], assets['asset-id'])">Withdraw</button></p></div>
            </div>
        </div>
    </div>
</template>

<script> 

export default {
    name: 'UserEscrows',
    props: ['addrToUse'],
    data() {
        return{
            escrows: [],
            pressed: false,
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
          this.pressed = !this.pressed
        },
        returnWithdrawParameters: function(escrowAddress, assetID) {
            this.withdrawForm.escrowAddress = escrowAddress
            this.withdrawForm.assetID = assetID
            this.$emit('returnWithdrawParameters', this.withdrawForm)
        }

    }
}

</script>

<style scoped>

li {
    display: inline-block;
}

.escrows {
    margin-top: 2%;
    border: 1px dotted black;
}

input {
    width: 10%;
}

</style>
