import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

import { init, execute, is_supported } from "../../util/wasm_helper.js";
import { get_join_password } from "../../api/join_password.js";
import { get_claims } from "../../service/login.js";

export default {
    inject: ["unmount"],
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            breadcrumb: [START, {
                name: "Play",
                icon: "play",
                link: "/play"
            }],
            busy: false,
            supported: is_supported()
        };
    },
    methods: {
        play: function() {
            get_join_password()
            .then(pw => {
                const claims = get_claims();
                this.unmount();
                init()
                .then(() => {
                    execute([
                        "--go",
                        "--address", "engine",
                        "--port", "30000",
                        "--name", claims.username,
                        "--password", pw
                    ]);
                });    
            });

        }
    },
    template: /*html*/`
        <default-layout title="Play" icon="play" :breadcrumb="breadcrumb">
            <div class="row">
                <div class="col-md-4"></div>
                <div class="col-md-4">
                    <div class="alert alert-info">
                        <i class="fa fa-info"></i>
                        <b>Note: </b> This is an experimental feature and may not yet work as expected
                    </div>
                    <a class="btn btn-success w-100" v-on:click="play" v-if="supported">
                        <i class="fa fa-spinner fa-spin" v-if="busy"></i>
                        <i class="fa fa-play" v-else></i>
                        Play
                    </a>
                    <div class="alert alert-warning" v-else>
                        <i class="fa fa-triangle-exclamation"></i>
                        Your browser or environment doesn't support
                        <a href="https://caniuse.com/?search=WebAssembly">webassembly</a>
                    </div>
                </div>
                <div class="col-md-4"></div>
            </div>
        </default-layout>
    `
};