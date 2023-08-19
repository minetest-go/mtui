import error_toast_store from "../store/error_toast.js";

export const protected_fetch = (url, opts) => fetch(url, opts)
    .then(r => {
        if (r.status == 200) {
            return r.json();
         } else {
            return r.text().then(msg => {
                error_toast_store.status = r.status;
                return Promise.reject(msg);
            });
         }
    })
    .catch(e => {
        error_toast_store.title = `HTTP fetch error`;
        error_toast_store.url = url;
        error_toast_store.message = e;
    });