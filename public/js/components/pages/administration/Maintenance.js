import { get_maintenance, get_stats } from "../../../service/stats.js";
import { enable_maintenance, disable_maintenance } from "../../../api/maintenance.js";
import { engine } from "../../../service/service.js";
import DefaultLayout from "../../layouts/DefaultLayout.js";
import { START, ADMINISTRATION } from "../../Breadcrumb.js";

export default {
    data: function() {
        return {
            breadcrumb: [START, ADMINISTRATION, {
                icon: "wrench",
                name: "Maintenance",
                link: "/maintenance"
            }]
        };
    },
    components: {
        "default-layout": DefaultLayout
    },
    computed: {
        maintenance: get_maintenance,
        is_engine_running: () => engine.is_running
    },
    methods: {
        enable_maintenance: function() {
            enable_maintenance()
            .then(() => get_stats());
        },
        disable_maintenance: function() {
            disable_maintenance()
            .then(() => window.location.reload());
        }
    },
    template: /*html*/`
    <default-layout icon="wrench" title="Maintenance" :breadcrumb="breadcrumb">
        <h4>Maintenance</h4>
        <table class="table">
            <tbody>
                <tr>
                    <td>Maintenance mode</td>
                    <td>
                        <span class="badge bg-success" v-if="!maintenance">Disabled</span>
                        <span class="badge bg-warning" v-if="maintenance">Enabled</span>
                    </td>
                </tr>
                <tr>
                    <td>Actions</td>
                    <td>
                        <button class="btn btn-warning" v-if="!maintenance" v-on:click="enable_maintenance" :disabled="is_engine_running">
                            <i class="fa-solid fa-triangle-exclamation"></i>
                            Enable maintenance mode
                        </button>
                        <button class="btn btn-success" v-if="maintenance" v-on:click="disable_maintenance">
                            Disable maintenance mode
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
        <div class="alert alert-info">
            <i class="fa fa-info"></i>
            The maintenance mode shuts down any database access in order to create and download consistent backups
        </div>
        <div class="alert alert-warning" v-if="is_engine_running && !maintenance">
            <i class="fa-solid fa-triangle-exclamation"></i>
            The <router-link to="/services/engine">minetest engine</router-link> is still running, please stop it to enable the maintenance mode
        </div>
    </default-layout>
    `
};