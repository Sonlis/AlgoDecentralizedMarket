<template>
    <div class="myAlgoConnection">
        <div>
            <button v-if="!addrToUse" class="right" v-on:click="connectToMyAlgo()" ><span>Connect to myAlgo</span></button>
        </div>
        <div class="picker" v-if="pressed">
            <p>Which address would you like to use ?</p>
            <ul>
            <li v-for="address in addresses" :key="address"><input type="radio" v-model="addrToUse" :id="address" :value="address">{{ address }}</li>
            </ul>
            <button class="center" v-on:click="returnAddress()">Done</button>
        </div>
        <div class="showaddress" v-if="addrToUse && !pressed">
            <button v-on:click="pressed=!pressed">{{addrToUse}}</button>
            <br />
            <button class="right" v-on:click="addrToUse=''">Disconnect</button>
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

.right {
    float: right;
}

.center {
    float: center;
}

.picker {
    border-style: double;
    background-color: rgb(228, 222, 222);
}
.myAlgoConnection {
    color: #000;
    position: absolute;
    top: 5%;
    right: 5%;
}

.showaddress {
    font-size: 0.8em;
}

ul {
    list-style-type: none;
}

li {
    padding-left: none;
    font-size: 0.8em;
}
</style>