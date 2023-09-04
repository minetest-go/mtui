import { count, search } from "../../api/log.js";
import format_time from "../../util/format_time.js";

const store = Vue.reactive({
    category: "minetest",
    event: "",
    username: "",
    nodename: "",
    posx: "",
    posy: "",
    posz: "",
    count: 0,
    geo_asn: "",
    ip_address: "",
    busy: false,
    from: new Date(Date.now() - (3600*1000*2)),
    to: new Date(Date.now() + (3600*1000*1)),
    logs: []
});

export default {
    data: () => store,
    computed: {
        events: function() {
            if (this.category == "minetest") {
                return [
                    "prejoin",
                    "join",
                    "leave",
                    "authplayer",
                    "dieplayer",
                    "cheat",
                    "chat",
                    "system",
                    "on_generated",
                    "protection_violation",
                    "placenode",
                    "dignode",
                    "punchnode",
                    "craft",
                    "logfile",
                    "logfile-error",
                    "logfile-warning",
                    "logfile-action",
                    "logfile-info",
                    "logfile-verbose"
                ];
            } else {
                return [];
            }
        }
    },
    methods: {
        search_query() {
            return {
                category: this.category,
                event: this.event != "" ? this.event : null,
                username: this.username != "" ? this.username : null,
                nodename: this.nodename != "" ? this.nodename : null,
                posx: this.posx != "" ? parseInt(this.posx) : null,
                posy: this.posy != "" ? parseInt(this.posy) : null,
                posz: this.posz != "" ? parseInt(this.posz) : null,
                ip_address: this.ip_address != "" ? this.ip_address : null,
                from_timestamp: +this.from,
                to_timestamp: +this.to,
                geo_asn: this.geo_asn != "" ? parseInt(this.geo_asn) : null
            };
        },
        update_count: function() {
            this.busy = true;
            count(this.search_query())
            .then(c => {
                this.count = c;
                this.busy = false;
            });
        },
        search: function() {
            this.busy = true;
            this.logs = [];
            search(this.search_query())
            .then(l => {
                this.logs = l;
                this.busy = false;
            });
        },
        search_specific: function(field, value) {
            store[field] = value;
            this.search();
        },
        search_specific_pos: function(posx, posy, posz) {
            store.posx = posx;
            store.posy = posy;
            store.posz = posz;
            this.search();
        },
        get_row_class: function(log) {
            return {
                'table-success': log.event == 'join',
                'table-warning': log.event=='leave',
                'table-danger': log.event == 'logfile-error'
            };
        },
        format_time: format_time
    },
    watch: {
        "event": "update_count",
        "from": "update_count",
        "to": "update_count",
        "username": "update_count",
        "nodename": "update_count",
        "posx": "update_count",
        "posy": "update_count",
        "posz": "update_count",
        "geo_asn": "update_count",
        "ip_address": "update_count"
    },
    mounted: function() {
        this.update_count();
    },
    template: /*html*/`
    <div>
        <div class="row">
            <div class="col-md-2">
                <label>Category</label>
                <select class="form-control" v-model="category">
                    <option value="minetest">Minetest</option>
                </select>
            </div>
            <div class="col-md-2">
                <label>Event</label>
                <select class="form-control" v-model="event">
                    <option value="">*</option>
                    <option v-for="e in events" :value="e">{{e}}</option>
                </select>
            </div>
            <div class="col-md-2">
                <label>Playername</label>
                <input type="text" class="form-control" v-model="username"/>
            </div>
            <div class="col-md-2">
                <label>From</label>
                <vue-datepicker v-model="from"/>
            </div>
            <div class="col-md-2">
                <label>To</label>
                <vue-datepicker v-model="to"/>
            </div>
            <div class="col-md-2">
                <label>Search</label>
                <a class="btn btn-primary w-100" v-on:click="search">
                    <i class="fa fa-magnifying-glass" v-if="!busy"></i>
                    <i class="fa-solid fa-spinner fa-spin" v-else></i>
                    Search
                    <span class="badge bg-secondary">{{count}}</span>
                </a>
            </div>
        </div>
        <div class="row">
            <div class="col-md-2">
                <label>Nodename</label>
                <input type="text" class="form-control" v-model="nodename"/>
            </div>
            <div class="col-md-2">
                <label>Pos X</label>
                <input type="number" class="form-control" v-model="posx"/>
            </div>
            <div class="col-md-2">
                <label>Pos Y</label>
                <input type="number" class="form-control" v-model="posy"/>
            </div>
            <div class="col-md-2">
                <label>Pos Z</label>
                <input type="number" class="form-control" v-model="posz"/>
            </div>
            <div class="col-md-2">
                <label>Geo ASN</label>
                <input type="number" class="form-control" v-model="geo_asn"/>
            </div>
            <div class="col-md-2">
                <label>IP Address</label>
                <input type="text" class="form-control" v-model="ip_address"/>
            </div>
        </div>
        <hr>
        <table class="table table-striped table-condensed">
            <thead>
                <tr>
                    <th>Event</th>
                    <th>Playername</th>
                    <th>Time</th>
                    <th>Message</th>
                    <th>Nodename</th>
                    <th>Position</th>
                    <th>IP Address</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="log in logs" v-bind:class="get_row_class(log)">
                    <td>
                        <i v-if="log.event == 'on_generated'" class="fa fa-map"></i>
                        <i v-if="log.event == 'prejoin'" class="fa fa-shield"></i>
                        <i v-if="log.event == 'authplayer'" class="fa fa-key"></i>
                        <i v-if="log.event == 'join'" class="fa fa-right-from-bracket"></i>
                        <i v-if="log.event == 'leave'" class="fa fa-right-to-bracket"></i>
                        <i v-if="log.event == 'placenode'" class="fa fa-plus"></i>
                        <i v-if="log.event == 'punchnode'" class="fa fa-hand-fist"></i>
                        <i v-if="log.event == 'dignode'" class="fa fa-minus"></i>
                        <i v-if="log.event == 'craft'" class="fa fa-cart-shopping"></i>
                        <i v-if="log.event == 'chat'" class="fa fa-message"></i>
                        <i v-if="log.event == 'system'" class="fa fa-gear"></i>
                        <i v-if="log.event == 'dieplayer'" class="fa fa-tombstone"></i>
                        <i v-if="log.event == 'cheat'" class="fa fa-question"></i>
                        <i v-if="log.event == 'protection_violation'" class="fa fa-question"></i>
                        <i v-if="log.event == 'logfile'" class="fa fa-file"></i>
                        <i v-if="log.event == 'logfile-error'" class="fa fa-file"></i>
                        <i v-if="log.event == 'logfile-warning'" class="fa fa-file"></i>
                        <i v-if="log.event == 'logfile-action'" class="fa fa-file"></i>
                        <i v-if="log.event == 'logfile-info'" class="fa fa-file"></i>
                        <i v-if="log.event == 'logfile-verbose'" class="fa fa-file"></i>
                        &nbsp;
                        <span class="badge bg-secondary">
                            {{log.event}}
                            <i class="fa fa-magnifying-glass" v-on:click="search_specific('event', log.event)"></i>
                        </span>
                    </td>
                    <td>
                        <i class="fa fa-magnifying-glass" v-on:click="search_specific('username', log.username)" v-if="log.username"></i>
                        &nbsp;
                        <router-link :to="'/profile/' + log.username" v-if="log.username">
                            {{log.username}}
                        </router-link>
                    </td>
                    <td>{{ format_time(log.timestamp/1000) }}</td>
                    <td>{{log.message}}</td>
                    <td>
                        <i class="fa fa-magnifying-glass" v-if="log.nodename" v-on:click="search_specific('nodename', log.nodename)"></i>
                        &nbsp;
                        {{log.nodename}}
                    </td>
                    <td>
                        <span v-if="log.posx != null">
                            <i class="fa fa-magnifying-glass" v-on:click="search_specific_pos(log.posx, log.posy, log.posz)"></i>
                            &nbsp;
                            {{log.posx}}/{{log.posy}}/{{log.posz}}
                        </span>
                    </td>
                    <td>
                        <i class="fa fa-magnifying-glass" v-on:click="search_specific('ip_address', log.ip_address)" v-if="log.ip_address"></i>
                        {{log.ip_address}}
                        <span v-if="log.geo_country || log.geo_city" class="badge bg-success">
                            {{log.geo_country}} <span v-if="log.geo_city">/ {{log.geo_city}}</span>
                        </span>
                        <span v-if="log.geo_asn" class="badge bg-info">
                            ASN:
                            <i class="fa fa-magnifying-glass" v-on:click="search_specific('geo_asn', log.geo_asn)"></i>
                            {{log.geo_asn}}
                        </span>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
    `
};