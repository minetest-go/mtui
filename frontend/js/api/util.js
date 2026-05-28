import error_toast_store from "../store/error_toast.js";

export async function protected_fetch(url, opts) {
    try {
        const r = await fetch(url, opts);
        if (r.status == 200) {
            return await r.json();
         } else {
            const msg = await r.text();
            error_toast_store.title = "HTTP status error";
            error_toast_store.url = url;
            error_toast_store.message = msg;
            throw new Error(msg);
         }
    } catch (e) {
        error_toast_store.title = `HTTP fetch error`;
        error_toast_store.url = url;
        error_toast_store.message = e;
        throw e;
    }
}