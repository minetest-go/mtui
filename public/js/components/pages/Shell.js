import { execute_chatcommand } from "../../api/chatcommand.js";
import { has_priv, get_claims } from "../../service/login.js";
import events, { EVENT_DIRECT_CHAT } from "../../events.js";

const store = Vue.reactive({
    command: "",
    success: false,
    error: false,
    busy: false,
    message: "",
    delay: 0
});

events.on(EVENT_DIRECT_CHAT, function(data) {
    store.message += data.text + "\n";
});

export default {
    data: () => store,
    methods: {
        has_priv: has_priv,
        predefined: function(cmd) {
            this.command = cmd;
            this.execute();
        },
        execute: function() {
            this.message = "";
            this.busy = true;
            this.error = false;
            this.success = false;
            this.delay = 0;
            const start = Date.now();

            execute_chatcommand(get_claims().username, this.command)
            .then(result => {
                this.busy = false;
                this.success = result.success;
                this.error = !this.success;
                this.message += result.message + "\n";
                this.delay = Date.now() - start;
            });
        },
        clear: function() {
            this.message = "";
        }
    },
    computed: {
        get_claims: get_claims
    },
    template: /*html*/`
    <div>
        <form @submit.prevent="execute" class="row">
            <div class="col-md-10">
                <input type="text" placeholder="Command" v-model="command" class="form-control"/>
            </div>
            <div class="col-md-2">
                <button class="btn btn-outline-primary w-100" type="submit" :disabled="!command">
                    Execute
                    &nbsp;
                    <i class="fa-solid fa-check" v-if="success" style="color: green;"></i>
                    <i class="fa-solid fa-xmark" v-if="error" style="color: red;"></i>
                    <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                    &nbsp;
                    <span v-if="delay">
                        ({{delay}} ms)
                    </span>
                </button>
            </div>
        </form>
        <div class="row">
            <div class="col-md-12">
                <div class="input-group">
                    <a class="btn btn-outline-secondary" v-on:click="predefined('status')" v-if="has_priv('interact')">status</a>
                    <a class="btn btn-outline-secondary" v-on:click="predefined('privs')" v-if="has_priv('interact')">privs</a>
                    <a class="btn btn-outline-secondary" v-on:click="predefined('days')" v-if="has_priv('interact')">days</a>
                    <a class="btn btn-outline-secondary" v-on:click="predefined('time 6000')" v-if="has_priv('settime')">time 6000</a>
                    <a class="btn btn-warning" v-on:click="predefined('shutdown 120 -r')" v-if="has_priv('server')">shutdown 120 -r</a>
                </div>
            </div>
        </div>
        <hr>
        <div class="row">
            <div class="col-md-12">
                <pre class="w-100" style="height: 300px; background-color: grey;">{{message}}</pre>
            </div>
        </div>
    </div>
    `
};