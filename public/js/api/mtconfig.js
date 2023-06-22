
export const get_all = () => fetch("api/mtconfig").then(r => r.json());

export const set = (key, value) => fetch(`api/mtconfig/${key}`, {
    method: "POST",
    body: value
}).then(r => r.json());

export const remove = key => fetch(`api/mtconfig/${key}`, {
    method: "DELETE",
});