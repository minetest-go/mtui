
export const search_packages = q => fetch(`api/cdb/search`, {
    method: "POST",
    body: JSON.stringify(q)
}).then(r => r.json());

export const get_package = (author, name) => fetch(`api/cdb/detail/${author}/${name}`)
.then(r => r.json());

export const get_dependencies = (author, name) => fetch(`api/cdb/detail/${author}/${name}/dependencies`)
.then(r => r.json());