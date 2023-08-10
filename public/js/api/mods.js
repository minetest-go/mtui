export const list = () => fetch("api/mods").then(r => r.json());

export const create = mod => fetch("api/mods", {
    method: "POST",
    body: JSON.stringify(mod)
})
.then(r => r.status == 200 ? r.json() : r.text())
.then(o => typeof(o) == "string" ? Promise.reject(o) : o);

export const remove = id => fetch(`api/mods/${id}`, {method: "DELETE"});
