import { protected_fetch } from "./util.js";

export const get_luacontroller = m => protected_fetch(`api/mesecons/luacontroller/get`, {
    method: "POST",
    body: JSON.stringify(m)
});

export const set_luacontroller = m => protected_fetch(`api/mesecons/luacontroller/set`, {
    method: "POST",
    body: JSON.stringify(m)
});
