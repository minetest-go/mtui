import { protected_fetch } from "./util.js";

export const get_all = () => protected_fetch("api/mtconfig");

export const set = (key, value) => protected_fetch(`api/mtconfig/${key}`, {
    method: "POST",
    body: value
});

export const remove = key => fetch(`api/mtconfig/${key}`, {
    method: "DELETE",
});