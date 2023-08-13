import { protected_fetch } from "./util.js";

export const list_mods = () => fetch("api/mods").then(r => r.json());

export const create_mod = mod => protected_fetch("api/mods", {
    method: "POST",
    body: JSON.stringify(mod)
});

export const remove_mod = id => fetch(`api/mods/${id}`, {method: "DELETE"});
