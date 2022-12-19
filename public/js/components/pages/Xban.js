
import { get_status, cleanup, get_records } from "../../api/xban.js";
import format_time from "../../util/format_time.js";

const store = Vue.reactive({
    status: null,
    busy: false,
    cleanup_result: null,
    banned_records: []
});

export default {
    data: () => store,
    mounted: function() {
        if (this.banned_records.length == 0){
            this.update();
        }
    },
    methods: {
        format_time: format_time,
        update: function(){
            this.busy = true;
            get_status()
            .then(s => this.status = s)
            .then(() => get_records())
            .then(r => this.banned_records = r)
            .then(() => this.busy = false);
        },
        cleanup: function() {
            this.busy = true;
            this.cleanup_result = null;
            cleanup().then(r => {
                this.busy = false;
                this.cleanup_result = r;
                this.update();
            });
        }
    },
    template: /*html*/`
    <div>
        <div class="row">
            <div class="col-md-10">
                <h3>
                    XBan database management
                    <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                </h3>
            </div>
            <div class="col-md-2 btn-group">
                <a class="btn btn-success" v-on:click="update">
                    <i class="fa-solid fa-rotate"></i>
                    Update
                </a>
            </div>
        </div>
        <ul v-if="status">
            <li>
                Total entries: {{status.total}}
            </li>
            <li>
                Banned entries: {{status.banned}}
            </li>
        </ul>
        <a class="btn btn-warning" :disabled="busy" v-on:click="cleanup">
            <i class="fa-solid fa-broom"></i>
            Remove unbanned entries
        </a>
        <div class="alert alert-success" v-if="cleanup_result">
            <b>Cleanup results:</b>
            Removed {{cleanup_result.removed}} and retained {{cleanup_result.retained}} entries
        </div>
        <hr>
        <h4>Ban entries</h4>
        <table class="table table-striped table-condensed" v-if="banned_records">
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Names</th>
                    <th>Records</th>
                    <th>Last seen</th>
                    <th>Banned until</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="record, i in banned_records">
                    <td>
                        <span class="badge bg-primary">{{i+1}}</span>
                    </td>
                    <td>
                        <ul>
                            <li v-for="_, name in record.names">
                                <router-link :to="'/profile/' + name">
                                    {{name}}
                                </router-link>
                            </li>
                        </ul>
                    </td>
                    <td>
                        <ul>
                            <li v-for="entry in record.record">
                                <b>Source:</b> <router-link :to="'/profile/' + entry.source">{{entry.source}}</router-link>
                                &nbsp;
                                <b>Reason:</b> {{entry.reason}}
                                &nbsp;
                                <b>Time:</b> {{format_time(entry.time)}}
                            </li>
                        </ul>
                    </td>
                    <td>{{format_time(record.last_seen)}}</td>
                    <td>{{format_time(record.time)}}</td>
                </tr>
            </tbody>
        </table>
    </div>
    `
};
