import DefaultLayout from "../layouts/DefaultLayout.js";
import PrivBadge from "../PrivBadge.js";
import { START, PLAYER_SEARCH, PLAYER_INFO } from "../Breadcrumb.js";

import { get as get_playerinfo } from "../../api/playerinfo.js";
import { get_priv_infos } from "../../api/uimod.js";
import { execute_chatcommand } from "../../api/chatcommand.js";
import { get_claims } from "../../service/login.js";

export default {
    props: ["name"],
    components: {
        "default-layout": DefaultLayout,
        "priv-badge": PrivBadge
    },
    data: function() {
        return {
            breadcrumb: [START, PLAYER_SEARCH, PLAYER_INFO(this.name), {
                name: "Privilege editor",
                icon: "award"
            }],
            playerinfo: null,
            privs: {},
            busy: {}
        };
    },
    mounted: async function() {
        this.privs = await get_priv_infos();
        await this.update();
    },
    methods: {
        update: async function() {
            this.playerinfo = await get_playerinfo(this.name);
        },
        has_priv: function(priv) {
            return this.playerinfo.privs.includes(priv);
        },
        revoke_priv: async function(priv) {
            this.busy[priv] = true;
            await execute_chatcommand(get_claims().username, `revoke ${this.name} ${priv}`);
            await this.update();
            this.busy[priv] = false;
        },
        grant_priv: async function(priv) {
            this.busy[priv] = true;
            await execute_chatcommand(get_claims().username, `grant ${this.name} ${priv}`);
            await this.update();
            this.busy[priv] = false;
        }
    },
    template: /*html*/`
        <default-layout title="Privilege editor" icon="award" :breadcrumb="breadcrumb">
            <table class="table table-condensed table-striped" v-if="playerinfo">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Description</th>
                        <th>State</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="(def, name) in privs" v-bind:class="{'table-success': has_priv(name)}">
                        <td>
                            <priv-badge :priv="name"/>
                        </td>
                        <td>{{def.description}}</td>
                        <td>
                            <button class="btn btn-warning" v-if="has_priv(name)" v-on:click="revoke_priv(name)" :disabled="busy[name]">
                                <i class="fa fa-spinner fa-spin" v-if="busy[name]"></i>
                                <i class="fa fa-times" v-else></i>
                                Revoke
                            </button>
                            <button class="btn btn-success" v-else v-on:click="grant_priv(name)">
                                <i class="fa fa-spinner fa-spin" v-if="busy[name]"></i>
                                <i class="fa fa-plus" v-else></i>
                                Grant
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </default-layout>
    `
};