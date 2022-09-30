import { list, scan } from '../../api/mods.js';

const store = Vue.reactive({
	list: []
});

export default {
	data: () => store,
	mounted: function() {
		if (this.list.length == 0) {
			this.update();
		}
	},
	methods: {
		scan: function() {
			scan()
			.then(() => this.update());
		},
		update: function() {
			list()
			.then(l => this.list = l);
		}
	},
	template: /*html*/`
		<div>
			<h3>Mod selection</h3>
			<table class="table table-condensed">
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
								</a>
								<a class="btn btn-danger">
									<i class="fa fa-times"></i>
								</a>
							</div>
						</td>
					</tr>
				</tbody>
				<tbody>
					<tr>
						<td>
							<select class="form-control">
								<option>Mod</option>
								<option>Game</option>
								<option>Textures</option>
							</select>
						</td>
						<td></td>
						<td>
							<select class="form-control">
								<option>Git</option>
								<option>ContentDB</option>
							</select>
						</td>
						<td>
							<input class="form-control" type="text" placeholder="Source url"/>
						</td>
						<td>
							<input class="form-control" type="text" placeholder="Version"/>
						</td>
						<td>
						</td>
						<td>
							<a class="btn btn-success">
								<i class="fa fa-save"></i>
							</a>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	`
};
