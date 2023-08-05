
export const get_versions = () => fetch(`api/service/engine/versions`).then(r => r.json());

export const get_status = () => fetch(`api/service/engine`).then(r => r.json());

export const create = opts => fetch(`api/service/engine`, {
    method: "POST",
    body: JSON.stringify(opts)
});

export const start = () => fetch(`api/service/engine/start`, {
    method: "POST"
});

export const stop = () => fetch(`api/service/engine/stop`, {
    method: "POST"
});

export const remove = () => fetch(`api/service/engine`, {
    method: "DELETE"
});

export const get_logs = (since, until) => fetch(`api/service/engine/logs/${since}/${until}`).then(r => r.json());