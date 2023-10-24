import { execute_lua } from "../../../api/lua.js";
import DefaultLayout from "../../layouts/DefaultLayout.js";
import CodeEditor from "../../CodeEditor.js";
import { START, ADMINISTRATION } from "../../Breadcrumb.js";

const store = Vue.reactive({
    code: "return minetest.features",
    success: false,
    error: false,
    busy: false,
    result: "",
    message: "",
    delay: 0,
    breadcrumb: [START, ADMINISTRATION, {
        name: "Lua",
        icon: "terminal",
        link: "/lua"
    }]
});

export default {
    components: {
        "default-layout": DefaultLayout,
        "code-editor": CodeEditor
    },
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
    <default-layout icon="terminal" title="Lua" :breadcrumb="breadcrumb">
        <form @submit.prevent="execute" class="row">
            <div class="col-md-10">
                <code-editor mode="lua" v-model="code" style="height: 500px;" class="w-100"/>
            </div>
            <div class="col-md-2">
                <button class="btn btn-outline-primary w-100" type="submit" :disabled="!code">
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
    </default-layout>
    `
};