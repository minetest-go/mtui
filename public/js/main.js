import { createApp } from 'vue';
import { createRouter, createWebHashHistory } from 'vue-router';
import VueDatePicker from '@vuepic/vue-datepicker';

import 'bootstrap/dist/css/bootstrap.min.css';
import '@fortawesome/fontawesome-free/css/all.min.css';
import '@vuepic/vue-datepicker/dist/main.css';
import 'codemirror/lib/codemirror.css';
import './style.css';

import App from './App.vue';
import routes from './routes.js';
import { check_login } from './service/login.js';
import { check_features } from './service/features.js';
import router_guards from './util/router_guards.js';
import { fetch_info } from './service/app_info.js';
import events, { EVENT_STARTUP } from './events.js';
import { start_polling, stop_polling } from './service/stats.js';

async function start(){
	// check features
	await check_features();

	await fetch_info();
	await check_login();

	// start stats polling
	start_polling();

	// create router instance
	const router = createRouter({
		history: createWebHashHistory(),
		routes: routes
	});

	// set up router guards
	router_guards(router);

	// trigger startup event
	events.emit(EVENT_STARTUP);

	// start vue
	const app = createApp(App);
	app.component('vue-datepicker', VueDatePicker);
	app.use(router);
	app.provide("unmount", () => {
		app.unmount();
		stop_polling();
	});
	app.mount("#app");
}

start().catch(e => console.error(e));
