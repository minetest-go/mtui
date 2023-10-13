import events, { EVENT_LOGGED_IN } from "../events.js";
import { has_feature } from "./features.js";
import { has_priv } from "./login.js";
import { get_versions, get_status, start as api_start, stop as api_stop, remove as api_remove, create as api_create } from "../api/service.js";

const service_map = {};

class Service {
    constructor(servicename) {
        this.servicename = servicename;
        this.store = Vue.reactive({
            versions: null,
            busy: false,
            status: null,
            version: ""
        });
        service_map[servicename] = this;

        events.on(EVENT_LOGGED_IN, () => {
            if (has_feature("docker") && has_priv("server")) {
                this.update_state();
                get_versions(servicename)
                .then(v => this.store.versions = v);
            }
        });
    }

    update_state() {
        this.store.busy = true;
        get_status(this.servicename)
        .then(s => this.store.status = s)
        .then(() => this.store.version = this.store.status.version)
        .finally(() => this.store.busy = false);
    }

    start() {
        this.store.busy = true;
        return api_start(this.servicename).then(() => this.update_state());
    }
    
    stop() {
        this.store.busy = true;
        return api_stop(this.servicename).then(() => this.update_state());
    }
    
    restart() {
        this.store.busy = true;
        return api_stop(this.servicename)
        .then(() => api_start(this.servicename))
        .then(() => this.update_state());
    }
    
    create() {
        this.store.busy = true;
        return api_create(this.servicename, {version: this.store.version}).then(() => this.update_state());
    }
    
    remove() {
        this.store.busy = true;
        return api_remove(this.servicename).then(() => this.update_state());
    }
    
    is_running() {
        return this.store.status ? this.store.status.running : false;
    }

}

export const engine = new Service("engine");
export const matterbridge = new Service("matterbridge");
export const mapserver = new Service("mapserver");
export const mtweb = new Service("mtweb");

export const get_service_by_name = name => service_map[name];