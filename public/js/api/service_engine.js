import { protected_fetch } from "./util.js";

export const get_versions = () => protected_fetch(`api/service/engine/versions`);

export const get_status = () => protected_fetch(`api/service/engine`);

export const create = opts => protected_fetch(`api/service/engine`, {
    method: "POST",
    body: JSON.stringify(opts)
});

export const start = () => protected_fetch(`api/service/engine/start`, {
    method: "POST"
});

export const stop = () => protected_fetch(`api/service/engine/stop`, {
    method: "POST"
});

export const remove = () => protected_fetch(`api/service/engine`, {
    method: "DELETE"
});

export const get_logs = (since, until) => protected_fetch(`api/service/engine/logs/${since}/${until}`);