
export const protected_fetch = (url, opts) => fetch(url, opts)
    .then(r => r.status == 200 ? r.json() : r.text())
    .then(o => typeof(o) == "string" ? Promise.reject(o) : o);
