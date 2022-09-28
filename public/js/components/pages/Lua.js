import { execute_lua } from "../../api/lua.js";
import login_store from '../../store/login.js';

const store = Vue.reactive({
    login_store: login_store,
    code: "",
    success: false,
    error: false,
    busy: false,
    result: "",
    message: "",
    delay: 0
});

export default {
    data: () => store,
    methods: {
        execute: function() {
            this.busy = true;
            this.error = false;
            this.message = "";
            this.result = "";
            this.success = false;
            this.delay = 0;
            const start = Date.now();

            execute_lua(this.code)
            .then(result => {
                this.busy = false;
                this.success = result.success;
                this.error = !this.success;
                this.message = result.message;
                this.result = result.result;
                this.delay = Date.now() - start;
            });
        }
    },
    template: /*html*/`
    <div>
        <form @submit.prevent="execute" class="row">
            <div class="col-md-10">
                <textarea rows="5" v-model="code" class="form-control"></textarea>
            </div>
            <div class="col-md-2">
                <button class="btn btn-outline-primary" type="submit" class="form-control" :disabled="!code">
                    Execute
                    &nbsp;
                    <i class="fa-solid fa-check" v-if="success" style="color: green;"></i>
                    <i class="fa-solid fa-xmark" v-if="error" style="color: red;"></i>
                    <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                    &nbsp;
                    <span v-if="delay">
                        ({{delay}} ms, {{JSON.stringify(result).length}} bytes)
                    </span>
                </button>
            </div>
        </form>
        <hr>
        <div class="alert alert-danger" v-if="error">
            <i class="fa-solid fa-triangle-exclamation"></i>
            {{message}}
        </div>
        <div class="row">
            <div class="col-md-12">
                <pre class="w-100" style="height: 500px; background-color: grey;">{{result}}</pre>
            </div>
        </div>
    </div>
    `
};