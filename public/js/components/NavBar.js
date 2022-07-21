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
			mail: mail_store
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
						<router-link to="/profile" class="nav-link">
							<i class="fa fa-user"></i> Profile
						</router-link>
					</li>
					<li class="nav-item" v-if="has_priv('interact')">
						<router-link to="/shell" class="nav-link">
							<i class="fa-solid fa-terminal"></i> Shell
						</router-link>
					</li>
					<li class="nav-item">
						<router-link to="/online-players" class="nav-link">
							<i class="fa fa-users"></i> Online players
							<span class="badge rounded-pill bg-info">
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
				</ul>
				<div class="d-flex">
					<stats-display class="navbar-text" style="padding-right: 10px;"/>
					<button class="btn btn-secondary" v-on:click="logout" v-if="login.loggedIn">
						<i class="fa-solid fa-right-from-bracket"></i>
						Logout
					</button>
				<div>
			</div>
		</nav>
	`
};
