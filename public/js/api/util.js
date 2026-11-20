import { error } from "../store/toast.js";

export async function protected_fetch(url, opts) {
    let r;
    try {
        r = await fetch(url, opts);
    } catch (e) {
        error("Connection error", "The server could not be reached", `URL: ${url}\n${e}`);
        throw e;
    }

    if (r.status != 200) {
        const msg = await r.text();
        error("Request failed", msg || r.statusText, `URL: ${url}\nStatus: ${r.status}`);
        throw new Error(msg);
    }

    try {
        return await r.json();
    } catch (e) {
        error("Invalid response", "The server sent an unexpected response", `URL: ${url}\n${e}`);
        throw e;
    }
}
