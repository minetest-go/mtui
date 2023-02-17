import { get_events, count, search } from "../../api/log.js";
import format_time from "../../util/format_time.js";

const store = Vue.reactive({
    category: "minetest",
    events: [],
    event: "",
    username: "",
    count: 0,
    busy: false,
    from: new Date(Date.now() - (3600*1000*2)),
    to: new Date(Date.now() + (3600*1000*1)),
    logs: []
});

export default {
    data: () => store,
    methods: {
        update_events: function() {
            this.busy = true;
            get_events(this.category)
            .then(e => {
                this.events = e;
                if (e.length) {
                    this.event = e[0];
                }
                this.update_count();
            });
        },
        search_query() {
            return {
                category: this.category,
                event: this.event != "" ? this.event : null,
                username: this.username != "" ? this.username : null,
                from_timestamp: +this.from,
                to_timestamp: +this.to
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
        format_time: format_time
    },
    mounted: function() {
        this.update_events();
    },
    watch: {
        "category": "update_events",
        "event": "update_count",
        "from": "update_count",
        "to": "update_count",
        "username": "update_count"
    },
    template: /*html*/`
    <div>
        <div class="row">
            <div class="col-md-2">
                <label>Category</label>
                <select class="form-control" v-model="category">
                    <option value="minetest">Minetest</option>
                    <option value="ui">UI</option>
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
        <hr>
        <table class="table table-striped table-condensed">
            <thead>
                <tr>
                    <th>Event</th>
                    <th>Playername</th>
                    <th>Time</th>
                    <th>Message</th>
                    <th>IP Address</th>
                    <th>Position</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="log in logs">
                    <td>
                        <span class="badge bg-secondary">{{log.event}}</span>
                    </td>
                    <td>
                        <i class="fa fa-magnifying-glass" v-on:click="search_specific('username', log.username)"></i>
                        &nbsp;
                        <router-link :to="'/profile/' + log.username" v-if="log.username">
                            {{log.username}}
                        </router-link>
                    </td>
                    <td>{{ format_time(log.timestamp/1000) }}</td>
                    <td>{{log.message}}</td>
                    <td>
                        {{log.ip_address}}
                        <span v-if="log.geo_country || log.geo_city" class="badge bg-success">
                            {{log.geo_country}} <span v-if="log.geo_city">/ {{log.geo_city}}</span>
                        </span>
                        <span v-if="log.geo_asn" class="badge bg-info">
                            ASN: {{log.geo_asn}}
                        <span>
                    </td>
                    <td>
                        <span v-if="log.posx != null">
                            {{log.posx}}/{{log.posy}}/{{log.posz}}
                        </span>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
    `
};