import { list, create, remove } from '../../../api/mods.js';

const store = Vue.reactive({
	list: [],
	busy: false,
	add_name: "",
	add_mod_type: "mod",
	add_source_type: "git",
	add_url: "",
	add_version: ""
});

export default {
	data: () => store,
	mounted: function() {
		if (this.list.length == 0) {
			this.update();
		}
	},
	methods: {
		add: function() {
			this.busy = true;
			create({
				name: this.add_name,
				mod_type: this.add_mod_type,
				source_type: this.add_source_type,
				url: this.add_url,
				branch: "",
				version: this.add_version
			})
			.then(() => {
				this.busy = false;
				this.add_name = "";
				this.add_url = "";
				this.add_version = "";
				this.update();
			});
		},
		update: function() {
			list()
			.then(l => this.list = l);
		},
		remove: function(id) {
			remove(id)
			.then(() => this.update());
		}
	},
	template: /*html*/`
		<div>
			<h3>Mod management</h3>
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
					<tr v-for="mod in list">
						<td>
							<span class="badge bg-secondary">{{mod.mod_type}}</span>
						</td>
						<td>{{mod.name}}</td>
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
									<i class="fa fa-times"></i>
									Remove
								</a>
							</div>
						</td>
					</tr>
					<tr>
						<td>
							<select class="form-control" v-model="add_mod_type">
								<option value="mod">Mod</option>
								<option value="game">Game</option>
								<option value="txp">Textures</option>
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
							<button class="btn btn-success" v-on:click="add" :disabled="busy">
								<i class="fa-brands fa-git-alt"></i>
								Add from git
							</button>
						</td>
					</tr>
					<tr>
						<td>
						</td>
						<td>
						</td>
						<td>
							<span class="badge bg-success">
								<i class="fa-solid fa-box-open"></i>
								ContentDB
							</span>
						</td>
						<td>
						</td>
						<td>
						</td>
						<td>
						</td>
						<td>
							<router-link to="/cdb/browse" class="btn btn-success" :disabled="busy">
								<i class="fa-solid fa-box-open"></i>
								Add from ContentDB
							</router-link>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	`
};
