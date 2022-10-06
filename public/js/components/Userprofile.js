import ChangePassword from './ChangePassword.js';
import { has_priv } from "../service/login.js";
import { get as get_playerinfo } from '../api/playerinfo.js';
import format_time from '../util/format_time.js';
import login_store from '../store/login.js';

export default {
    props: ["username"],
    data: function() {
        return {
            playerinfo: null
        };
    },
    mounted: function() {
        get_playerinfo(this.username)
        .then(pi => this.playerinfo = pi);
    },
    components: {
        "change-password": ChangePassword
    },
    computed: {
        can_change_pw: function() {
            return (login_store.claims && login_store.claims.username == this.username) || has_priv("password");
        }
    },
    methods: {
        format_time: format_time,
        getPrivBadgeClass: function(priv) {
            if (priv == "server" || priv == "privs") {
                return { "badge": true, "bg-danger": true };
            } else if (priv == "ban" || priv == "kick") {
                return { "badge": true, "bg-primary": true };
            } else {
                return { "badge": true, "bg-secondary": true };
            }
        }
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
                        <ul v-if="playerinfo.auth_entry">
                            <li v-for="priv in playerinfo.privs">
                                <span v-bind:class="getPrivBadgeClass(priv)">{{ priv }}</span>
                            </li>
                        </ul>
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
                                <b>Health:</b> {{playerinfo.health}}
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
            </div>
        </div>
    </div>
    `
};