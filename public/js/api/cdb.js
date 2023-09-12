import { protected_fetch } from "./util.js";

export const search_packages = q => protected_fetch(`api/cdb/search`, {
    method: "POST",
    body: JSON.stringify(q)
});

export const resolve_package = data => protected_fetch(`api/cdb/resolve`, {
    method: "POST",
    body: JSON.stringify(data)
});

export const get_package = (author, name) => protected_fetch(`api/cdb/detail/${author}/${name}`);

export const get_dependencies = (author, name) => protected_fetch(`api/cdb/detail/${author}/${name}/dependencies`);