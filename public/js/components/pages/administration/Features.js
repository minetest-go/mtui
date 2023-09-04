import { set_feature, store } from '../../../service/features.js';
import DefaultLayout from '../../layouts/DefaultLayout.js';
import { START, ADMINISTRATION } from '../../Breadcrumb.js';

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            features: store,
            breadcrumb: [START, ADMINISTRATION, {
                name: "Features",
                icon: "tags",
                link: "/features"
            }]
        };
    },
    methods: {
        getDescription: function(name) {
            const h = store[name];
            return h.description || "<no description>";
        },
        getMods: function(name) {
            const h = store[name];
            return h.mods || [];
        },
        set_feature: function(name, enabled) {
            set_feature(name, enabled);
        },
        is_experimental: function(name) {
            const h = store[name];
            return h.experimental;
        }
    },
    template: /*html*/`
    <default-layout icon="tags" title="Features" :breadcrumb="breadcrumb">
        <table class="table table-condensed table-striped">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Status</th>
                    <th>Action</th>
                    <th>Description</th>
                    <th>Required mods</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(feature, name) in features">
                    <td>
                        {{name}}
                        <i class="fa-solid fa-flask" v-if="is_experimental(name)" title="Experimental feature"></i>
                    </td>
                    <td v-if="feature.enabled">
                        <i class="fa-solid fa-check" style="color: green;"></i>
                    </td>
                    <td v-if="!feature.enabled">
                        <i class="fa-solid fa-times" style="color: red;"></i>
                    </td>
                    <td v-if="feature.enabled">
                        <button class="btn btn-sm btn-danger" v-on:click="set_feature(name, false)">
                            Disable
                        </button>
                    </td>
                    <td v-if="!feature.enabled">
                        <button class="btn btn-sm btn-primary" v-on:click="set_feature(name, true)">
                            Enable
                        </button>
                    </td>
                    <td>{{ getDescription(name) }}</td>
                    <td>
                        <span class="badge bg-success" v-for="mod in getMods(name)">
                            {{mod}}
                        </span>
                    </td>
                </tr>
            </tbody>
        </table>
    </default-layout>
    `
};