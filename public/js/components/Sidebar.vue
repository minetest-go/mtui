<script>
import nav, { is_active } from "../util/nav.js";
import { get_claims, logout } from "../service/login.js";
import { get_servername } from "../service/app_info.js";
import { get_player_count } from "../service/stats.js";
import { get_unread_count } from "../service/mail.js";

import SkinPreview from "./SkinPreview.vue";

export default {
	props: ["open"],
	emits: ["close"],
	components: {
		"skin-preview": SkinPreview
	},
	computed: {
		sections: function() {
			return nav
				.map(section => Object.assign({}, section, {
					items: section.items.filter(item => !item.visible || item.visible())
				}))
				.filter(section => section.items.length > 0);
		},
		badges: function() {
			return {
				players: get_player_count(),
				mail: get_unread_count()
			};
		},
		claims: get_claims,
		servername: get_servername
	},
	methods: {
		is_active: function(item) {
			return is_active(item, this.$route.path);
		},
		logout: function() {
			logout().then(() => this.$router.push("/login"));
		}
	}
};
</script>

<template>
  <aside class="app-sidebar bg-body-tertiary border-end d-flex flex-column" :class="{open: open}">
  	<div class="d-flex align-items-center px-3 py-2 border-bottom app-sidebar-header">
  		<router-link to="/" class="text-decoration-none text-body-emphasis me-auto text-truncate">
  			<div class="fw-semibold">Luanti Web UI</div>
  			<div class="small text-body-secondary text-truncate" v-if="servername">{{servername}}</div>
  		</router-link>
  		<button type="button" class="btn-close d-lg-none" aria-label="Close menu" v-on:click="$emit('close')"></button>
  	</div>

  	<nav class="flex-grow-1 overflow-y-auto py-2" aria-label="Main navigation">
  		<div v-for="section in sections" :key="section.name">
  			<div class="app-sidebar-section">{{section.name}}</div>
  			<ul class="nav nav-pills flex-column px-2">
  				<li class="nav-item" v-for="item in section.items" :key="item.link">
  					<router-link
  						:to="item.link"
  						class="nav-link d-flex align-items-center"
  						:class="is_active(item) ? 'active' : 'link-body-emphasis'"
  						:aria-current="is_active(item) ? 'page' : null">
  						<i :class="item.icon" class="app-sidebar-icon"></i>
  						<span class="me-auto">{{item.name}}</span>
  						<span class="badge rounded-pill text-bg-info" v-if="item.badge && badges[item.badge]">
  							{{badges[item.badge]}}
  						</span>
  					</router-link>
  				</li>
  			</ul>
  		</div>
  	</nav>

  	<div class="border-top p-2 d-flex align-items-center" v-if="claims">
  		<router-link
  			to="/profile"
  			class="nav-link d-flex align-items-center flex-grow-1 text-truncate rounded px-2 py-1"
  			:class="$route.path == '/profile' ? 'active text-bg-primary' : 'link-body-emphasis'"
  			title="My profile">
  			<skin-preview :playername="claims.username" class="me-2"/>
  			<span class="text-truncate">{{claims.username}}</span>
  		</router-link>
  		<button type="button" class="btn btn-outline-secondary btn-sm ms-2" v-on:click="logout" title="Log out">
  			<i class="fa-solid fa-right-from-bracket"></i>
  			<span class="visually-hidden">Log out</span>
  		</button>
  	</div>
  </aside>
</template>
