import App from './app.js';
import routes from './routes.js';
import { check_login, has_priv } from './service/login.js';
import { check_features, has_feature } from './service/features.js';
import router_guards from './util/router_guards.js';
import { fetch_info } from './service/app_info.js';
import { update as update_modlist } from './service/mods.js';
import events, { EVENT_STARTUP } from './events.js';
import { start_polling, get_stats } from './service/stats.js';

function start(){
	// fetch app info
	fetch_info();

	// start stats polling
	start_polling();

	if (has_feature("modmanagement") && has_priv("server")) {
		update_modlist();
	}

	// create router instance
	const router = VueRouter.createRouter({
		history: VueRouter.createWebHashHistory(),
		routes: routes
	});

	// set up router guards
	router_guards(router);

	// trigger startup event
	events.emit(EVENT_STARTUP);

	// start vue
	const app = Vue.createApp(App);
	app.component('vue-datepicker', VueDatePicker);
	app.use(router);
	app.mount("#app");
}

get_stats()
.then(stats => {
	if (stats.maintenance) {
		// skip feature checking
		return;
	}
	return check_features();
})
.then(() => check_login())
.then(() => start());
