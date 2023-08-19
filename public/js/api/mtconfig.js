import { protected_fetch } from "./util.js";


export const get_settingtypes = () => protected_fetch("api/mtconfig/settingtypes");

export const get_all = () => protected_fetch("api/mtconfig/settings");

export const set = (key, setting) => protected_fetch(`api/mtconfig/settings/${key}`, {
    method: "POST",
    body: JSON.stringify(setting)
});

export const remove = key => fetch(`api/mtconfig/settings/${key}`, {
    method: "DELETE",
});