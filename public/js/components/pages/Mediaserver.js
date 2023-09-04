import { scan, stats } from "../../api/media.js";
import format_count from "../../util/format_count.js";
import format_size from "../../util/format_size.js";
import DefaultLayout from "../layouts/DefaultLayout.js";
import { START, ADMINISTRATION } from "../Breadcrumb.js";

export default {
    data: function() {
        return {
            busy: false,
            stats: null,
            breadcrumb: [START, ADMINISTRATION, {
                icon: "photo-film",
                name: "Media server",
                link: "/mediaserver"
            }]
        };
    },
    components: {
        "default-layout": DefaultLayout
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
        },
        media_url: function() {
            return location.protocol + "//" + location.host + location.pathname + "api/media/";
        }
    },
    template: /*html*/`
    <default-layout icon="photo-film" title="Media server" :breadcrumb="breadcrumb">
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
        <div class="alert alert-info">
            URL for the "remote_media" setting:
            <b>{{ media_url() }}</b>
        </div>
    </default-layout>
    `
};