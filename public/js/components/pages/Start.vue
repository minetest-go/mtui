<script>
import SkinPreview from "../SkinPreview.vue";

import { get_servername, get_version } from "../../service/app_info.js";
import { has_priv, is_logged_in, get_claims } from "../../service/login.js";
import { has_feature } from "../../service/features.js";
import { store as stats, get_stats } from "../../service/stats.js";
import { get_unread_count } from "../../service/mail.js";
import { get_all as get_mods, get_git_mod, update as update_mods } from "../../service/mods.js";
import { count as count_logs } from "../../api/log.js";
import format_duration from "../../util/format_duration.js";

const DAY = 24 * 3600;

export default {
	components: {
		"skin-preview": SkinPreview
	},
	data: function() {
		return {
			stats,
			stats_loaded: false,
			mods_loaded: false,
			engine_errors: 0,
			mod_errors: 0
		};
	},
	mounted: function() {
		get_stats()
			.then(() => this.stats_loaded = true)
			.catch(() => {});

		if (this.can_manage_mods) {
			update_mods()
				.then(() => this.mods_loaded = true)
				.catch(() => {});
		}

		if (has_priv("ban")) {
			const since = { from_timestamp: Date.now() - DAY * 1000, to_timestamp: Date.now() + 3600 * 1000 };
			count_logs(Object.assign({ category: "minetest", event: "logfile-error" }, since))
				.then(c => this.engine_errors = c)
				.catch(() => {});
			count_logs(Object.assign({ category: "service", event: "engine-mod-errors" }, since))
				.then(c => this.mod_errors = c)
				.catch(() => {});
		}
	},
	methods: {
		has_priv,
		has_feature
	},
	computed: {
		is_logged_in,
		servername: get_servername,
		version: get_version,
		username: function() {
			const c = get_claims();
			return c && c.username;
		},
		online: function() {
			return this.stats.uptime > 0;
		},
		status: function() {
			if (!this.stats_loaded) {
				return { text: "Checking…", detail: "", color: "secondary", icon: "fa-solid fa-spinner fa-spin" };
			}
			if (!this.online) {
				return {
					text: "Offline",
					detail: "The game server is not responding. It may be starting, restarting or stopped.",
					color: "danger",
					icon: "fa-solid fa-circle-xmark"
				};
			}
			const ms = Math.floor(this.stats.max_lag * 1000);
			if (this.stats.max_lag < 0.2) {
				return { text: "Running smoothly", detail: `${ms} ms lag`, color: "success", icon: "fa-solid fa-circle-check" };
			}
			if (this.stats.max_lag < 0.5) {
				return { text: "Some lag", detail: `${ms} ms lag`, color: "warning", icon: "fa-solid fa-circle-exclamation" };
			}
			return { text: "Heavy lag", detail: `${ms} ms lag`, color: "danger", icon: "fa-solid fa-triangle-exclamation" };
		},
		uptime: function() {
			return format_duration(this.stats.uptime) || "just started";
		},
		clock: function() {
			const hours = this.stats.time_of_day * 24;
			const h = Math.floor(hours);
			const m = Math.floor((hours - h) * 60);
			return {
				text: `${h}:${m < 10 ? "0" + m : m}`,
				day: h >= 6 && h < 18
			};
		},
		players: function() {
			return this.stats.players || [];
		},
		can_manage_mods: function() {
			return has_priv("server") && has_feature("modmanagement");
		},
		show_attention: function() {
			return has_priv("server") || has_priv("ban");
		},
		// things an admin or moderator should look at, most important first
		attention: function() {
			const items = [];

			if (has_priv("server") && this.stats_loaded && !this.online) {
				items.push({
					color: "danger", icon: "fa-solid fa-power-off",
					text: "The game server is offline or not responding",
					link: "/log?category=service&event=engine&since=3600", action: "View engine log"
				});
			}

			if (this.can_manage_mods && this.mods_loaded) {
				const mods = get_mods();
				const failed = mods.filter(m => m.mod_status == "error");
				if (failed.length > 0) {
					items.push({
						color: "danger", icon: "fa-solid fa-cubes",
						text: failed.length == 1 ?
							`The mod "${failed[0].name}" has an error` :
							`${failed.length} mods have errors`,
						link: "/mods", action: "Open mods"
					});
				}
				if (!get_git_mod("mtui")) {
					items.push({
						color: "warning", icon: "fa-solid fa-puzzle-piece",
						text: "The mtui mod is not installed, some features will not work",
						link: "/mods", action: "Install"
					});
				}
				const updates = mods.filter(m => m.latest_version && m.latest_version != m.version);
				if (updates.length > 0) {
					items.push({
						color: "info", icon: "fa-solid fa-download",
						text: updates.length == 1 ?
							`An update is available for "${updates[0].name}"` :
							`Updates are available for ${updates.length} mods`,
						link: "/mods", action: "Review updates"
					});
				}
			}

			if (this.mod_errors > 0) {
				items.push({
					color: "warning", icon: "fa-solid fa-bug",
					text: `${this.mod_errors} mod error report(s) in the last 24 hours`,
					link: `/log?category=service&event=engine-mod-errors&since=${DAY}`, action: "View"
				});
			}

			if (this.engine_errors > 0) {
				items.push({
					color: "warning", icon: "fa-solid fa-file-circle-exclamation",
					text: `${this.engine_errors} engine error(s) logged in the last 24 hours`,
					link: `/log?category=minetest&event=logfile-error&since=${DAY}`, action: "View"
				});
			}

			return items;
		},
		quick_actions: function() {
			return [{
				name: "Chat", icon: "fa-solid fa-comment", link: "/chat",
				visible: has_priv("shout") && has_feature("chat")
			}, {
				name: "Mail", icon: "fa-solid fa-envelope", link: "/mail",
				visible: has_feature("mail"), badge: get_unread_count()
			}, {
				name: "Skin", icon: "fa-solid fa-user-astronaut", link: "/skin",
				visible: has_feature("skinsdb")
			}, {
				name: "Commands", icon: "fa-solid fa-terminal", link: "/shell",
				visible: has_priv("interact") && has_feature("shell")
			}, {
				name: "My profile", icon: "fa-solid fa-user", link: "/profile",
				visible: has_priv("interact")
			}].filter(a => a.visible);
		}
	}
};
</script>

<template>
  <div class="dashboard mx-auto">
  	<div class="mb-4">
  		<h3 class="mb-0" v-if="is_logged_in">Welcome back, {{username}}</h3>
  		<h3 class="mb-0" v-else>{{servername || "Luanti server"}}</h3>
  		<div class="text-body-secondary" v-if="is_logged_in && servername">{{servername}}</div>
  	</div>

  	<!-- status tiles -->
  	<div class="row g-3 mb-3">
  		<div class="col-md-4">
  			<div class="card h-100" :class="'border-' + status.color">
  				<div class="card-body">
  					<div class="small text-body-secondary">Server status</div>
  					<div class="fs-5 fw-semibold" :class="'text-' + status.color">
  						<i :class="status.icon"></i>
  						{{status.text}}
  					</div>
  					<div class="small text-body-secondary">{{status.detail}}</div>
  				</div>
  			</div>
  		</div>
  		<div class="col-md-4">
  			<div class="card h-100">
  				<div class="card-body">
  					<div class="small text-body-secondary">Players online</div>
  					<div class="fs-5 fw-semibold">
  						<i class="fa-solid fa-users"></i>
  						{{online ? stats.player_count : "–"}}
  					</div>
  					<router-link to="/online-players" class="small" v-if="online && stats.player_count > 0">Show all</router-link>
  				</div>
  			</div>
  		</div>
  		<div class="col-md-4">
  			<div class="card h-100">
  				<div class="card-body">
  					<div class="small text-body-secondary">In-game time</div>
  					<div class="fs-5 fw-semibold" v-if="online">
  						<i class="fa-solid fa-sun text-warning" v-if="clock.day"></i>
  						<i class="fa-solid fa-moon text-info" v-else></i>
  						{{clock.text}}
  					</div>
  					<div class="fs-5 fw-semibold" v-else>–</div>
  					<div class="small text-body-secondary" v-if="online">Up for {{uptime}}</div>
  				</div>
  			</div>
  		</div>
  	</div>

  	<div class="row g-3">
  		<div class="col-lg-8">
  			<!-- needs attention (admins and moderators) -->
  			<div class="card mb-3" v-if="is_logged_in && show_attention">
  				<div class="card-header fw-semibold">Needs attention</div>
  				<ul class="list-group list-group-flush" v-if="attention.length > 0">
  					<li class="list-group-item d-flex align-items-center gap-2" v-for="item in attention" :key="item.text">
  						<i :class="[item.icon, 'text-' + item.color]" class="fa-fw"></i>
  						<span class="me-auto">{{item.text}}</span>
  						<router-link :to="item.link" class="btn btn-sm btn-outline-secondary text-nowrap">{{item.action}}</router-link>
  					</li>
  				</ul>
  				<div class="card-body text-success" v-else>
  					<i class="fa-solid fa-circle-check"></i>
  					All good, nothing needs your attention.
  				</div>
  			</div>

  			<!-- quick actions -->
  			<div class="card" v-if="is_logged_in && quick_actions.length > 0">
  				<div class="card-header fw-semibold">Quick actions</div>
  				<div class="card-body">
  					<div class="row g-2">
  						<div class="col-6 col-md-4" v-for="action in quick_actions" :key="action.link">
  							<router-link :to="action.link" class="btn btn-outline-primary w-100 py-3 position-relative">
  								<div class="fs-4 mb-1"><i :class="action.icon"></i></div>
  								{{action.name}}
  								<span class="position-absolute top-0 end-0 m-2 badge rounded-pill text-bg-danger" v-if="action.badge">
  									{{action.badge}}
  								</span>
  							</router-link>
  						</div>
  					</div>
  				</div>
  			</div>

  			<!-- logged out -->
  			<div class="card" v-if="!is_logged_in">
  				<div class="card-body">
  					<p>Log in to chat, read your mail, change your skin and more.</p>
  					<router-link to="/login" class="btn btn-primary me-2">
  						<i class="fa-solid fa-right-to-bracket"></i> Log in
  					</router-link>
  					<router-link to="/signup" class="btn btn-outline-primary" v-if="has_feature('signup')">
  						<i class="fa-solid fa-user-plus"></i> Create account
  					</router-link>
  				</div>
  			</div>
  		</div>

  		<!-- online players -->
  		<div class="col-lg-4">
  			<div class="card">
  				<div class="card-header fw-semibold">Online now</div>
  				<ul class="list-group list-group-flush" v-if="players.length > 0">
  					<li class="list-group-item" v-for="player in players" :key="player.name">
  						<router-link :to="'/profile/' + player.name" class="d-flex align-items-center gap-2 text-decoration-none" v-if="is_logged_in">
  							<skin-preview :playername="player.name"/>
  							{{player.name}}
  						</router-link>
  						<span class="d-flex align-items-center gap-2" v-else>
  							<skin-preview :playername="player.name"/>
  							{{player.name}}
  						</span>
  					</li>
  				</ul>
  				<div class="card-body text-body-secondary" v-else>
  					Nobody is online right now.
  				</div>
  			</div>
  		</div>
  	</div>

  	<div class="small text-body-secondary text-center mt-4">
  		Luanti Web UI
  		<span v-if="version">{{version}}</span>
  		·
  		<a href="https://github.com/luanti-golang/mtui" target="_blank" rel="noopener" class="link-secondary">
  			<i class="fa-brands fa-github"></i> Source
  		</a>
  	</div>
  </div>
</template>

<style scoped>
.dashboard {
	max-width: 72rem;
}
</style>
