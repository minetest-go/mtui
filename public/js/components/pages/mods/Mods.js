import { add, remove, get_all, is_busy, get_git_mod } from '../../../service/mods.js';
import FeedbackButton from '../../FeedbackButton.js';
import DefaultLayout from '../../layouts/DefaultLayout.js';
import CDBPackageLink from '../../CDBPackageLink.js';
import { START, ADMINISTRATION, MODS } from '../../Breadcrumb.js';

export default {
	components: {
		"feedback-button": FeedbackButton,
		"default-layout": DefaultLayout,
		"cdb-package-link": CDBPackageLink
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
			add({
				name: "mtui",
				mod_type: "mod",
				source_type: "git",
				url: "https://github.com/minetest-go/mtui_mod.git",
				branch: "refs/heads/master"
			});
		},
		remove: remove,
		get_mods: get_all,
		get_git_mod: get_git_mod
	},
	computed: {
		busy: is_busy
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
			<table class="table table-condensed table-striped">
				<thead>
					<tr>
						<th>Type</th>
						<th>Name</th>
						<th>Source-Type</th>
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
								<option value="game">Game</option>
								<option value="txp" v-if="false">Textures</option>
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
						<td>
						</td>
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
						<td>
							<router-link to="/cdb/browse" class="btn btn-success" :disabled="busy">
								<i class="fa-solid fa-box-open"></i>
								Add from ContentDB
							</router-link>
						</td>
					</tr>
					<tr v-for="mod in get_mods()">
						<td>
							<span class="badge bg-secondary">{{mod.mod_type}}</span>
						</td>
						<td>
							<cdb-package-link v-if="mod.source_type == 'cdb'" :author="mod.author" :name="mod.name"/>
							<span v-else>{{mod.name}}</span>
						</td>
						<td>
							<span class="badge bg-success">
								<i class="fa-solid fa-box-open" v-if="mod.source_type == 'cdb'"></i>
								<i class="fa-brands fa-git-alt" v-if="mod.source_type == 'git'"></i>
								{{mod.source_type}}
							</span>
						</td>
						<td>
							<a :href="mod.url" v-if="mod.source_type == 'git'">{{mod.url}}</a>
						</td>
						<td>
							<span class="badge bg-secondary" v-if="mod.version">{{mod.version}}</span>
						</td>
						<td>
							<a class="btn btn-success" v-if="mod.auto_update">
								<i class="fa fa-check"></i>
								Enabled
							</a>
							<a class="btn btn-secondary" v-if="!mod.auto_update">
								<i class="fa fa-times"></i>
								Disabled
							</a>
						</td>
						<td>
							<div class="btn-group">
								<a class="btn btn-primary">
									<i class="fa fa-edit"></i>
									Edit
								</a>
								<a class="btn btn-danger" v-on:click="remove(mod.id)">
									<i class="fa fa-trash"></i>
									Remove
								</a>
							</div>
						</td>
					</tr>
				</tbody>
			</table>
		</default-layout>
	`
};
