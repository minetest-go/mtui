<script>
import Sidebar from './components/Sidebar.vue';
import TopBar from './components/TopBar.vue';
import Toasts from './components/Toasts.vue';
import ConfirmDialog from './components/ConfirmDialog.vue';
import { is_logged_in } from './service/login.js';

export default {
	components: {
		"side-bar": Sidebar,
		"top-bar": TopBar,
		"toasts": Toasts,
		"confirm-dialog": ConfirmDialog
	},
	data: function() {
		return {
			sidebar_open: false
		};
	},
	computed: {
		is_logged_in
	},
	watch: {
		"$route": function() {
			// close the mobile menu after navigating
			this.sidebar_open = false;
		}
	}
};
</script>

<template>
  <div class="app-shell">
    <side-bar v-if="is_logged_in" :open="sidebar_open" v-on:close="sidebar_open = false"/>
    <div class="app-sidebar-backdrop d-lg-none" v-if="is_logged_in && sidebar_open" v-on:click="sidebar_open = false"></div>
    <div class="app-main">
      <top-bar v-on:toggle-sidebar="sidebar_open = !sidebar_open"/>
      <toasts/>
      <confirm-dialog/>
      <main class="container-fluid py-3">
        <router-view></router-view>
      </main>
    </div>
  </div>
</template>
