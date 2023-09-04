import { get_maintenance, get_stats } from "../../service/stats.js";
import { enable_maintenance, disable_maintenance } from "../../api/maintenance.js";
import { check_features } from "../../service/features.js";
import { is_running } from "../../service/engine.js";
import events, { EVENT_LOGGED_IN} from "../../events.js";
import { get_claims } from "../../service/login.js";

export default {
    computed: {
        maintenance: get_maintenance,
        is_engine_running: is_running
    },
    methods: {
        enable_maintenance: function() {
            enable_maintenance()
            .then(() => get_stats());
        },
        disable_maintenance: function() {
            disable_maintenance()
            .then(() => Promise.all([get_stats(), check_features()]))
            .then(() => events.emit(EVENT_LOGGED_IN, get_claims()));
        }
    },
    template: /*html*/`
        <h3>Maintenance</h3>
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
            The maintenance mode shuts down any database-related service (ui and engine) in order to create and download consistent backups
        </div>
        <div class="alert alert-warning" v-if="is_engine_running && !maintenance">
            <i class="fa-solid fa-triangle-exclamation"></i>
            The <router-link to="/services/engine">minetest engine</router-link> is still running, please stop it to enable the maintenance mode
        </div>
    `
};