export const list = () => fetch("api/mods").then(r => r.json());

export const scan = () => fetch("api/mods/scan", {method: "POST"}).then(r => r.json());

export const create = mod => fetch("api/mods", {
    method: "POST",
    body: JSON.stringify(mod)
}).then(r => r.json());

export const remove = id => fetch(`api/mods/${id}`, {method: "DELETE"});
