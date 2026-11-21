<script>
import { remove, update_mod, update_mod_version } from '../../../service/mods.js';
import { confirm } from '../../../store/confirm.js';
import { success } from '../../../store/toast.js';

export default {
	props: ["mod"],
	components: {
	},
	methods: {
		remove: async function(mod) {
			const ok = await confirm({
				title: `Remove ${mod.name}?`,
				message: `The ${mod.mod_type} "${mod.name}" and its files will be deleted from the server.`,
				confirm_label: "Remove",
				confirm_icon: "fa fa-trash"
			});
			if (ok) {
				await remove(mod.id);
				success(`Removed ${mod.name}`);
			}
		},
		update_mod_version: async function(mod, version) {
			await update_mod_version(mod, version);
			success(`Updating ${mod.name} to ${version}`);
		},
		update_mod: async function(mod) {
			await update_mod(mod);
			success(`Auto-update ${mod.auto_update ? "enabled" : "disabled"} for ${mod.name}`);
		}
	}
};
</script>

<template>
  <td>
  	<span class="badge bg-secondary">{{mod.mod_type}}</span>
  </td>
  <td>
  	<span>{{mod.name}}</span>
  </td>
  <td>
  	<a :href="mod.url" v-if="mod.source_type == 'git'">{{mod.url}}</a>
  	<span class="badge bg-secondary" v-if="mod.source_type == 'manual'">manual</span>
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
  		<button class="btn btn-danger" v-on:click="remove(mod)" v-if="mod.mod_status != 'processing'">
  			<i class="fa fa-trash"></i>
  			Remove
  		</button>
  	</div>
  	<i class="fa fa-spin fa-spinner" v-if="mod.mod_status == 'processing'"></i>
  </td>
</template>
