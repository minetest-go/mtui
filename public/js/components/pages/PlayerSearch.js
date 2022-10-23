
import { search, count } from "../../api/playerinfo.js";
import format_time from '../../util/format_time.js';

const store = Vue.reactive({
    busy: false,
    count: 0,
    searchterm: "",
    list: []
});

export default {
    data: () => store,
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
            this.count = 0;
            this.list = [];
            if (this.searchterm == "")
                return;

            this.busy = true;
            count(this.query())
            .then(c => {
                this.count = c;
            })
            .finally(() => this.busy = false);
        }
    },
    watch: {
        searchterm: "count"
    },
    template: /*html*/`
    <div>
        <div class="row">
            <div class="col-md-10">
                <input type="text" placeholder="Player-name" class="form-control" v-model="searchterm"/>
            </div>
            <div class="col-md-2">
                <button class="btn btn-primary w-100" :disabled="searchterm == '' || busy" v-on:click="search">
                    <i class="fa fa-magnifying-glass" v-if="!busy"></i>
                    <i class="fa-solid fa-spinner fa-spin" v-else></i>
                    Search
                    <span class="badge bg-secondary">{{count}}</span>
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
    </div>
    `
};