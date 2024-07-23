import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";
import SkinPreview from "../SkinPreview.js";

import { get_servername, get_version } from "../../service/app_info.js";
import { has_priv, is_logged_in } from "../../service/login.js";
import { has_feature } from "../../service/features.js";
import { get_players } from "../../service/stats.js";

export default {
	components: {
		"default-layout": DefaultLayout,
		"skin-preview": SkinPreview
	},
	data: function() {
		return {
			breadcrumb: [START]
		};
	},
	methods: {
		has_priv,
		has_feature
	},
	computed: {
		is_logged_in,
		servername: get_servername,
		version: get_version,
		players: get_players
	},
	template: /*html*/`
	<default-layout icon="home" title="Start" :breadcrumb="breadcrumb">
		<div class="text-center">
			<h3>
				Minetest Web UI
				<small class="text-muted" v-if="servername">{{servername}}</small>
			</h3>
			<span v-if="version">
				Version: <span class="badge bg-primary">{{ version }}</span>
			</span>
			<hr/>
			<router-link to="/shell" class="btn btn-primary" v-if="has_priv('interact')">
				<i class="fa-solid fa-terminal"></i> Shell
			</router-link>
			&nbsp;
			<router-link to="/profile" class="btn btn-primary" v-if="is_logged_in">
				<i class="fa fa-user"></i> Profile
			</router-link>
			<router-link to="/login" class="btn btn-primary" v-else>
				<i class="fa-solid fa-right-to-bracket"></i> Login
			</router-link>
			&nbsp;
			<router-link to="/help" class="btn btn-primary" v-if="has_priv('server')">
				<i class="fa-solid fa-circle-question"></i> Help
			</router-link>
			&nbsp;
			<router-link to="/wizard/1" class="btn btn-primary" v-if="has_priv('server') && has_feature('docker') && has_feature('minetest_config') && has_feature('modmanagement')">
				<i class="fa-solid fa-wand-magic-sparkles"></i> Setup wizard
			</router-link>
			&nbsp;
			<a class="btn btn-secondary" href="https://github.com/minetest-go/mtui" target="new">
				<i class="fa-brands fa-github"></i> Source
			</a>
			<hr/>
			<h3 v-if="players && players.length">Online players</h3>
			<div class="container">
				<router-link
					class="btn btn-secondary m-1"
					:to="is_logged_in ? '/profile/'+player.name : ''"
					v-for="player in players"
					:key="player.name">
					{{player.name}}
					&nbsp;
					<skin-preview :playername="player.name"/>
				</router-link>
			</div>
		</div>
	</default-layout>
	`
};
