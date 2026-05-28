import { protected_fetch } from "./util.js";

export const get_events = category => protected_fetch(`api/log/events/${category}`);

export const count = s => protected_fetch("api/log/count", {
    method: "POST",
    body: JSON.stringify(s)
});

export const search = s => protected_fetch("api/log/search", {
    method: "POST",
    body: JSON.stringify(s)
});