<script>
import { is_logged_in } from "../service/login.js";

import StatsDisplay from "./StatsDisplay.vue";
import ThemeSwitcher from "./ThemeSwitcher.vue";

export default {
	emits: ["toggle-sidebar"],
	components: {
		"stats-display": StatsDisplay,
		"theme-switcher": ThemeSwitcher
	},
	computed: {
		is_logged_in
	}
};
</script>

<template>
  <header class="app-topbar navbar bg-body border-bottom px-3">
  	<button
  		type="button"
  		class="btn btn-outline-secondary d-lg-none me-2"
  		v-if="is_logged_in"
  		v-on:click="$emit('toggle-sidebar')"
  		aria-label="Open menu">
  		<i class="fa-solid fa-bars"></i>
  	</button>
  	<router-link to="/" class="navbar-brand" :class="{'d-lg-none': is_logged_in}">
  		Luanti Web UI
  	</router-link>
  	<div class="d-flex align-items-center ms-auto gap-2">
  		<stats-display class="navbar-text d-none d-sm-inline text-nowrap"/>
  		<router-link to="/login" class="btn btn-primary" v-if="!is_logged_in && $route.path != '/login'">
  			<i class="fa-solid fa-right-to-bracket"></i>
  			Login
  		</router-link>
  		<theme-switcher/>
  	</div>
  </header>
</template>
