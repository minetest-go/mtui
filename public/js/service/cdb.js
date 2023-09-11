// functions for cached cdb access
import { get_dependencies as cdb_get_dependencies, get_package as cdb_get_package } from "../api/cdb.js";

const get_key = (author, name) => `${author}/${name}`;
const package_cache = {};
const package_dependency_cache = {};

export const get_package = (author, name) => {
    const key = get_key(author, name);
    if (package_cache[key]) {
        // cached
        return Promise.resolve(package_cache[key]);
    }
    // fetch from api
    return cdb_get_package(author, name)
    .then(p => {
        // populate cache
        package_cache[key] = p;
        return p;
    });
};

export const get_dependencies = (author, name) => {
    const key = get_key(author, name);
    if (package_dependency_cache[key]){
        // cached
        return Promise.resolve(package_dependency_cache[key]);
    }
    // fetch from api
    return cdb_get_dependencies(author, name)
    .then(d => {
        // populate cache
        package_dependency_cache[key] = d;
        return d;
    });
};
