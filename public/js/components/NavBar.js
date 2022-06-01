import { has_priv } from "../store/login.js";

export default {
	methods: {
		has_priv: has_priv
	},
	template: /*html*/`
		<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
			<div class="container-fluid">
				<router-link to="/" class="navbar-brand">MT Admin</router-link>
				<ul class="navbar-nav me-auto mb-2 mb-lg-0">
					<li class="nav-item" v-if="has_priv('interact')">
						<router-link to="/" class="nav-link">
							<i class="fa fa-home"></i> Home
						</router-link>
					</li>
					<li class="nav-item" v-if="has_priv('server')">
						<router-link to="/mods" class="nav-link">
							<i class="fa fa-bell" style="color: yellow;"></i>
							<i class="fa fa-puzzle-piece"></i>
							Mods
						</router-link>
					</li>
					<li class="nav-item" v-if="has_priv('server')">
						<router-link to="/engine" class="nav-link">
							<i class="fa fa-gears"></i>
							<i class="fa fa-bell" style="color: yellow;"></i>
							<i class="fa fa-play" style="color: green;"></i>
							Minetest engine
						</router-link>
					</li>
				</ul>
			</div>
		</nav>
	`
};
