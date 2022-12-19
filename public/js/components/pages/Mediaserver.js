import { scan, stats } from "../../api/media.js";
import format_count from "../../util/format_count.js";
import format_size from "../../util/format_size.js";

export default {
    data: function() {
        return {
            busy: false,
            stats: null
        };
    },
    mounted: function() {
        this.update();
    },
    methods: {
        format_count: format_count,
        format_size: format_size,
        scan: function() {
            this.busy = true;
            scan()
            .then(() => {
                this.busy = false;
                this.update();
            });
        },
        update: function() {
            stats().then(s => this.stats = s);
        }
    },
    template: /*html*/`
    <div>
        <button class="btn btn-primary" :disabled="busy" v-on:click="scan">
            <i class="fa fa-refresh"/>
            Scan
        </button>
        <table class="table table-condensed table-striped" v-if="stats">
            <tbody>
                <tr>
                    <td>Total size</td>
                    <td>{{format_size(stats.size)}}</td>
                </tr>
                <tr>
                    <td>Total count</td>
                    <td>{{format_count(stats.count)}}</td>
                </tr>
                <tr>
                    <td>Transferred</td>
                    <td>{{format_size(stats.transferredbytes)}}</td>
                </tr>
            </tbody>
        </table>
    </div>
    `
};