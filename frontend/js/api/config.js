
export const get_config = key => fetch(`api/config/${key}`).then(r => r.text());

export const set_config = (key, value) => fetch(`api/config/${key}`, {
    method: "POST",
    body: value
});