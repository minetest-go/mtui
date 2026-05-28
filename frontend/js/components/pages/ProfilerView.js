import DefaultLayout from "../layouts/DefaultLayout.js";
import { START, FILEBROWSER } from "../Breadcrumb.js";

import { download_text } from "../../api/filebrowser.js";

import format_count from "../../util/format_count.js";

const InstrumentRow = {
    props: ["instrument", "name"],
    computed: {
        avg_micros: function() {
            return Math.floor(this.instrument.time_all / this.instrument.samples);
        },
        tableClass: function() {
            if (this.avg_micros > 50000) {
                return { "table-danger": true };
            } else if (this.avg_micros > 25000) {
                return { "table-warning": true };
            }
        }
    },
    methods: {
        format_count
    },
    template: /*html*/`
    <tr v-bind:class="tableClass">
        <td>
            <span class="badge bg-secondary">{{name}}</span>
        </td>
        <td>{{instrument.time_min}}</td>
        <td>{{avg_micros}}</td>
        <td>{{instrument.time_max}}</td>
        <td>{{ format_count(instrument.samples) }}</td>
    </tr>
    `
};

export default {
    props: ["pathMatch"],
    components: {
        "default-layout": DefaultLayout,
        "instrument-row": InstrumentRow
    },
    data: function() {
        return {
            breadcrumb: [START, FILEBROWSER, {
                name: "View profile",
                icon: "chart-line",
                link: "/profiler-view/" + this.pathMatch
            }],
            busy: false,
            profile: {}
        };
    },
    mounted: function() {
        this.busy = true;
        download_text(this.pathMatch)
        .then(p => {
            this.profile = JSON.parse(p);
            this.busy = false;
        });
    },
    template: /*html*/`
        <default-layout title="View profile" icon="chart-line" :breadcrumb="breadcrumb">
            <h4>
                View profile
                <small class="muted">{{pathMatch}}</small>
                <i class="fa fa-spinner fa-spin" v-if="busy"></i>
            </h4>
            <hr>
            <div v-for="(entry, modname) in profile">
                <h5>{{modname}}</h5>
                <table class="table table-condensed table-striped">
                    <thead>
                        <tr>
                            <th style="width: 20%;">Instrument</th>
                            <th style="width: 20%;">Min time [&micro;s]</th>
                            <th style="width: 20%;">Avg time [&micro;s]</th>
                            <th style="width: 20%;">Max time [&micro;s]</th>
                            <th style="width: 20%;">Samples</th>
                        </tr>
                    </thead>
                    <tbody>
                        <instrument-row v-for="(instrument, name) in entry.instruments" :instrument="instrument" :name="name" :key="name"/>
                    </tbody>
                </table>
            </div>
        </default-layout>
    `
};