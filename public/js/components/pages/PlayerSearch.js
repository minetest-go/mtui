import DefaultLayout from "../layouts/DefaultLayout.js";
import { START, PLAYER_SEARCH } from "../Breadcrumb.js";
import SkinPreview from "../SkinPreview.js";

import { search, count } from "../../api/playerinfo.js";
import format_time from '../../util/format_time.js';

const store = Vue.reactive({
    busy: false,
    result_count: 0,
    searchterm: "",
    list: [],
    breadcrumb: [START, PLAYER_SEARCH]
});

export default {
    data: () => store,
	components: {
		"default-layout": DefaultLayout,
        "skin-preview": SkinPreview
	},
    methods: {
        format_time: format_time,
        getRole: function(player) {
            if (player.privs.includes("server")) {
                return "admin";
            } else if (player.privs.includes("kick") || player.privs.includes("ban")) {
                return "moderator";
            } else {
                return "player";
            }
        },
        query: function() {
            return {
                usernamelike: `%${this.searchterm}%`,
                limit: 100
            };
        },
        search: function() {
            this.busy = true;
            this.list = [];
            search(this.query())
            .then(l => {
                this.list = l;
            })
            .finally(() => this.busy = false);
        },
        count: function() {
            this.result_count = 0;
            this.list = [];
            if (this.searchterm == "")
                return;

            this.busy = true;
            count(this.query())
            .then(c => this.result_count = c)
            .finally(() => this.busy = false);
        }
    },
    watch: {
        searchterm: "count"
    },
    template: /*html*/`
    <default-layout icon="magnifying-glass" title="Player search" :breadcrumb="breadcrumb">
        <div class="row">
            <div class="col-md-10">
                <input type="text" placeholder="Player-name" class="form-control" v-model="searchterm"/>
            </div>
            <div class="col-md-2">
                <button class="btn btn-primary w-100" :disabled="searchterm == '' || busy" v-on:click="search">
                    <i class="fa fa-magnifying-glass" v-if="!busy"></i>
                    <i class="fa-solid fa-spinner fa-spin" v-else></i>
                    Search
                    <span class="badge bg-secondary">{{result_count}}</span>
                </button>
            </div>
        </div>
        <table class="table table-striped table-condensed">
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Health</th>
                    <th>Privs</th>
                    <th>First login</th>
                    <th>Last login</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="p in list">
                    <td>
                        <span class="badge bg-secondary">{{p.auth_id}}</span>
                    </td>
                    <td>
                        <router-link :to="'/profile/' + p.name">
                            {{p.name}}
                        </router-link>
                        &nbsp;
                        <skin-preview :playername="p.name"/>
                    </td>
                    <td>
                        <i class="fa-solid fa-heart" style="color: red;"></i>
                        {{p.health}}
                    </td>
                    <td>
                        <span class="badge bg-danger" v-if="getRole(p) == 'admin'">Admin</span>
                        <span class="badge bg-primary" v-if="getRole(p) == 'moderator'">Moderator</span>
                        <span class="badge bg-secondary" v-if="getRole(p) == 'player'">Player</span>
                    </td>
                    <td>{{ format_time(p.first_login) }}</td>
                    <td>{{ format_time(p.last_login) }}</td>
                </tr>
            </tbody>
        </table>
    </default-layout>
    `
};