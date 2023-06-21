import ChangePassword from './ChangePassword.js';
import { has_priv } from "../service/login.js";
import { has_feature } from "../service/features.js";
import { get as get_playerinfo } from '../api/playerinfo.js';
import { get_record } from '../api/xban.js';
import format_time from '../util/format_time.js';
import format_duration from '../util/format_duration.js';
import format_count from '../util/format_count.js';
import login_store from '../store/login.js';
import { execute_chatcommand } from "../api/chatcommand.js";

export default {
    props: ["username"],
    data: function() {
        return {
            playerinfo: null,
            xban_record: null,
            new_priv: ""
        };
    },
    mounted: function() {
        this.update();
    },
    components: {
        "change-password": ChangePassword
    },
    computed: {
        can_change_pw: function() {
            return (login_store.claims && login_store.claims.username == this.username) || has_priv("password");
        },
        is_moderator: function() {
            return has_priv("ban");
        }
    },
    methods: {
        update_playerinfo: function() {
            get_playerinfo(this.username)
            .then(pi => this.playerinfo = pi);
        },
        update: function() {
            this.playerinfo = null;
            this.update_playerinfo();

            this.xban_record = null;
            if (this.is_moderator && this.has_feature("xban")) {
                get_record(this.username).then(r => this.xban_record = r);
            }
        },
        has_priv: has_priv,
        has_feature: has_feature,
        revoke_priv: function(priv) {
            execute_chatcommand(login_store.claims.username, `revoke ${this.username} ${priv}`)
            .then(() => this.update_playerinfo());
        },
        grant_priv: function() {
            execute_chatcommand(login_store.claims.username, `grant ${this.username} ${this.new_priv}`)
            .then(() => {
                this.new_priv = "";
                this.update_playerinfo();
            });
        },
        format_time: format_time,
        format_duration: format_duration,
        format_count: format_count,
        getPrivBadgeClass: function(priv) {
            if (priv == "server" || priv == "privs") {
                return { "badge": true, "bg-danger": true };
            } else if (priv == "ban" || priv == "kick") {
                return { "badge": true, "bg-primary": true };
            } else if (priv == "otp_enabled") {
                return { "badge": true, "bg-success": true };
            } else {
                return { "badge": true, "bg-secondary": true };
            }
        }
    },
    watch: {
        "username": "update"
    },
    template: /*html*/`
    <div v-if="playerinfo">
        <h3>
            Profile for
            <small class="text-muted">
                {{ playerinfo.name }}
            </small>
        </h3>
        <div class="alert alert-danger" v-if="!playerinfo.auth_entry || !playerinfo.player_entry">
            <i class="fa fa-triangle-exclamation"></i>
            <b>Warning:</b>
            <ul>
                <li v-if="!playerinfo.auth_entry">no entry in the auth-database found!</li>
                <li v-if="!playerinfo.player_entry">no entry in the player-database found!</li>
            </ul>
        </div>
        <div class="row">
            <div class="col-md-4">
                <div class="card">
                    <div class="card-header">
                        Privileges
                    </div>
                    <div class="card-body">
                        <table>
                            <tr v-for="priv in playerinfo.privs">
                                <td>
                                    <span v-bind:class="getPrivBadgeClass(priv)">
                                        {{ priv }}
                                    </span>
                                </td>
                                <td>
                                    <a class="btn btn-sm btn-warning" v-if="has_priv('privs')" v-on:click="revoke_priv(priv)">
                                        <i class="fa fa-times"></i>
                                        Revoke
                                    </a>
                                </td>
                            </tr>
                            <tr v-if="has_priv('privs')">
                                <td>
                                    <input type="text" class="form-control" v-model="new_priv"/>
                                </td>
                                <td>
                                    <button class="btn btn-sm btn-success"
                                        v-if="has_priv('privs')"
                                        :disabled="new_priv == ''"
                                        v-on:click="grant_priv">
                                        <i class="fa fa-plus"></i>
                                        Grant
                                    </button>
                                </td>
                            </tr>
                        </table>
                    </div>
                </div>
            </div>
            <div class="col-md-4">
                <div class="card">
                    <div class="card-header">
                        Login stats
                    </div>
                    <div class="card-body" v-if="playerinfo.player_entry">
                        <ul>
                            <li>
                                <b>First login:</b> {{ format_time(playerinfo.first_login) }}
                            </li>
                            <li>
                                <b>Last login:</b> {{ format_time(playerinfo.last_login) }}
                            </li>
                            <li>
                                <b>Health:</b> <i class="fa-solid fa-heart" style="color: red;"></i> {{playerinfo.health}}
                            </li>
                            <li>
                                <b>Breath:</b> {{playerinfo.breath}}
                            </li>
                        </ul>
                    </div>
                    <div class="card-body" v-else>
                        No entries found
                    </div>
                </div>
                <br>
                <div class="card">
                    <div class="card-header">
                        Ingame stats
                    </div>
                    <div class="card-body" v-if="playerinfo.stats">
                        <ul>
                            <li v-if="playerinfo.stats.played_time">
                                <b>Playtime:</b> {{ format_duration(playerinfo.stats.played_time) }}
                            </li>
                            <li v-if="playerinfo.stats.digged_nodes">
                                <b>Digged nodes:</b> {{ format_count(playerinfo.stats.digged_nodes) }}
                            </li>
                            <li v-if="playerinfo.stats.placed_nodes">
                                <b>Placed nodes:</b> {{ format_count(playerinfo.stats.placed_nodes) }}
                            </li>
                            <li v-if="playerinfo.stats.died">
                                <b>Died:</b> {{ format_count(playerinfo.stats.died) }}
                            </li>
                            <li v-if="playerinfo.stats.crafted">
                                <b>Crafted:</b> {{ format_count(playerinfo.stats.crafted) }}
                            </li>
                        </ul>
                    </div>
                    <div class="card-body" v-else>
                        No stats found
                    </div>
                </div>
            </div>
            <div class="col-md-4" v-if="can_change_pw">
                <div class="card">
                    <div class="card-header">
                        Password
                    </div>
                    <div class="card-body">
                        <change-password :username="username"/>
                    </div>
                </div>
                <br>
                <div class="card" v-if="is_moderator && has_feature('xban') && xban_record">
                    <div class="card-header">
                        XBan record
                        <span v-if="xban_record.banned" class="badge bg-danger">Banned</span>
                    </div>
                    <div class="card-body">
                        <p v-if="xban_record.reason">
                            <b>Current ban-reason:</b> {{xban_record.reason}}
                        </p>
                        <p v-if="xban_record.last_seen">
                            <b>Last seen:</b> {{format_time(xban_record.last_seen)}}
                        </p>
                        <h5>Names</h5>
                        <ul>
                            <li v-for="_, name in xban_record.names">
                                <router-link :to="'/profile/' + name">
                                    {{name}}
                                </router-link>
                            </li>
                        </ul>
                        <h5>Records</h5>
                        <ul v-if="xban_record.record">
                            <li v-for="record in xban_record.record">
                                <b>Source:</b> <router-link :to="'/profile/' + record.source">{{record.source}}</router-link>
                                <b>Reason:</b> {{record.reason}}
                                <b>Time:</b> {{format_time(record.time)}}
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>
    `
};