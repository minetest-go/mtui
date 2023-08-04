import App from './app.js';
import routes from './routes.js';
import messages from './messages.js';
import { check_login } from './service/login.js';
import { check_features } from './service/features.js';
import router_guards from './util/router_guards.js';
import { fetch_info } from './service/app_info.js';
import { connect } from './ws.js';

function start(){
	// fetch app info
	fetch_info();

	// create router instance
	const router = VueRouter.createRouter({
		history: VueRouter.createWebHashHistory(),
		routes: routes
	});

	// set up router guards
	router_guards(router);

	const i18n = VueI18n.createI18n({
		fallbackLocale: 'en',
		messages: messages
	});

	// set up websocket events
	connect();

	// start vue
	const app = Vue.createApp(App);
	app.component('vue-datepicker', VueDatePicker);
	app.use(router);
	app.use(i18n);
	app.mount("#app");
}

check_features()
.then(() => check_login())
.then(() => start());
