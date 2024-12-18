import { get_maintenance, get_stats } from "../../../service/stats.js";
import { enable_maintenance, disable_maintenance } from "../../../api/maintenance.js";
import { engine } from "../../../service/service.js";
import { upload_chunked } from "../../../service/uploader.js";
import { unzip, remove } from "../../../api/filebrowser.js";
import { get_mods_by_type, remove as remove_mod, add_mtui, add_beerchat, add_mapserver } from "../../../service/mods.js";

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
            }],
            restore_active: false,
            restore_message: "",
            restore_progress: 0
        };
    },
    computed: {
        maintenance: get_maintenance,
        is_engine_running: () => engine.is_running()
    },
    methods: {
        enable_maintenance: async function() {
            await enable_maintenance();
            get_stats();
        },
        disable_maintenance: async function() {
            await disable_maintenance();
            window.location.reload();
        },
        restore: async function() {
            const file = this.$refs.input_upload.files[0];
            if (!file) {
                // no file selected
                return;
            }

            this.restore_message = "Starting to upload archive";
            this.restore_active = true;

            await upload_chunked("/", "restore.zip", file, progress => {
                this.restore_progress = progress;
                this.restore_message = `Uploading: ${Math.floor(progress*100)}% done`;
            });

            this.restore_message = "Unzipping archive...";
            await unzip("/restore.zip");

            this.restore_message = "Removing temporary archive";
            await remove("/restore.zip");

            this.restore_active = false;
            
            this.restore_message = "Disableing maintenance mode";
            await disable_maintenance();

            this.restore_message = "Reconfiguring system-relevant mods";
            const mods = get_mods_by_type("mod");

            for (let i=0; i<mods.length; i++) {
                const mod = mods[i];
                // remove and reinstall
                if (mod.name == "mtui"){
                    await remove_mod(mod.id);
                    await add_mtui();
                }
                if (mod.name == "mapserver"){
                    await remove_mod(mod.id);
                    await add_mapserver();
                }
                if (mod.name == "beerchat"){
                    await remove_mod(mod.id);
                    await add_beerchat();
                }
            }

            window.location.reload();
        }
    },
    template: /*html*/`
        <default-layout title="Backup/Restore" icon="upload" :breadcrumb="breadcrumb">
            <div class="row">
                <div class="col-md-6">
                    <div class="card">
                        <div class="card-header">
                            Download backup <i class="fa fa-download"></i>
                        </div>
                        <div class="card-body">
                            <a class="btn btn-primary" href="api/filebrowser/zip?dir=/">
                                <i class="fa fa-file-zipper"></i>
                                Download world-backup as zip-file
                            </a>
                        </div>
                    </div>
                </div>
                <div class="col-md-6">
                    <div class="card">
                        <div class="card-header">
                            Restore from backup <i class="fa fa-upload"></i>
                        </div>
                        <div class="card-body">
                            <div class="alert alert-info" v-if="!maintenance">
                                <i class="fa-solid fa-info"></i>
                                The maintenance mode must be enabled and all the services stopped to restore from a backup
                            </div>

                            <div class="alert alert-warning" v-if="is_engine_running && !maintenance">
                                <i class="fa-solid fa-triangle-exclamation"></i>
                                The <router-link to="/services/engine">minetest engine</router-link> is still running, please stop it to enable the maintenance mode
                            </div>

                            <div class="alert alert-warning" v-if="!is_engine_running && maintenance">
                                <i class="fa-solid fa-triangle-exclamation"></i>
                                <b>Warning:</b> All existing world-data will be overwritten by a backup-restore!
                            </div>

                            <div class="input-group" v-if="!restore_active">
                                <button class="btn btn-warning" v-if="!maintenance" v-on:click="enable_maintenance" :disabled="is_engine_running">
                                    <i class="fa-solid fa-triangle-exclamation"></i>
                                    Enable maintenance mode
                                </button>
                                <button class="btn btn-success" v-if="maintenance" v-on:click="disable_maintenance">
                                    Disable maintenance mode
                                </button>
                                <input ref="input_upload" type="file" class="form-control" :disabled="!maintenance" accept=".zip"/>
                                <button class="btn btn-danger" v-on:click="restore" :disabled="!maintenance">
                                    <i class="fa fa-file-zipper"></i>
                                    Restore from zipfile backup
                                </button>
                            </div>

                            <div class="progress" v-if="restore_active">
                                <div class="progress-bar overflow-visible progress-bar-striped progress-bar-animated" v-bind:style="{ width: (restore_progress*100)+'%' }">
                                    {{restore_message}}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </default-layout>
    `
};