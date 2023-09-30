import { add, remove, get_all, is_busy, get_git_mod, update_mod, update_mod_version, check_updates, add_mtui } from '../../../service/mods.js';
import { update_settings } from '../../../service/mtconfig.js';

import FeedbackButton from '../../FeedbackButton.js';
import DefaultLayout from '../../layouts/DefaultLayout.js';
import CDBPackageLink from '../../CDBPackageLink.js';
import { START, ADMINISTRATION, MODS } from '../../Breadcrumb.js';

const ModRow = {
	props: ["mod", "busy"],
	components: {
		"cdb-package-link": CDBPackageLink
	},
	methods: {
		remove,
		toggle_autoupdate: function(mod) {
			mod.auto_update = !mod.auto_update;
			update_mod(mod);
		},
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
			<span class="badge bg-success" v-if="mod.source_type == 'cdb'">
				<i class="fa-solid fa-box-open"></i>
				ContentDB
			</span>
			<span class="badge bg-success" v-if="mod.source_type == 'git'">
				<i class="fa-brands fa-git-alt"></i>
				Git
			</span>
		</td>
		<td>
			<a :href="mod.url" v-if="mod.source_type == 'git'">{{mod.url}}</a>
		</td>
		<td>
			<span class="badge bg-secondary" v-if="mod.version">{{mod.version}}</span>
		</td>
		<td>
			<span class="badge bg-secondary" v-if="mod.latest_version">{{mod.latest_version}}</span>
		</td>
		<td>
			<button class="btn btn-success" v-if="mod.auto_update" v-on:click="toggle_autoupdate(mod)" :disabled="busy">
				<i class="fa fa-check"></i>
				Enabled
			</button>
			<button class="btn btn-secondary" v-if="!mod.auto_update" v-on:click="toggle_autoupdate(mod)" :disabled="busy">
				<i class="fa fa-times"></i>
				Disabled
			</button>
		</td>
		<td>
			<div class="btn-group">
				<button class="btn btn-primary" v-on:click="update_mod_version(mod, mod.latest_version)" :disabled="busy || mod.version == mod.latest_version">
					<i class="fa fa-download"></i>
					Update
				</button>
				<button class="btn btn-danger" v-on:click="remove(mod.id)" :disabled="busy">
					<i class="fa fa-trash"></i>
					Remove
				</button>
			</div>
		</td>
	`
};

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
			breadcrumb: [START, ADMINISTRATION, MODS]
		};
	},
	methods: {
		add: function() {
			return add({
				name: this.add_name,
				mod_type: this.add_mod_type,
				source_type: this.add_source_type,
				url: this.add_url,
				branch: "",
				version: this.add_version
			})
			.then(() => {
				this.add_name = "";
				this.add_url = "";
				this.add_version = "";
			});
		},
		add_mtui_mod: function() {
			add_mtui().then(update_settings);
		},
		update_mod_version: update_mod_version,
		get_git_mod: get_git_mod,
		check_updates: check_updates
	},
	computed: {
		busy: is_busy,
		games: () => get_all().filter(m => m.mod_type == "game"),
		mods: () => get_all().filter(m => m.mod_type == "mod"),
		txps: () => get_all().filter(m => m.mod_type == "txp")
	},
	template: /*html*/`
		<default-layout icon="cubes" title="Mods" :breadcrumb="breadcrumb">
			<h4>
				Mod management
				<i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
			</h4>
			<div class="alert alert-warning" v-if="!get_git_mod('mtui')">
				<div class="row">
					<div class="col-12">
						<i class="fa-solid fa-triangle-exclamation"></i>
						<b>Warning:</b>
						The <i>mtui</i> mod is not installed, some features may not work properly
						<button class="btn btn-primary float-end" :disabled="busy" v-on:click="add_mtui_mod">
							<i class="fa fa-plus"></i>
							Install "mtui" mod
						</button>
					</div>
				</div>
			</div>
			<div class="row">
				<div class="col-10"></div>
				<div class="col-2">
					<button class="btn btn-primary w-100" v-on:click="check_updates" :disabled="busy">
						<i class="fa fa-refresh"></i>
						Check for updates
					</button>
				</div>
			</div>
			<table class="table table-condensed table-striped">
				<thead>
					<tr>
						<th>Type</th>
						<th>Name</th>
						<th>Source-Type</th>
						<th>Source</th>
						<th>Version</th>
						<th>Latest Version</th>
						<th>Auto-update</th>
						<th>Actions</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>
							<select class="form-control" v-model="add_mod_type">
								<option value="mod">Mod</option>
								<option value="game">Game</option>
								<option value="txp">Texturepack</option>
							</select>
						</td>
						<td>
							<input class="form-control" type="text" placeholder="Mod name" v-model="add_name"/>
						</td>
						<td>
							<span class="badge bg-success">
								<i class="fa-brands fa-git-alt"></i>
								Git
							</span>
						</td>
						<td>
							<input class="form-control" type="text" placeholder="Source url" v-model="add_url"/>
						</td>
						<td>
							<input class="form-control" type="text" placeholder="Version" v-model="add_version"/>
						</td>
						<td></td>
						<td></td>
						<td>
							<feedback-button type="success" :fn="add" :disabled="!add_name || !add_url">
								<i class="fa-brands fa-git-alt"></i>
								Add from git
							</feedback-button>
						</td>
					</tr>
					<tr>
						<td></td>
						<td></td>
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
						<mod-row :mod="mod" :busy="busy"/>
					</tr>
					<tr v-if="txps.length > 0" class="table-secondary">
						<td colspan="8">
							<h4>Texture-packs</h4>
						</td>
					</tr>
					<tr v-for="mod in txps" :key="mod.id">
						<mod-row :mod="mod" :busy="busy"/>
					</tr>
					<tr v-if="mods.length > 0" class="table-secondary">
						<td colspan="8">
							<h4>Mods</h4>
						</td>
					</tr>
					<tr v-for="mod in mods" :key="mod.id">
						<mod-row :mod="mod" :busy="busy"/>
					</tr>
				</tbody>
			</table>
		</default-layout>
	`
};
