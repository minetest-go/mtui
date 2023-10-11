import { has_priv, is_logged_in, get_claims, logout } from "../service/login.js";
import { has_feature } from "../service/features.js";
import { get_player_count, get_maintenance } from "../service/stats.js";
import { get_unread_count } from '../service/mail.js';
import { engine, matterbridge, mapserver } from "../service/service.js";

import StatsDisplay from './StatsDisplay.js';
import ServiceStatus from "./pages/services/ServiceStatus.js";
import NavDropdown from "./NavDropdown.js";

export default {
	data: function() {
		return {
			engine,
			matterbridge,
			mapserver
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
		"service-status": ServiceStatus,
		"nav-dropdown": NavDropdown
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
					<li class="nav-item" v-if="has_priv('shout') && has_feature('chat') && !maintenance">
						<router-link to="/chat" class="nav-link">
							<i class="fa-solid fa-comment"></i> Chat
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
					<nav-dropdown v-if="has_priv('ban') && !maintenance" icon="hammer" name="Moderation">
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
					</nav-dropdown>
					<nav-dropdown v-if="has_feature('docker') && has_priv('server') && !maintenance" icon="gears" name="Services">
						<li>
							<router-link to="/services/engine" class="dropdown-item">
								<i class="fa fa-gear"></i>
								Minetest engine
								<service-status :status="engine.store.status"/>
							</router-link>
						</li>
						<li>
							<router-link to="/services/matterbridge" class="dropdown-item">
								<i class="fa fa-gear"></i>
								Matterbridge
								<service-status :status="matterbridge.store.status"/>
							</router-link>
						</li>
						<li>
							<router-link to="/services/mapserver" class="dropdown-item">
								<i class="fa fa-gear"></i>
								Mapserver
								<service-status :status="mapserver.store.status"/>
							</router-link>
						</li>
					</nav-dropdown>
					<nav-dropdown v-if="has_priv('server')" icon="screwdriver-wrench" name="Administration">
						<li>
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
					</nav-dropdown>
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
