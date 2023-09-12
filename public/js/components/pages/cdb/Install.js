import DefaultLayout from "../../layouts/DefaultLayout.js";
import { add } from "../../../service/mods.js";
import { get_dependencies, get_package } from "../../../service/cdb.js";
import { search_packages } from "../../../api/cdb.js";
import { validate } from "../../../api/mods.js";
import { START, ADMINISTRATION, MODS, CDB, CDB_DETAIL } from "../../Breadcrumb.js";
import CDBPackageLink from "../../CDBPackageLink.js";

const DependencyInstallRow = {
    props: ["dep", "selected_dep"],
    methods: {
        select_dep: function(modname, dep) {
            if (modname == "") {
                delete this.selected_deps[modname];
            } else {
                const parts = dep.split("/");
                this.selected_deps[modname] = {
                    author: parts[0],
                    name: parts[1]
                };
            }
        },
    },
    template: /*html*/`
    <tr>
        <td>{{dep.name}}</td>
        <td>
            <select class="form-control" v-on:change="select_dep(dep.name, $event.target.value)" v-if="dep.choices.length > 0">
                <option v-for="choice in dep.choices" :selected="selected_dep == choice">{{choice}}</option>
            </select>
            <span class="badge bg-danger" v-if="dep.choices.length == 0 && !dep.installed">
                <i class="fa-solid fa-triangle-exclamation"></i>
                No installation candidate found!
            </span>
            <span class="badge bg-success" v-if="dep.choices.length == 0 && dep.installed">
                <i class="fa fa-check"></i>
                Already installed
            </span>
        </td>
    </tr>
    `
};

export default {
    components: {
        "default-layout": DefaultLayout,
        "cdb-package-link": CDBPackageLink,
        "dependency-install-row": DependencyInstallRow
    },
    data: function() {
        const author = this.$route.params.author;
        const name = this.$route.params.name;

        return {
            author: author,
            name: name,
            pkg: null,
            selected_deps: {}, // modname => {author,name}
            deps: [],
            packages: [],
            installed_mods: {}, // modname => true
            breadcrumb: [START, ADMINISTRATION, MODS, CDB, CDB_DETAIL(author, name), {
                name: `Install`,
                icon: "plus",
                link: `/cdb/install/${author}/${name}`
            }]
        };
    },
    created: function() {
        get_package(this.author, this.name)
        .then(p => this.pkg = p);

        validate()
        .then(r => r.installed.forEach(modname => this.installed_mods[modname] = true))
        .then(() => search_packages({ type: ["mod"] }))
        .then(pkgs => this.packages = pkgs)
        .then(() => this.resolve_deps(this.author, this.name));
    },
    methods: {
        get_key: function(author, name) {
            return `${author}/${name}`;
        },
        get_author_name: function(pgkname) {
            return pgkname.split("/");
        },
        fetch_dependencies: function(author, name) {
            const key = this.get_key(author, name);
            if (this.deps[key]) {
                // already fetched
                return Promise.resolve(this.deps[key]);
            } else {
                // fetch
                return get_dependencies(author, name)
                .then(deps => Object.keys(deps).forEach(k => this.deps[k] = deps[k]))
                .then(() => this.deps[key]);
            }
        },
        resolve_deps: function(author, name) {
            // fetch dependency info
            this.fetch_dependencies(author, name)
            .then(deps => {
                deps
                .filter(dep => !dep.is_optional) // not-optional
                .forEach(dep => {
                    if (this.installed_mods[dep.name]) {
                        // already installed
                        this.deps.push({
                            name: dep.name,
                            choices: [],
                            installed: true
                        });
                        return;
                    }

                    // fetch all package infos and provide package choices
                    const choices = [];
                    dep.packages.forEach(pkgname => {
                        const detail = this.packages.find(p => this.get_key(p.author, p.name) == pkgname);
                        if (detail) {
                            // mod found
                            choices.push(pkgname);
                        }
                    });

                    this.deps.push({
                        name: dep.name,
                        choices: choices
                    });

                    if (choices.length > 0) {
                        // select first choice as fallback
                        let choice = choices[0];
                        choices.forEach(c => {
                            const author_name = this.get_author_name(c);
                            if (author_name[1] == dep.name) {
                                // exact match
                                choice = c;
                            }
                        });

                        this.selected_deps[dep.name] = choice;
                        const author_name = this.get_author_name(choice);
                        this.fetch_dependencies(author_name[0], author_name[1]);
                    }
                });
            });
        },
        install: function() {
            return add({
				name: this.pkg.name,
                author: this.pkg.author,
				mod_type: this.pkg.type,
				source_type: "cdb"
			});
        }
    },
    template: /*html*/`
        <default-layout :breadcrumb="breadcrumb" title="Install package" icon="plus">
            <table class="table table-striped">
                <thead>
                    <tr>
                        <th>Modname</th>
                        <th>Provided by package</th>
                    </tr>
                </thead>
                <tbody v-if="pkg">
                    <tr>
                        <td>{{pkg.name}}</td>
                        <td>
                            <cdb-package-link :pkg="pkg"/>
                        </td>
                    </tr>
                </tbody>
                <tbody v-for="dep in deps">
                    <dependency-install-row :dep="dep" :selected_dep="selected_deps[dep.name]"/>
                </tbody>
            </table>
        </default-layout>
    `
};