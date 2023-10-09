import { protected_fetch } from "./util.js";

export const get_versions = (servicename) => protected_fetch(`api/service/${servicename}/versions`);

export const get_status = (servicename) => protected_fetch(`api/service/${servicename}`);

export const create = (servicename, opts) => protected_fetch(`api/service/${servicename}`, {
    method: "POST",
    body: JSON.stringify(opts)
});

export const start = (servicename) => protected_fetch(`api/service/${servicename}/start`, {
    method: "POST"
});

export const stop = (servicename) => protected_fetch(`api/service/${servicename}/stop`, {
    method: "POST"
});

export const remove = (servicename) => protected_fetch(`api/service/${servicename}`, {
    method: "DELETE"
});

export const get_logs = (servicename, since, until) => protected_fetch(`api/service/${servicename}/logs/${since}/${until}`);