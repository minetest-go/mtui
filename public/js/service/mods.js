import {
    list_mods,
    create_mod,
    remove_mod,
    update_mod as api_update_mod,
    update_mod_version as api_update_mod_version,
    check_updates as api_check_updates,
    create_mtui_mod,
    create_beerchat_mod,
    create_mapserver_mod
} from '../api/mods.js';

import events, { EVENT_LOGGED_IN } from '../events.js';
import { has_priv } from './login.js';
import { has_feature } from './features.js';

const store = Vue.reactive({
    list: [],
    busy: false,
    has_mtui_mod: false,
    error_msg: ""
});

export async function update() {
    store.busy = true;
    store.list = await list_mods();
    store.has_mtui_mod = store.list.find(m => m.name == "mtui");
    store.busy = false;
}

export async function update_mod(m) {
    store.busy = true;
    await api_update_mod(m);
    await update();
}

export async function update_mod_version(m, v) {
    store.busy = true;
    await api_update_mod_version(m, v);
    await update();
};

export const is_busy = () => store.busy;

export async function add(m) {
    await create_mod(m);
    await update();
}

export async function add_mtui() {
    store.busy = true;
    await create_mtui_mod();
    await update();
}

export async function add_beerchat() {
    store.busy = true;
    await create_beerchat_mod();
    await update();
}

export async function add_mapserver() {
    store.busy = true;
    await create_mapserver_mod();
    await update();
}

export const remove = async id => await remove_mod(id).then(update);

export const get_all = () => store.list;

export const get_mods_by_type = type => store.list.filter(m => m.mod_type == type);

export const get_cdb_mod = (author, name) => store.list.find(m => m.name == name && m.author == author);

export const get_git_mod = name => store.list.find(m => m.name == name);

export const get_game = () => store.list.find(m => m.mod_type == "game");

export async function check_updates() {
    store.busy = true;
    await api_check_updates();
    await update();
}

events.on(EVENT_LOGGED_IN, function() {
    if (!has_priv("server") || !has_feature("modmanagement")){
        return;
    }

    update();
});
