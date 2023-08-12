import { get_logs } from "../../../api/service_engine.js";
import { store as engine_store } from "./Engine.js";

const store = Vue.reactive({
    busy: false,
    logs: "",
    linecount: 0,
    live: true,
    since: Date.now() - (1000*60*60),
    until: Date.now() + (1000*60*60),
    logs_live_since: Date.now() - (1000*60*60) // shifting window for live-view
});

export default {
    data: function(){
        return store;
    },
    created: function(){
        this.log_update_handle = setInterval(this.update_logs.bind(this), 1000);
    },
    unmounted: function() {
        clearInterval(this.log_update_handle);
    },
    methods: {
        clear_logs: function() {
            store.linecount = 0;
            store.logs = "";
        },
        insert_logs: function(l) {
            if (l.out){
                store.logs += l.out;
                store.linecount += l.out.split("\n").length;
            }
            if (l.err){
                store.logs += l.err;
                store.linecount += l.err.split("\n").length;
            }
        },
        update_logs: function(){
            if (!this.live || !engine_store.status || !engine_store.status.running) {
                // skip log-fetching if not enabled or not live
                return;
            }
        
            // fetch and shift window
            const now = Date.now();
            get_logs(store.logs_live_since, now)
            .then(l => {
                this.insert_logs(l);
                store.logs_live_since = now + 1;
            });
        },
        fetch_logs: function() {
            this.clear_logs();
            this.busy = true;
            get_logs(+this.since, +this.until)
            .then(l => this.insert_logs(l))
            .finally(() => this.busy = false);
        }
    },
    template: /*html*/`
    <div class="card">
        <div class="card-header">
            Logs
            <span class="badge bg-primary">{{linecount}}</span>
            <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
        </div>
        <div class="card-body">
            <div class="row">
                <div class="col-2">
                    <label>From</label>
                    <vue-datepicker v-model="since"/>
                </div>
                <div class="col-2">
                    <label>Until</label>
                    <vue-datepicker v-model="until"/>
                </div>
                <div class="col-2">
                    <label>Live logs</label>
                    <button class="btn btn-outline-secondary w-100" v-on:click="live = true" v-if="!live">
                        Disabled
                    </button>
                    <button class="btn btn-success w-100" v-on:click="live = false" v-if="live">
                        Enabled
                    </button>
                </div>
                <div class="col-2">
                    <label>Log search</label>
                    <button class="btn btn-primary w-100" v-on:click="fetch_logs" :disabled="live">
                        Search
                    </button>
                </div>
                <div class="col-2">
                    <label>Clear logs</label>
                    <button class="btn btn-secondary w-100" v-on:click="clear_logs">
                        Clear
                    </button>
                </div>
            </div>
            <hr>
            <pre style="height: 400px; background: lightgray;">{{logs}}</pre>
        </div>
    </div>
	`
};
