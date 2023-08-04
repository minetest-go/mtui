import app_info from "../../store/app_info.js";

export default {
	data: () => app_info,
	template: /*html*/`
	<div>
		<div class="text-center">
			<h3>
				Minetest Web UI
				<small class="text-muted" v-if="servername">{{servername}}</small>
			</h3>
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
