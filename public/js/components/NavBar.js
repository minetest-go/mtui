import { has_priv } from "../service/login.js";
import { has_feature } from "../service/features.js";
import { logout } from '../service/login.js';
import login_store from '../store/login.js';
import stats_store from '../store/stats.js';
import mail_store from '../store/mail.js';
import StatsDisplay from './StatsDisplay.js';

export default {
	data: function() {
		return {
			login: login_store,
			stats: stats_store,
			mail: mail_store,
			admin_menu: false,
			mod_menu: false
		};
	},
	methods: {
		has_priv: has_priv,
		has_feature: has_feature,
		logout: function(){
			logout().then(() => this.$router.push("/login"));
		}
	},
	components: {
		"stats-display": StatsDisplay
	},
	template: /*html*/`
		<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
			<div class="container-fluid">
				<router-link to="/" class="navbar-brand">Minetest Web UI</router-link>
				<ul class="navbar-nav me-auto mb-2 mb-lg-0" v-if="login.loggedIn">
					<li class="nav-item" v-if="has_priv('interact')">
						<router-link to="/" class="nav-link">
							<i class="fa fa-home"></i> Home
						</router-link>
					</li>
					<li class="nav-item" v-if="has_priv('interact')">
						<router-link to="/playersearch" class="nav-link">
							<i class="fa fa-magnifying-glass"></i> Player search
						</router-link>
					</li>
					<li class="nav-item" v-if="has_priv('interact') && has_feature('shell')">
						<router-link to="/shell" class="nav-link">
							<i class="fa-solid fa-terminal"></i> Shell
						</router-link>
					</li>
					<li class="nav-item">
						<router-link to="/online-players" class="nav-link">
							<i class="fa fa-users"></i> Online players
							<span class="badge rounded-pill bg-info" v-if="stats.player_count">
								{{stats.player_count}}
							</span>
						</router-link>
					</li>
					<li class="nav-item" v-if="has_feature('mail')">
						<router-link to="/mail" class="nav-link">
							<i class="fa-solid fa-envelope"></i> Mail
							<span class="badge rounded-pill bg-info" v-if="mail.unread_count">
								{{mail.unread_count}}
							</span>
						</router-link>
					</li>
					<li class="nav-item" v-if="has_feature('skinsdb')">
						<router-link to="/skin" class="nav-link">
							<i class="fa-solid fa-user-astronaut"></i> Skin
						</router-link>
					</li>
					<li class="nav-item dropdown" v-if="has_priv('ban')" v-on:mouseleave="mod_menu = false">
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
					<li class="nav-item dropdown" v-if="has_priv('server')" v-on:mouseleave="admin_menu = false">
						<a class="nav-link dropdown-toggle" v-on:click="admin_menu = true" v-on:mouseover="admin_menu = true">
							<i class="fa-solid fa-screwdriver-wrench"></i>
							Administration
						</a>		
						<ul class="dropdown-menu" v-bind:class="{'show': admin_menu}">
							<li>
								<router-link to="/features" class="dropdown-item">
									<i class="fa fa-tags"></i> Features
								</router-link>
							</li>
							<li>
								<router-link to="/oauth-apps" class="dropdown-item">
									<i class="fa fa-passport"></i> OAuth apps
								</router-link>
							</li>
							<li v-if="has_feature('luashell')">
								<router-link to="/lua" class="dropdown-item">
									<i class="fa-solid fa-terminal"></i> Lua
								</router-link>
							</li>
							<li v-if="has_feature('modmanagement')">
								<router-link to="/mods" class="dropdown-item">
									<i class="fa fa-cubes"></i> Mods
								</router-link>
							</li>
							<li v-if="has_feature('mediaserver')">
								<router-link to="/mediaserver" class="dropdown-item">
									<i class="fa fa-photo-film"></i> Mediaserver
								</router-link>
							</li>
						</ul>
					</li>
				</ul>
				<div class="d-flex">
					<stats-display class="navbar-text" style="padding-right: 10px;"/>
					<div class="btn-group">
						<button class="btn btn-outline-secondary" v-if="login.claims">
							<router-link to="/profile">
								<i class="fas fa-user"></i>
								<span>
									Logged in as <b>{{login.claims.username}}</b>
								</span>
							</router-link>
						</button>
						<button class="btn btn-secondary" v-on:click="logout" v-if="login.loggedIn">
							<i class="fa-solid fa-right-from-bracket"></i>
							Logout
						</button>
					</div>
				<div>
			</div>
		</nav>
	`
};
