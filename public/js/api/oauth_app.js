
export const get_all = () => fetch(`api/oauth_app`).then(r => r.json());

export const get_by_id = id => fetch(`api/oauth_app/${id}`).then(r => r.json());

export const save = (q) => fetch(`api/oauth_app`, {
    method: "POST",
    body: JSON.stringify(q)
}).then(r => r.json());

export const remove = id => fetch(`api/oauth_app/${id}`, {
    method: "DELETE"
});
