import { START, SERVICES } from "../../Breadcrumb.js";
import { add_beerchat, get_git_mod } from "../../../service/mods.js";
import { update_settings } from '../../../service/mtconfig.js';

import DefaultLayout from "../../layouts/DefaultLayout.js";
import ServiceView from "./ServiceView.js";

export default {
	components: {
		"default-layout": DefaultLayout,
		"service-view": ServiceView
	},
	computed: {
		breadcrumb: function() {
			return [START, SERVICES, {
				name: "Matterbridge",
				icon: "gear",
				link: "/services/matterbridge"
			}];
		}
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
                    <button class="btn btn-primary float-end" v-on:click="add_beerchat_mod">
                        <i class="fa fa-plus"></i>
                        Install and configure the "beerchat" mod
                    </button>
                </div>
            </div>
        </div>
		<service-view servicename="matterbridge"/>
	</default-layout>
	`
};
