import { protected_fetch } from "./util.js";

export const get_mesecon_controls = () => protected_fetch(`api/mesecons`);

export const set_mesecon = m => fetch(`api/mesecons`, {
    method: "POST",
    body: JSON.stringify(m)
});

export const delete_mesecon = m => fetch(`api/mesecons/${m.poskey}`, {
    method: "DELETE"
});