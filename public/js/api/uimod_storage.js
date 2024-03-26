
export const get_uimod_storage = key => fetch(`api/uimod/storage/${key}`).then(r => r.text());

export const set_uimod_storage = (key, value) => fetch(`api/uimod/storage/${key}`, {
    method: "POST",
    body: value
}).then(r => r.text());
