
export const upload = (id, data) => fetch(`api/skin/${id}`, {
    method: "POST",
    body: data
});

export const get = id => fetch(`api/skin/${id}`).then(r => r.blob());

export const remove = id => fetch(`api/skin/${id}`, {
    method: "DELETE"
});
