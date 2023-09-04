import { get_maintenance, get_stats } from "../../service/stats.js";
import { enable_maintenance, disable_maintenance } from "../../api/maintenance.js";
import { check_features } from "../../service/features.js";

export default {
    computed: {
        maintenance: get_maintenance
    },
    methods: {
        enable_maintenance: function() {
            enable_maintenance()
            .then(() => get_stats());
        },
        disable_maintenance: function() {
            disable_maintenance()
            .then(() => Promise.all([get_stats(), check_features()]));
        }
    },
    template: /*html*/`
        <div>
            <h3>Maintenance</h3>
            Maintenance mode:
            <span class="badge bg-success" v-if="!maintenance">Disabled</span>
            <span class="badge bg-warning" v-if="maintenance">Enabled</span>
            <br>
            <button class="btn btn-warning" v-if="!maintenance" v-on:click="enable_maintenance">
                Enable maintenance mode
            </button>
            <button class="btn btn-success" v-if="maintenance" v-on:click="disable_maintenance">
                Disable maintenance mode
            </button>
        </div>
    `
};