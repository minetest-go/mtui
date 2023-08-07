
export const get_packages = type => fetch(`api/cdb/packages/${type}`).then(r => r.json());

