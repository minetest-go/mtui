<script>
import { remove, update_mod, update_mod_version } from '../../../service/mods.js';
import CDBPackageLink from '../../CDBPackageLink.vue';

export default {
	props: ["mod"],
	components: {
		"cdb-package-link": CDBPackageLink
	},
	methods: {
		remove,
		update_mod_version,
		update_mod
	}
};
</script>

<template>
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
</template>
