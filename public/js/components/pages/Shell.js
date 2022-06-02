import { execute_chatcommand } from "../../api/chatcommand.js";
import login_store from '../../store/login.js';

export default {
    data: function() {
        return {
            claims: login_store.claims,
            command: "",
            success: false,
            error: false,
            busy: false,
            message: ""
        };
    },
    methods: {
        execute: function() {
            this.busy = true;
            this.error = false;
            this.success = false;

            execute_chatcommand(this.claims.username, this.command)
            .then(result => {
                this.busy = false;
                this.success = result.success;
                this.error = !this.success;
                this.message = result.message;
            });
        }
    },
    template: /*html*/`
    <div>
        <form @submit.prevent="execute" class="row">
            <div class="col-md-10">
                <input type="text" placeholder="Command" v-model="command" class="form-control"/>
            </div>
            <div class="col-md-2">
                <button class="btn btn-primary" type="submit" class="form-control" :disabled="!command">
                    Execute
                    <i class="fa-solid fa-check" v-if="success" style="color: green;"></i>
                    <i class="fa-solid fa-xmark" v-if="error" style="color: red;"></i>
                    <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                </button>
            </div>
        </form>
        <hr>
        <div class="row">
            <div class="col-md-12">
                <pre class="w-100" style="height: 300px; background-color: grey;">{{message}}</pre>
            </div>
        </div>
    </div>
    `
};