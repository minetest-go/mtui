import { list_mods, create_mod, remove_mod } from '../api/mods.js';

import store from '../store/mods.js'

export const update = () => {
    store.busy = true;
    list_mods()
    .then(l => store.list = l)
    .then(() => store.has_mtui_mod = store.list.find(m => m.name == "mtui"))
    .finally(() => store.busy = false);
};

export const add = m => {
    store.busy = true;
    return create_mod(m)
    .then(update)
	.catch(e => store.error_msg = e)
    .finally(() => store.busy = false);
};

export const remove = id => remove_mod(id).then(update);
