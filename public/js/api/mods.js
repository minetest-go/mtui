import { protected_fetch } from "./util.js";

export const list_mods = () => protected_fetch("api/mods");

export const validate = () => protected_fetch("api/mods/validate");

export const create_mod = mod => protected_fetch("api/mods", {
    method: "POST",
    body: JSON.stringify(mod)
});

export const remove_mod = id => fetch(`api/mods/${id}`, {method: "DELETE"});
