import { get_balance, atm_transfer } from '../api/atm.js';
import { get_claims, is_logged_in } from "../service/login.js";


export default {
    props: ["username"],
    data: function() {
        return {
            balance: 0,
            target: "",
            amount: 0,
            errmsg: "",
            busy: false,
            done: false
        };
    },
    mounted: function() {
        this.update();
    },
    methods: {
        update: function() {
            get_balance(this.username).then(r => this.balance = r.balance);
        },
        transfer: function() {
            this.errmsg = "";
            this.done = false;
            this.busy = true;

            atm_transfer({
                target: this.target,
                amount: this.amount
            })
            .then(r => {
                if (r.success) {
                    this.done = true;
                    this.amount = 0;
                    this.balance = r.source_balance;
                } else {
                    this.errmsg = r.errmsg;
                }

                this.busy = false;
            });
        },
        format_money: function(x) {
            // https://stackoverflow.com/questions/2901102/how-to-format-a-number-with-commas-as-thousands-separators
            return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, "'");
        }
    },
    computed: {
        can_transfer: function() {
            return (is_logged_in() && get_claims().username == this.username);
        }
    },
    template: /*html*/`
    <div>
        Your balance: <b>$ {{ format_money(balance) }}</b>
        <div v-if="can_transfer">
            <hr>
            <h5>Wiretransfer</h5>
            <span class="input-group">
                <input type="text" class="form-control" placeholder="Receiving player" v-model="target" :disabled="busy"/>
                <input type="number" min="0" class="form-control" placeholder="Amount" v-model="amount" :disabled="busy"/>
                <button class="btn btn-secondary" v-on:click="transfer" :disabled="!target || !amount || busy">
                    <i class="fa fa-spinner fa-spin" v-if="busy"></i>
                    <i class="fa-solid fa-credit-card" v-else></i>
                    Transfer
                    <i class="fa-solid fa-check" v-if="done"></i>
                </button>
            </span>
            <div class="alert alert-danger" v-if="errmsg">
                <i class="fa fa-exclamation-mark"></i>
                <b>Wiretransfer error: </b> {{errmsg}}
            </div>
        </div>
    </div>
    `
};
