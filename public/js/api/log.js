
export const get_events = category => fetch(`api/log/events/${category}`).then(r => r.json());

export const count = s => fetch("api/log/count", {
    method: "POST",
    body: JSON.stringify(s)
})
.then(r => r.json());

export const search = s => fetch("api/log/search", {
    method: "POST",
    body: JSON.stringify(s)
})
.then(r => r.json());