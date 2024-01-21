import App from './app.js';
import routes from './routes.js';
import { check_login } from './service/login.js';
import { check_features } from './service/features.js';
import router_guards from './util/router_guards.js';
import { fetch_info } from './service/app_info.js';
import events, { EVENT_STARTUP } from './events.js';
import { start_polling, stop_polling, get_stats } from './service/stats.js';

function start(){
	// start stats polling
	start_polling();

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
	app.provide("unmount", () => {
		app.unmount();
		stop_polling();
	});
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
.then(() => fetch_info())
.then(() => check_login())
.then(() => start());
