import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

import { init, execute } from "../../util/wasm_helper.js";
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
            }]
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
                    <a class="btn btn-success w-100" v-on:click="play">
                        <i class="fa fa-play"></i>
                        Play
                    </a>
                </div>
                <div class="col-md-4"></div>
            </div>
        </default-layout>
    `
};