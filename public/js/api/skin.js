
export const upload = data => fetch(`api/skin`, {
    method: "POST",
    body: data
});

export const get = () => fetch(`api/skin`).then(r => r.blob());