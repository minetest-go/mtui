import ChangePassword from './ChangePassword.js';
import ATMDisplay from './ATMDisplay.js';
import SkinPreview from '../SkinPreview.js';
import PrivBadge from '../PrivBadge.js';

import { has_priv, get_claims, is_logged_in } from "../../service/login.js";
import { has_feature } from "../../service/features.js";
import { get as get_playerinfo } from '../../api/playerinfo.js';
import { get_record } from '../../api/xban.js';
import format_time from '../../util/format_time.js';
import format_duration from '../../util/format_duration.js';
import format_count from '../../util/format_count.js';

export default {
    props: ["username","show_token_link"],
    data: function() {
        return {
            playerinfo: null,
            xban_record: null
        };
    },
    mounted: function() {
        this.update();
    },
    components: {
        "change-password": ChangePassword,
        "atm-display": ATMDisplay,
        "skin-preview": SkinPreview,
        "priv-badge": PrivBadge
    },
    computed: {
        can_change_pw: function() {
            return (is_logged_in() && get_claims().username == this.username) || has_priv("password");
        },
        is_moderator: function() {
            return has_priv("ban") || has_priv("server");
        },
        can_administer: function() {
            return has_priv("privs");
        }
    },
    methods: {
        format_time,
        format_duration,
        format_count,
        has_priv,
        has_feature,
        update_playerinfo: function() {
            get_playerinfo(this.username).then(pi => this.playerinfo = pi);
        },
        update: function() {
            this.playerinfo = null;
            this.update_playerinfo();

            this.xban_record = null;
            if (this.is_moderator && this.has_feature("xban")) {
                get_record(this.username).then(r => this.xban_record = r);
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
            &nbsp;
            <skin-preview :playername="username"/>
        </h3>
        <div class="alert alert-warning" v-if="!playerinfo.auth_entry || !playerinfo.player_entry">
            <i class="fa fa-triangle-exclamation"></i>
            <b>Warning:</b>
            <ul>
                <li v-if="!playerinfo.auth_entry">no entry in the auth-database found!</li>
                <li v-if="!playerinfo.player_entry">no entry in the player-database found! (the player hasn't logged in ingame yet)</li>
            </ul>
        </div>
        <div class="row">
            <div class="col-md-4">
                <div class="card">
                    <div class="card-header">
                        <i class="fa-solid fa-award"></i>
                        Privileges
                        <router-link :to="'/profile/' + username + '/priveditor'" class="btn btn-sm btn-primary float-end" v-if="can_administer">
                            <i class="fa fa-edit"></i>
                            Edit
                        </router-link>
                    </div>
                    <div class="card-body">
                        <div class="container">
                            <priv-badge :priv="priv" v-for="priv in playerinfo.privs" class="m-1"/>
                        </div>
                    </div>
                </div>
            </div>
            <div class="col-md-4">
                <div class="card">
                    <div class="card-header">
                        <i class="fa-solid fa-door-open"></i>
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
                <div class="card" v-if="has_feature('atm')">
                    <div class="card-header">
                        <i class="fa-solid fa-money-bill"></i>
                        ATM
                    </div>
                    <div class="card-body">
                        <atm-display :username="username"/>
                    </div>
                </div>
                <br v-if="has_feature('atm')">
                <div class="card">
                    <div class="card-header">
                        <i class="fa-solid fa-info"></i>
                        Ingame stats
                    </div>
                    <div class="card-body" v-if="playerinfo.stats">
                        <ul>
                            <li v-if="is_moderator">
                                <b>Position: </b> {{ parseInt(playerinfo.posx)/10 }} / {{ parseInt(playerinfo.posy)/10 }} / {{ parseInt(playerinfo.posz)/10 }}
                            </li>
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
            <div class="col-md-4">
                <div class="card" v-if="can_change_pw">
                    <div class="card-header">
                        <i class="fa-solid fa-key"></i>
                        Password
                    </div>
                    <div class="card-body">
                        <change-password :username="username"/>
                    </div>
                </div>
                <br>
                <div class="card" v-if="show_token_link && has_feature('api')">
                    <div class="card-header">
                        <i class="fa-solid fa-key"></i>
                        Access tokens
                    </div>
                    <div class="card-body">
                        <router-link to="/token" class="btn btn-secondary">
                            Configure access tokens
                        </router-link>
                    </div>
                </div>
                <br v-if="has_feature('api')">
                <div class="card" v-if="is_moderator && has_feature('xban') && xban_record">
                    <div class="card-header">
                        <i class="fa-solid fa-clipboard"></i>
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