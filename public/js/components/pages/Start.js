
const store = Vue.reactive({ version: "" });
import { get_version } from "../../api/version.js";

export default {
	data: () => store,
	mounted: function() {
		get_version().then(v => this.version = v);
	},
	template: /*html*/`
	<div>
		<div class="text-center">
			<h4>Start page</h4>
			<span v-if="version">
				Version: <span class="badge bg-primary">{{ version }}</span>
			</span>
			<hr/>
			<router-link to="/shell" class="btn btn-primary">
				<i class="fa-solid fa-terminal"></i> Shell
			</router-link>
			&nbsp;
			<router-link to="/profile" class="btn btn-primary">
				<i class="fa fa-user"></i> Profile
			</router-link>
			&nbsp;
			<a class="btn btn-secondary" href="https://github.com/minetest-go/mtui" target="new">
				<i class="fa-brands fa-github"></i> Source
			</a>
		</div>
	</div>
	`
};
