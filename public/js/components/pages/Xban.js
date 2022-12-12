
import { get_status, cleanup } from "../../api/xban.js";

export default {
    data: function() {
        return {
            status: null,
            busy: false,
            cleanup_result: null
        };
    },
    mounted: function() {
        this.update_status();
    },
    methods: {
        update_status: function(){
            get_status().then(s => this.status = s);
        },
        cleanup: function() {
            this.busy = true;
            this.cleanup_result = null;
            cleanup().then(r => {
                this.busy = false;
                this.cleanup_result = r;
            });
        }
    },
    template: /*html*/`
    <div>
        <h4>XBan database management</h4>
        <ul v-if="status">
            <li>
                Total entries: {{status.total}}
            </li>
            <li>
                Banned entries: {{status.banned}}
            </li>
        </ul>
        <a class="btn btn-warning" :disabled="busy" v-on:click="cleanup">
            <i class="fa-solid fa-broom" v-if="!busy"></i>
            <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
            Remove unbanned entries
        </a>
        <div class="alert alert-success" v-if="cleanup_result">
            <b>Cleanup results:</b>
            Removed {{cleanup_result.removed}} and retained {{cleanup_result.retained}} entries
        </div>
    </div>
    `
};
