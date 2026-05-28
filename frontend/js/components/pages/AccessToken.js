import format_duration from "../../util/format_duration.js";
import { get_privs } from "../../service/login.js";
import { generate_token } from "../../api/token.js";

import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            breadcrumb: [START, {
                name: `Access tokens`,
                icon: "key",
                link: `/token`
            }],
            expiry: new Date(Date.now() + (3600*1000*24*30)),
            privs: get_privs(),
            selected_privs: [],
            copy_success: false,
            token: "",
            base_url: window.location.href.split("#")[0]
        };
    },
    watch: {
        expiry: "clear_token",
        selected_privs: "clear_token",
    },
    methods: {
        format_duration,
        check_priv: function(priv, enabled) {
            if (enabled) {
                this.selected_privs.push(priv);
            } else {
                this.selected_privs = this.selected_privs.filter(p => p != priv);
            }
        },
        clear_token: function() {
            this.token = "";
            this.copy_success = false;
        },
        generate_token: async function() {
            this.copy_success = false;
            this.token = "";
            this.token = await generate_token(+this.expiry, this.selected_privs);
        },
        copy: async function() {
            await navigator.clipboard.writeText(this.token);
            this.copy_success = true;
        }
    },
    template: /*html*/`
        <default-layout title="Access tokens" icon="key" :breadcrumb="breadcrumb">
            <table class="table table-striped table-condensed">
                <tbody>
                    <tr>
                        <td>Expiration date</td>
                        <td>
                            <vue-datepicker v-model="expiry" auto-apply :min-date="new Date()"/>
                        </td>
                    </tr>
                    <tr>
                        <td>Expires in</td>
                        <td>{{format_duration((expiry - Date.now())/1000)}}</td>
                    </tr>
                    <tr>
                        <td>Privileges</td>
                        <td>
                            <div class="form-check" v-for="priv in privs">
                                <input class="form-check-input" type="checkbox" v-on:change="e => check_priv(priv, e.target.checked)">
                                <label class="form-check-label">
                                    {{priv}}
                                </label>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
            
            <button class="btn btn-primary w-100" v-on:click="generate_token" :disabled="selected_privs.length == 0">
                <i class="fa fa-key"></i>
                Generate token
            </button>
            <div class="form-floating">
                <textarea class="form-control" disabled v-model="token"></textarea>
                <label>Access token</label>
            </div>
            <button class="btn btn-secondary w-100" :disabled="!token" v-on:click="copy">
                <i class="fa fa-clipboard"></i>
                Copy to clipboard
                <i class="fa fa-check" style="color: green;" v-if="copy_success"></i>
            </button>

            <hr>

            <h4>Api examples with cURL</h4>

            <h5 class="text-muted">Execute chatcommand</h5>

            <pre>
curl -H "Authorization: Bearer \${TOKEN}" --data '{"command":"status"}' {{base_url}}api/bridge/execute_chatcommand
            </pre>

            <h5 class="text-muted">Execute lua command (with "server" priv only)</h5>

            <pre>
curl -H "Authorization: Bearer \${TOKEN}" --data '{"code":"return core.features"}' {{base_url}}api/bridge/lua
            </pre>
        </default-layout>
    `
};