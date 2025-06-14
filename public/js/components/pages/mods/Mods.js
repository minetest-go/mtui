import {
	add, remove,
	get_mods_by_type,
	get_git_mod,
	update_mod,
	update_mod_version,
	check_updates,
	add_mtui,
	update
} from '../../../service/mods.js';
import { update_settings } from '../../../service/mtconfig.js';

import FeedbackButton from '../../FeedbackButton.js';
import DefaultLayout from '../../layouts/DefaultLayout.js';
import CDBPackageLink from '../../CDBPackageLink.js';
import { START, ADMINISTRATION, MODS } from '../../Breadcrumb.js';
import { has_feature } from '../../../service/features.js';

const ModRow = {
	props: ["mod"],
	components: {
		"cdb-package-link": CDBPackageLink
	},
	methods: {
		remove,
		update_mod_version,
		update_mod
	},
	template: /*html*/`
		<td>
			<span class="badge bg-secondary">{{mod.mod_type}}</span>
		</td>
		<td>
			<cdb-package-link v-if="mod.source_type == 'cdb'" :author="mod.author" :name="mod.name"/>
			<span v-else>{{mod.name}}</span>
		</td>
		<td>
			<a :href="mod.url" v-if="mod.source_type == 'git'">{{mod.url}}</a>
			<span class="badge bg-success" v-if="mod.source_type == 'cdb'">
				<i class="fa-solid fa-box-open"></i>
				ContentDB
			</span>
			<div v-if="mod.mod_status == 'error'" class="badge bg-danger">
				{{mod.message}}
			</div>
		</td>
		<td>
			<span class="badge bg-secondary" v-if="mod.version">{{mod.version}}</span>
			<span class="badge bg-warning" v-if="mod.latest_version && mod.latest_version != mod.version">Latest: {{mod.latest_version}}</span>
		</td>
		<td>
			<div class="form-check" v-if="mod.mod_status == 'installed'">
				<input class="form-check-input" type="checkbox" v-model="mod.auto_update" v-on:change="update_mod(mod)"/>
				<label class="form-check-label">Auto-update</label>
			</div>
		</td>
		<td>
			<div class="btn-group">
				<button class="btn btn-primary" v-on:click="update_mod_version(mod, mod.latest_version)" :disabled="mod.version == mod.latest_version" v-if="mod.mod_status == 'installed'">
					<i class="fa fa-download"></i>
					Update
				</button>
				<button class="btn btn-danger" v-on:click="remove(mod.id)" v-if="mod.mod_status != 'processing'">
					<i class="fa fa-trash"></i>
					Remove
				</button>
			</div>
			<i class="fa fa-spin fa-spinner" v-if="mod.mod_status == 'processing'"></i>
		</td>
	`
};

const modname_regex = /^[a-zA-Z0-9_]*$/;

export default {
	components: {
		"feedback-button": FeedbackButton,
		"default-layout": DefaultLayout,
		"mod-row": ModRow
	},
	data: () => {
		return {
			add_name: "",
			add_mod_type: "mod",
			add_source_type: "git",
			add_url: "",
			add_version: "",
			breadcrumb: [START, ADMINISTRATION, MODS],
			update_handle: null,
			busy: false
		};
	},
	created: function() {
        this.update_handle = setInterval(update, 2000);
    },
    unmounted: function() {
        clearInterval(this.update_handle);
    },
	methods: {
		add: async function() {
			await add({
				name: this.add_name,
				mod_type: this.add_mod_type,
				source_type: this.add_source_type,
				url: this.add_url,
				branch: "",
				version: this.add_version
			});
			this.add_name = "";
			this.add_url = "";
			this.add_version = "";
		},
		add_mtui_mod: async function() {
			await add_mtui();
			if (has_feature("minetest_config")) {
				await update_settings();
			}
		},
		get_git_mod,
		check_updates: async function() {
			this.busy = true;
			await check_updates();
			this.busy = false;
		}
	},
	computed: {
		games: () => get_mods_by_type("game"),
		mods: () => get_mods_by_type("mod"),
		worldmods: () => get_mods_by_type("worldmods"),
		txps: () => get_mods_by_type("txp"),
		add_name_valid: function() {
			return (this.add_name == "" || modname_regex.test(this.add_name));
		}
	},
	template: /*html*/`
		<default-layout icon="cubes" title="Mods" :breadcrumb="breadcrumb">
			<h4>
				Mod management
			</h4>
			<div class="alert alert-warning" v-if="!get_git_mod('mtui')">
				<div class="row">
					<div class="col-12">
						<i class="fa-solid fa-triangle-exclamation"></i>
						<b>Warning:</b>
						The <i>mtui</i> mod is not installed, some features may not work properly
						<button class="btn btn-primary float-end" :disabled="busy" v-on:click="add_mtui_mod">
							<i class="fa fa-plus"></i>
							Install the "mtui" mod
							<i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
						</button>
					</div>
				</div>
			</div>
			<div class="row">
				<div class="col-md-10"></div>
				<div class="col-md-2">
					<button class="btn btn-primary w-100" v-on:click="check_updates" :disabled="busy">
						<i class="fa fa-refresh" v-bind:class="{'fa-spin': busy}"></i>
						Check for updates
					</button>
				</div>
			</div>
			<table class="table table-condensed table-striped">
				<thead>
					<tr>
						<th>Type</th>
						<th>Name</th>
						<th>Source</th>
						<th>Version</th>
						<th>Auto-update</th>
						<th>Actions</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>
							<select class="form-control" v-model="add_mod_type">
								<option value="mod">Mod</option>
								<option value="worldmods">Worldmods</option>
								<option value="game">Game</option>
								<option value="txp">Texturepack</option>
							</select>
						</td>
						<td>
							<div class="input-group has-validation">
								<input
									class="form-control"
									type="text"
									placeholder="Mod name"
									v-model="add_name"
									v-bind:class="{'is-invalid': !add_name_valid}"/>
								<div class="invalid-feedback" v-if="!add_name_valid">
									Allowed modname characters: a-z A-Z 0-9 _
								</div>
							</div>
						</td>
						<td>
							<input class="form-control" type="text" placeholder="Source url" v-model="add_url"/>
						</td>
						<td>
							<input class="form-control" type="text" placeholder="Version" v-model="add_version"/>
						</td>
						<td></td>
						<td>
							<feedback-button type="success" :fn="add" :disabled="!add_name || !add_url">
								<i class="fa-brands fa-git-alt"></i>
								Add from git
							</feedback-button>
						</td>
					</tr>
					<tr>
						<td>
							<span class="badge bg-success">
								<i class="fa-solid fa-box-open"></i>
								ContentDB
							</span>
						</td>
						<td></td>
						<td></td>
						<td></td>
						<td></td>
						<td>
							<router-link to="/cdb/browse" class="btn btn-success">
								<i class="fa-solid fa-box-open"></i>
								Add from ContentDB
							</router-link>
						</td>
					</tr>
					<tr v-if="games.length > 0" class="table-secondary">
						<td colspan="8">
							<h4>Game</h4>
						</td>
					</tr>
					<tr v-for="mod in games" :key="mod.id">
						<mod-row :mod="mod"/>
					</tr>
					<tr v-if="txps.length > 0" class="table-secondary">
						<td colspan="8">
							<h4>Texture-packs</h4>
						</td>
					</tr>
					<tr v-for="mod in txps" :key="mod.id">
						<mod-row :mod="mod"/>
					</tr>
					<tr v-if="mods.length > 0" class="table-secondary">
						<td colspan="8">
							<h4>Mods</h4>
						</td>
					</tr>
					<tr v-for="mod in mods" :key="mod.id">
						<mod-row :mod="mod"/>
					</tr>
					<tr v-if="worldmods.length > 0" class="table-secondary">
						<td colspan="8">
							<h4>Worldmods collection</h4>
						</td>
					</tr>
					<tr v-for="worldmod in worldmods" :key="worldmod.id">
						<mod-row :mod="worldmod"/>
					</tr>
				</tbody>
			</table>
		</default-layout>
	`
};
