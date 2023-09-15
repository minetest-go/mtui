import events, { EVENT_LOGGED_IN } from "../events.js";
import { has_feature } from "./features.js";
import { has_priv } from "./login.js";
import { get_versions, get_status, start as api_start, stop as api_stop, remove as api_remove, create as api_create } from "../api/service_engine.js";

export const store = Vue.reactive({
	versions: null,
	busy: false,
	status: null,
	version: ""
});

export const update_state = () => {
	store.busy = true;
	get_status()
	.then(s => store.status = s)
	.then(() => store.version = store.status.version)
	.finally(() => store.busy = false);
};

events.on(EVENT_LOGGED_IN, function() {
	if (has_feature("docker") && has_priv("server")) {
		update_state();
		get_versions()
		.then(v => store.versions = v);
	}
});

export const start = () => {
	store.busy = true;
	api_start()
	.then(() => update_state());
};

export const stop = () => {
	store.busy = true;
	api_stop()
	.then(() => update_state());
};

export const create = () => {
	store.busy = true;
	api_create({version: store.version})
	.then(() => update_state());
};

export const remove = () => {
	store.busy = true;
	api_remove()
	.then(() => update_state());
};

export const is_running = () => store.status ? store.status.running : false;
