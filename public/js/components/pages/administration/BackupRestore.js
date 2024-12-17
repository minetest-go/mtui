import { get_maintenance, get_stats } from "../../../service/stats.js";
import { enable_maintenance, disable_maintenance } from "../../../api/maintenance.js";
import { engine } from "../../../service/service.js";

import DefaultLayout from "../../layouts/DefaultLayout.js";
import { START } from "../../Breadcrumb.js";

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            breadcrumb: [START, {
                name: "Backup/Restore",
                icon: "upload",
                link: "/backup"
            }]
        };
    },
    computed: {
        maintenance: get_maintenance,
        is_engine_running: () => engine.is_running()
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
        <default-layout title="Backup/Restore" icon="upload" :breadcrumb="breadcrumb">
            <div class="card">
                <div class="card-header">
                    Download backup <i class="fa fa-download"></i>
                </div>
                <div class="card-body">
                    <a class="btn btn-primary">
                        Download backup
                    </a>
                </div>
            </div>
            <hr/>
            <div class="card">
                <div class="card-header">
                    Restore from backup <i class="fa fa-upload"></i>
                </div>
                <div class="card-body">
                    <div class="alert alert-warning" v-if="is_engine_running && !maintenance">
                        <i class="fa-solid fa-triangle-exclamation"></i>
                        The <router-link to="/services/engine">minetest engine</router-link> is still running, please stop it to enable the maintenance mode
                    </div>

                    <button class="btn btn-warning" v-if="!maintenance" v-on:click="enable_maintenance" :disabled="is_engine_running">
                        <i class="fa-solid fa-triangle-exclamation"></i>
                        Enable maintenance mode
                    </button>
                    <button class="btn btn-success" v-if="maintenance" v-on:click="disable_maintenance">
                        Disable maintenance mode
                    </button>
                    <i class="fa fa-chevron-right"></i>
                    <button class="btn btn-danger" :disabled="!maintenance">
                        Restore from backup
                    </button>
                </div>
            </div>
        </default-layout>
    `
};