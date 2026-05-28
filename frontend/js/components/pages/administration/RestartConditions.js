import DefaultLayout from "../../layouts/DefaultLayout.js";
import { RESTART_CONDITIONS } from "../../Breadcrumb.js";
import { get_uimod_storage, set_uimod_storage } from "../../../api/uimod.js";

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            breadcrumb: [RESTART_CONDITIONS],
            busy: false,
            ready: false,
            // values fetched after load
            condition_empty: null,
            condition_mods_changed: null,
            max_uptime: null,
            restart_method: null
        };
    },
    mounted: function() {
        this.update();
    },
    methods: {
        update: async function() {
            this.condition_empty = await get_uimod_storage("restart_condition_empty") == "1";
            this.condition_mods_changed = await get_uimod_storage("restart_condition_mods_changed") == "1";
            this.max_uptime = parseInt(await get_uimod_storage("restart_condition_max_uptime")) || 0;
            this.restart_method = await get_uimod_storage("restart_method") || "DELAY";
            this.ready = true;
        },
        save: async function() {
            this.busy = true;
            await set_uimod_storage("restart_condition_empty", this.condition_empty ? "1" : "0");
            await set_uimod_storage("restart_condition_mods_changed", this.condition_mods_changed ? "1" : "0");
            await set_uimod_storage("restart_condition_max_uptime", "" + this.max_uptime);
            await set_uimod_storage("restart_method", this.restart_method);
            this.busy = false;
        }
    },
    template: /*html*/`
        <default-layout title="" icon="" :breadcrumb="breadcrumb">
            <h4>
                Configure restart conditions
                <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
            </h4>
            <div v-if="!ready">
                <i class="fa-solid fa-spinner fa-spin"></i>
                Loading configuration
            </div>
            <div v-else>
                <div class="form-check">
                    <input class="form-check-input" type="checkbox" v-model="condition_empty"/>
                    <label class="form-check-label">
                        Restart when the server is empty next time (condition resets after restart)
                    </label>
                </div>
                <div class="form-check">
                    <input class="form-check-input" type="checkbox" v-model="condition_mods_changed"/>
                    <label class="form-check-label">
                        Restart after mods have been updated
                    </label>
                </div>
                <div class="form-group">
                    <label class="form-control-label">
                            Restart after uptime, in seconds (0 = disabled)
                    </label>
                    <input class="form-control" type="number" min="0" step="0" v-model="max_uptime"/>
                </div>
                <div class="form-group">
                    <label class="form-control-label">
                            Restart method if the above conditions match
                    </label>
                    <select class="form-control" v-model="restart_method">
                        <option value="INSTANT">Instant, no delay</option>
                        <option value="DELAY">With delay and warning in chat</option>
                        <option value="EMPTY">Only if empty</option>
                    </select>
                </div>
                <button class="btn btn-primary" :disabled="busy" v-on:click="save">
                    <i class="fa fa-floppy-disk"></i>
                    Save
                </button>
            </div>
        </default-layout>
    `
};