import App from './app.js';
import routes from './routes.js';
import messages from './messages.js';
import { check_login, has_priv } from './service/login.js';
import { check_features, has_feature } from './service/features.js';
import router_guards from './util/router_guards.js';
import { fetch_info } from './service/app_info.js';
import { update as update_modlist } from './service/mods.js';
import { connect } from './ws.js';

function start(){
	// fetch app info
	fetch_info();

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
