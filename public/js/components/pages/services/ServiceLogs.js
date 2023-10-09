import { get_logs } from "../../../api/service.js";

export default {
    props: ["running", "servicename"],
    data: function(){
        return {
            busy: false,
            logs: "",
            linecount: 0,
            live: true,
            since: Date.now() - (1000*60*60),
            until: Date.now() + (1000*60*60),
            logs_live_since: Date.now() - (1000*60*60) // shifting window for live-view
        };
    },
    created: function(){
        this.update_logs();
        this.log_update_handle = setInterval(this.update_logs.bind(this), 1000);
    },
    unmounted: function() {
        clearInterval(this.log_update_handle);
    },
    methods: {
        clear_logs: function() {
            this.linecount = 0;
            this.logs = "";
        },
        insert_logs: function(l) {
            if (l.out){
                this.logs += l.out;
                this.linecount += l.out.split("\n").length;
            }
            if (l.err){
                this.logs += l.err;
                this.linecount += l.err.split("\n").length;
            }
            this.scroll_to_bottom();
        },
        scroll_to_bottom: function() {
            this.$refs.log_pre.scrollTop = this.$refs.log_pre.scrollTopMax;
        },
        update_logs: function(){
            if (!this.live || !this.running) {
                // skip log-fetching if not enabled or not live
                return;
            }
        
            // fetch and shift window
            const now = Date.now();
            get_logs(this.servicename, this.logs_live_since, now)
            .then(l => {
                this.insert_logs(l);
                this.logs_live_since = now + 1;
            });
        },
        fetch_logs: function() {
            this.clear_logs();
            this.busy = true;
            get_logs(this.servicename, +this.since, +this.until)
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
                        <i class="fa fa-pause"></i>
                    </button>
                    <button class="btn btn-success w-100" v-on:click="live = false" v-if="live">
                        Enabled
                        <i class="fa fa-play"></i>
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
            <pre ref="log_pre" style="height: 400px; background: gray;">{{logs}}</pre>
        </div>
    </div>
	`
};
