import { protected_fetch } from "./util.js";

export const list_mods = () => protected_fetch("api/mods");

export const validate = () => protected_fetch("api/mods/validate");

export const create_mod = mod => protected_fetch("api/mods", {
    method: "POST",
    body: JSON.stringify(mod)
});

export const create_mtui_mod = () => protected_fetch("api/mods/create_mtui", {
    method: "POST"
});

export const create_beerchat_mod = () => protected_fetch("api/mods/create_beerchat", {
    method: "POST"
});

export const create_mapserver_mod = () => protected_fetch("api/mods/create_mapserver", {
    method: "POST"
});

export const update_mod = mod => protected_fetch(`api/mods/${mod.id}`, {
    method: "POST",
    body: JSON.stringify(mod)
});

export const update_mod_version = (mod, version) => protected_fetch(`api/mods/${mod.id}/update/${version}`, {
    method: "POST"
});

export const check_updates = () => protected_fetch("api/mods/checkupdates", {
    method: "POST"
});

export const remove_mod = id => fetch(`api/mods/${id}`, {method: "DELETE"});
