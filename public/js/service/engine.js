import events, { EVENT_LOGGED_IN } from "../events.js";
import { has_feature } from "./features.js";
import { has_priv } from "./login.js";
import { get_versions, get_status, start as api_start, stop as api_stop, remove as api_remove, create as api_create } from "../api/service.js";

const servicename = "engine";

export const store = Vue.reactive({
	versions: null,
	busy: false,
	status: null,
	version: ""
});

export const update_state = () => {
	store.busy = true;
	get_status(servicename)
	.then(s => store.status = s)
	.then(() => store.version = store.status.version)
	.finally(() => store.busy = false);
};

events.on(EVENT_LOGGED_IN, function() {
	if (has_feature("docker") && has_priv("server")) {
		update_state();
		get_versions(servicename)
		.then(v => store.versions = v);
	}
});

export const start = () => {
	store.busy = true;
	return api_start(servicename).then(() => update_state());
};

export const stop = () => {
	store.busy = true;
	return api_stop(servicename).then(() => update_state());
};

export const restart = () => {
	store.busy = true;
	return api_stop(servicename)
	.then(() => api_start(servicename))
	.then(() => update_state());
};

export const create = () => {
	store.busy = true;
	return api_create(servicename, {version: store.version}).then(() => update_state());
};

export const remove = () => {
	store.busy = true;
	return api_remove(servicename).then(() => update_state());
};

export const is_running = () => store.status ? store.status.running : false;
