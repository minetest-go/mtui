import { scan, stats } from "../../api/media.js";
import prettycount from "../../util/prettycount.js";
import prettysize from "../../util/prettysize.js";

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
        prettycount: prettycount,
        prettysize: prettysize,
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
                    <td>{{prettysize(stats.size)}}</td>
                </tr>
                <tr>
                    <td>Total count</td>
                    <td>{{prettycount(stats.count)}}</td>
                </tr>
                <tr>
                    <td>Transferred</td>
                    <td>{{prettysize(stats.transferredbytes)}}</td>
                </tr>
            </tbody>
        </table>
    </div>
    `
};