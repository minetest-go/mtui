import { get_appinfo } from "../api/app_info.js";

const store = Vue.reactive({
	version: "",
	servername: ""
});

export const fetch_info = () => {
    get_appinfo().then(i => Object.assign(store, i));
};

export const get_version = () => store.version;
export const get_servername = () => store.servername;