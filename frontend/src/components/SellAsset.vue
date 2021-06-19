<template>
    <div>
        <div class="button">
            <button v-on:click="getAssets()" ><span>Sell an Asset</span></button>
        </div>
        <template v-if="pressed">
            <h4> Which asset to sell ?</h4>
            <div class="chooseAsset" v-for="(asset) in assets.assets" :key="asset['asset-id']">
                <label v-if="asset['amount']" v-bind:for="asset['asset-id']" :name="asset['asset-id']">asset {{asset["asset-id"]}}, held: {{asset["amount"]}}</label> 
                <input class="{ 'form-group--error': $v.sellForm.assetID.$error }" v-if="asset['amount']" type="radio" v-model.number="$v.sellForm.assetID.$model" :name="asset['asset-id']" :id="asset['asset-id']" :value="asset['asset-id']">
            </div>
            <div class="error" v-if="$v.sellForm.assetID.$invalid && submitStatus === 'ERROR'">Please select an asset to sell</div>
            <br />
            <div class="form">
                <div class="form-line">
                    <h5> Amount of asset {{ sellForm.assetID }} to sell</h5>
                    <label for="assetAmount">Amount of asset to sell</label>
                    <input type="assetAmount"  id="assetAmount" placeholder=0
                        v-model.number="sellForm.assetAmount">
                    <div class="error" v-if="$v.sellForm.assetAmount.$invalid && submitStatus === 'ERROR'">Please enter an amount of asset to sell</div>

                    <label for="algoAmount">MicroAlgo needed to buy the asset</label>
                    <input type="number"  id="algoAmount" placeholder="10000"
                        v-model.number="sellForm.algoAmount">
                    <div class="error" v-if="$v.sellForm.algoAmount.$invalid && submitStatus === 'ERROR'">Please enter an amount of algo to buy the asset</div>
                    <p>You can also sell the asset for up to two other payment assets</p>
                    <input type="radio" v-model.number="paymentAssets" id="none" value="0">
                    <label for="none"> No asset</label> 
                    <input type="radio" v-model.number="paymentAssets" id="one" value="1">
                    <label for="one"> 1 Asset</label> 
                    <input type="radio" v-model.number="paymentAssets" id="two" value="2">
                    <label for="two"> 2 Assets</label> 
                </div>
                <div class="form-line" v-if="paymentAssets == 1 || paymentAssets == 2">
                    <h5> First payment asset</h5>
                    <label for="USDC">USDC</label>
                    <input type="radio" id="USDC" v-model.number="sellForm.paymentAssetID" value="31566704">
                    <label for="USDt">USDt</label>
                    <input type="radio" id="USDt" v-model.number="sellForm.paymentAssetID" value="312769">
                    <label for="EUR">Monerium EUR</label>
                    <input type="radio" id="EUR" v-model.number="sellForm.paymentAssetID" value="83209012">
                    <label for="paymentAssetID">Other (insert asset ID)</label>
                    <input type="number"  id="paymentAssetID" placeholder=1000
                        v-model.number="sellForm.paymentAssetID">
                    <label for="paymentAssetAmount">First payment asset amount</label>
                    <input type="number"  id="paymentAssetAmount" 
                        v-model.number="sellForm.paymentAssetAmount">
                </div>
                <div class="form-line" v-if="paymentAssets == 2">
                    <h5> Second payment asset</h5>
                    <label for="USDC">USDC</label>
                    <input type="radio" id="USDC" v-model.number="sellForm.secondPaymentAssetID" value="31566704">
                    <label for="USDt">USDt</label>
                    <input type="radio" id="USDt" v-model.number="sellForm.secondPaymentAssetID" value="312769">
                    <label for="EUR">Monerium EUR</label>
                    <input type="radio" id="EUR" v-model.number="sellForm.secondPaymentAssetID" value="83209012">
                    <label for="secondPaymentAssetID">Other (insert asset ID)</label>
                    <input type="number"  id="secondPaymentAssetID" 
                        v-model.number="sellForm.secondPaymentAssetID">

                    <label for="secondPaymentAssetAmount">Second payment asset amount</label>
                    <input type="number"  id="secondPaymentAssetAmount" placeholder=1000
                        v-model.number="sellForm.secondPaymentAssetAmount">
                </div>
            </div>
            <button v-on:click="returnSellParameters()">Sell the asset</button>
            <p class="error" v-if="submitStatus === 'ERROR'">Please fill the form correctly.</p>
            
        </template>
    </div>
</template>


<script>
import { required, minValue } from 'vuelidate/lib/validators'

export default {
    name: 'SellAsset',
    props: ['addrToUse'],
    data(){
        return{
            assets: [''],
            pressed: false,
            submitStatus: null,
            paymentAssets: 0,
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
    validations: {
        sellForm: {
            assetID: { required, minValue: minValue(1), },
            assetAmount: { required, minValue: minValue(1) },
            algoAmount: { required, minValue: minValue(1) },
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
            this.$v.$touch()
            if (this.$v.$invalid) {
                this.submitStatus = 'ERROR'
            } else {
                this.submitStatus = 'PENDING'
                this.$emit('returnSellParameters', this.sellForm)
        }
    }
}
}


</script>

<style scoped>

.chooseAsset {
    width: 50%;
    display: inline-block;
}

.form-line {
    margin-top: 2%;
    border-bottom: 1px dotted black;
}

input {
    width: 8%;
    margin-right: 3%;
  }
</style>
