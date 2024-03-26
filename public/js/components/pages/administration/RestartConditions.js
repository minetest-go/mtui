import DefaultLayout from "../../layouts/DefaultLayout.js";
import { RESTART_CONDITIONS } from "../../Breadcrumb.js";
import { get_uimod_storage, set_uimod_storage } from "../../../api/uimod_storage.js";

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            busy: false,
            ready: false,
            condition_empty: false,
            breadcrumb: [RESTART_CONDITIONS]
        };
    },
    mounted: function() {
        this.update();
    },
    methods: {
        update: function() {
            get_uimod_storage("restart_condition_empty").then(value => {
                this.condition_empty = value == "1";
                this.ready = true;
            });
        },
        set_condition_empty: function() {
            this.busy = true;
            set_uimod_storage("restart_condition_empty", this.condition_empty ? "1" : "0")
            .then(() => this.busy = false);
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
            </div>
            <div v-else>
                <div class="form-check">
                    <input class="form-check-input" type="checkbox" v-model="condition_empty" v-on:change="set_condition_empty">
                    <label class="form-check-label">
                        restart if the server is empty
                    </label>
                </div>
            </div>
        </default-layout>
    `
};