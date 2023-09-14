import { list_mods, create_mod, remove_mod, update_mod as api_update_mod, check_updates as api_check_updates } from '../api/mods.js';

const store = Vue.reactive({
    list: [],
    busy: false,
    has_mtui_mod: false,
    error_msg: ""
});

export const update = () => {
    store.busy = true;
    list_mods()
    .then(l => store.list = l)
    .then(() => store.has_mtui_mod = store.list.find(m => m.name == "mtui"))
    .finally(() => store.busy = false);
};

export const update_mod = (m, v) => {
    store.busy = true;
    api_update_mod(m, v)
    .then(() => update());
};

export const is_busy = () => store.busy;

export const add = m => create_mod(m).then(update);

export const remove = id => remove_mod(id).then(update);

export const get_all = () => store.list;

export const get_cdb_mod = (author, name) => store.list.find(m => m.name == name && m.author == author);

export const get_git_mod = name => store.list.find(m => m.name == name);

export const get_game = () => store.list.find(m => m.mod_type == "game");

export const check_updates = () => {
    store.busy = true;
    api_check_updates()
    .then(() => update());
};

