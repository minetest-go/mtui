import { has_priv, is_logged_in, get_claims, logout } from "../service/login.js";
import { has_feature } from "../service/features.js";
import { get_player_count, get_maintenance } from "../service/stats.js";
import { get_unread_count } from '../service/mail.js';
import StatsDisplay from './StatsDisplay.js';
import EngineStatus from "./pages/services/EngineStatus.js";

export default {
	data: function() {
		return {
			admin_menu: false,
			mod_menu: false,
			services_menu: false
		};
	},
	methods: {
		has_priv: has_priv,
		has_feature: has_feature,
		logout: function(){
			logout().then(() => this.$router.push("/login"));
		}
	},
	computed: {
		get_player_count: get_player_count,
		is_logged_in: is_logged_in,
		get_claims: get_claims,
		get_unread_count: get_unread_count,
		maintenance: get_maintenance
	},
	components: {
		"stats-display": StatsDisplay,
		"engine-status": EngineStatus
	},
	template: /*html*/`
		<nav class="navbar navbar-expand-lg navbar-dark" v-bind:class="{'bg-dark': !maintenance, 'bg-warning': maintenance}">
			<div class="container-fluid">
				<router-link to="/" class="navbar-brand">Minetest Web UI</router-link>
				<ul class="navbar-nav me-auto mb-2 mb-lg-0" v-if="is_logged_in">
					<li class="nav-item" v-if="has_priv('interact')">
						<router-link to="/" class="nav-link">
							<i class="fa fa-home"></i> Home
						</router-link>
					</li>
					<li class="nav-item" v-if="has_priv('interact') && !maintenance">
						<router-link to="/playersearch" class="nav-link">
							<i class="fa fa-magnifying-glass"></i> Player search
						</router-link>
					</li>
					<li class="nav-item" v-if="has_priv('interact') && has_feature('shell') && !maintenance">
						<router-link to="/shell" class="nav-link">
							<i class="fa-solid fa-terminal"></i> Shell
						</router-link>
					</li>
					<li class="nav-item">
						<router-link to="/online-players" class="nav-link" v-if="!maintenance">
							<i class="fa fa-users"></i> Online players
							<span class="badge rounded-pill bg-info">
								{{get_player_count}}
							</span>
						</router-link>
					</li>
					<li class="nav-item" v-if="has_feature('mail') && !maintenance">
						<router-link to="/mail" class="nav-link">
							<i class="fa-solid fa-envelope"></i> Mail
							<span class="badge rounded-pill bg-info" v-if="get_unread_count">
								{{get_unread_count}}
							</span>
						</router-link>
					</li>
					<li class="nav-item" v-if="has_feature('skinsdb') && !maintenance">
						<router-link to="/skin" class="nav-link">
							<i class="fa-solid fa-user-astronaut"></i> Skin
						</router-link>
					</li>
					<li class="nav-item dropdown" v-if="has_priv('ban') && !maintenance" v-on:mouseleave="mod_menu = false">
						<a class="nav-link dropdown-toggle" v-on:click="mod_menu = true" v-on:mouseover="mod_menu = true">
							<i class="fa-solid fa-hammer"></i>
							Moderation
						</a>		
						<ul class="dropdown-menu" v-bind:class="{'show': mod_menu}">
							<li v-if="has_feature('xban')">
								<router-link to="/xban" class="dropdown-item">
									<i class="fa fa-ban"></i> XBan
								</router-link>
							</li>
							<li>
								<router-link to="/log" class="dropdown-item">
									<i class="fa fa-magnifying-glass"></i> Logs
								</router-link>
							</li>	
						</ul>
					</li>
					<li class="nav-item dropdown" v-if="has_feature('docker') && !maintenance" v-on:mouseleave="services_menu = false">
						<a class="nav-link dropdown-toggle" v-on:click="services_menu = true" v-on:mouseover="services_menu = true">
							<i class="fa-solid fa-gears"></i>
							Services
						</a>		
						<ul class="dropdown-menu" v-bind:class="{'show': services_menu}">
							<li>
								<router-link to="/services/engine" class="dropdown-item">
									<i class="fa fa-gear"></i>
									Minetest engine
									<engine-status/>
								</router-link>
							</li>
						</ul>
					</li>
					<li class="nav-item dropdown" v-if="has_priv('server')" v-on:mouseleave="admin_menu = false">
						<a class="nav-link dropdown-toggle" v-on:click="admin_menu = true" v-on:mouseover="admin_menu = true">
							<i class="fa-solid fa-screwdriver-wrench"></i>
							Administration
						</a>		
						<ul class="dropdown-menu" v-bind:class="{'show': admin_menu}">
							<li v-if="!maintenance && has_feature('filebrowser')">
								<router-link to="/filebrowser/" class="dropdown-item">
									<i class="fa-solid fa-folder"></i> Filebrowser
								</router-link>
							</li>
							<li v-if="!maintenance">
								<router-link to="/ui/settings" class="dropdown-item">
									<i class="fa-solid fa-list-check"></i> UI Settings
								</router-link>
							</li>
							<li v-if="!maintenance">
								<router-link to="/features" class="dropdown-item">
									<i class="fa fa-tags"></i> Features
								</router-link>
							</li>
							<li v-if="!maintenance">
								<router-link to="/oauth-apps" class="dropdown-item">
									<i class="fa fa-passport"></i> OAuth apps
								</router-link>
							</li>
							<li v-if="has_feature('luashell') && !maintenance">
								<router-link to="/lua" class="dropdown-item">
									<i class="fa-solid fa-terminal"></i> Lua
								</router-link>
							</li>
							<li v-if="has_feature('minetest_config') && !maintenance">
								<router-link to="/minetest-config" class="dropdown-item">
									<i class="fa fa-cog"></i> Minetest config
								</router-link>
							</li>
							<li v-if="has_feature('modmanagement') && !maintenance">
								<router-link to="/mods" class="dropdown-item">
									<i class="fa fa-cubes"></i> Mods
								</router-link>
							</li>
							<li v-if="has_feature('mediaserver') && !maintenance">
								<router-link to="/mediaserver" class="dropdown-item">
									<i class="fa fa-photo-film"></i> Mediaserver
								</router-link>
							</li>
							<li>
								<router-link to="/maintenance" class="dropdown-item">
									<i class="fa fa-wrench"></i> Maintenance
								</router-link>
							</li>
						</ul>
					</li>
				</ul>
				<div class="d-flex">
					<stats-display class="navbar-text" style="padding-right: 10px;"/>
					<div class="btn-group" v-if="is_logged_in">
						<button class="btn btn-outline-secondary">
							<router-link to="/profile">
								<i class="fas fa-user"></i>
								<span>
									Logged in as <b>{{get_claims.username}}</b>
								</span>
							</router-link>
						</button>
						<button class="btn btn-secondary" v-on:click="logout" v-if="!maintenance">
							<i class="fa-solid fa-right-from-bracket"></i>
							Logout
						</button>
					</div>
				</div>
			</div>
		</nav>
	`
};
