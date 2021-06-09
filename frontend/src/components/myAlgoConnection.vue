<template>
    <div class="myAlgoConnection">
        <button v-on:click="connectToMyAlgo()" ><span>Connect to myAlgo</span></button>
        <div v-if="pressed">
            <p>Which address would you like to use ?</p>
            <li v-for="address in addresses" :key="address"><input type="radio" v-model="addrToUse" :id="address" :value="address">{{ address }}</li>
            <button v-on:click="returnAddress()">Done</button>
        </div>
    </div>
</template>


<script> 
import MyAlgo from '@randlabs/myalgo-connect';
const myAlgoWallet = new MyAlgo();

export default {
    name: "myAlgoConnection",
    data() {
        return {
            addrToUse: '',
            addresses: [],
            pressed: false
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
            if (this.pressed == false) {
                this.pressed = !this.pressed
            }
        },
        returnAddress: function () {
            this.$emit('clicked', this.addrToUse)
            this.pressed = !this.pressed
        }
    }

}
</script>

<style scoped>
</style>