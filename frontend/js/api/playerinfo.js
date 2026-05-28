import { protected_fetch } from "./util.js";

export const get = name => protected_fetch(`api/player/info/${name}`);

export const get_skin = name => fetch(`api/player/skin/${name}`).then(r => r.blob());

export const search = q => protected_fetch(`api/player/search`, {
    method: "POST",
    body: JSON.stringify(q)
});

export const count = q => protected_fetch(`api/player/count`, {
    method: "POST",
    body: JSON.stringify(q)
});