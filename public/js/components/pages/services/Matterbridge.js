import { START, SERVICES } from "../../Breadcrumb.js";
import { add_beerchat, get_git_mod } from "../../../service/mods.js";
import { update_settings } from '../../../service/mtconfig.js';
import { is_busy } from "../../../service/mods.js";

import DefaultLayout from "../../layouts/DefaultLayout.js";
import ServiceView from "./ServiceView.js";
import FileEditor from "../filebrowser/FileEditor.js";

export default {
	components: {
		"default-layout": DefaultLayout,
		"service-view": ServiceView,
        "file-editor": FileEditor
	},
	computed: {
		breadcrumb: function() {
			return [START, SERVICES, {
				name: "Matterbridge",
				icon: "gear",
				link: "/services/matterbridge"
			}];
		},
        busy: is_busy
	},
    methods: {
        get_git_mod,
        add_beerchat_mod: function() {
            add_beerchat().then(update_settings);
        }
    },
	template: /*html*/`
	<default-layout icon="gear" title="Matterbridge" :breadcrumb="breadcrumb">
        <div class="alert alert-warning" v-if="!get_git_mod('beerchat')">
            <div class="row">
                <div class="col-12">
                    <i class="fa-solid fa-triangle-exclamation"></i>
                    <b>Warning:</b>
                    The <i>beerchat</i> mod is not installed, some features may not work properly
                    <button class="btn btn-primary float-end" :disabled="busy" v-on:click="add_beerchat_mod">
                        <i class="fa fa-plus"></i>
                        Install and configure the "beerchat" mod
                        <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                    </button>
                </div>
            </div>
        </div>
        <div class="row">
            <div class="col-6">
		        <service-view servicename="matterbridge"/>
            </div>
            <div class="col-6">
                <file-editor filename="/matterbridge.toml"/>
            </div>
        </div>
	</default-layout>
	`
};
