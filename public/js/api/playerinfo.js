export const get = name => fetch(`api/player/info/${name}`).then(r => r.json());

export const search = q => fetch(`api/player/search`, {
    method: "POST",
    body: JSON.stringify(q)
})
.then(r => r.json());

export const count = q => fetch(`api/player/count`, {
    method: "POST",
    body: JSON.stringify(q)
})
.then(r => r.json());