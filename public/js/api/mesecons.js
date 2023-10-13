import { protected_fetch } from "./util.js";

export const get_mesecon_controls = () => protected_fetch(`api/mesecons`);

export const set_mesecon = m => protected_fetch(`api/mesecons`, {
    method: "POST",
    body: JSON.stringify(m)
});

export const delete_mesecon = poskey => fetch(`api/mesecons/${poskey}`, {
    method: "DELETE"
});