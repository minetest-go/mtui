import { has_priv } from "../store/login.js";
import { logout } from '../service/login.js';

export default {
	methods: {
		has_priv: has_priv,
		logout: function(){
			logout().then(() => this.$router.push("/login"));
		}
	},
	template: /*html*/`
		<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
			<div class="container-fluid">
				<router-link to="/" class="navbar-brand">MT UI</router-link>
				<ul class="navbar-nav me-auto mb-2 mb-lg-0">
					<li class="nav-item" v-if="has_priv('interact')">
						<router-link to="/" class="nav-link">
							<i class="fa fa-home"></i> Home
						</router-link>
					</li>
					<li class="nav-item" v-if="has_priv('interact')">
						<router-link to="/profile" class="nav-link">
							<i class="fa fa-user"></i> Profile
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
				<div class="d-flex">
					<button class="btn btn-secondary" v-on:click="logout">
						<i class="fa-solid fa-right-from-bracket"></i>
						Logout
					</button>
				<div>
			</div>
		</nav>
	`
};
