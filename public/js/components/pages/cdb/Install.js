import DefaultLayout from "../../layouts/DefaultLayout.js";
import { add } from "../../../service/mods.js";
import { get_dependencies, get_package } from "../../../service/cdb.js";
import { validate } from "../../../api/mods.js";
import { START, ADMINISTRATION, MODS, CDB, CDB_DETAIL } from "../../Breadcrumb.js";
import CDBPackageLink from "../../CDBPackageLink.js";

export default {
    components: {
        "default-layout": DefaultLayout,
        "cdb-package-link": CDBPackageLink
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
            installed_mods: {}, // modname => true
            breadcrumb: [START, ADMINISTRATION, MODS, CDB, CDB_DETAIL(author, name), {
                name: `Install package '${author}/${name}'`,
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
        .then(this.resolve_deps(this.author, this.name));
    },
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
        resolve_deps: function(author, name) {
            const key = `${author}/${name}`;

            // fetch dependency info
            get_dependencies(author, name)
            .then(deps => {
                deps[key]
                .filter(dep => !dep.is_optional) // not-optional
                .filter(dep => !this.installed_mods[dep.name]) // not installed
                .forEach(dep => {
                    if (this.installed_mods[dep.name]) {
                        // already installed
                        this.deps.push({
                            name: deps.name,
                            installed: true
                        });
                    }

                    // fetch all package infos and provide package choices
                    const pl = dep.packages.map(p => {
                        const parts = p.split("/");
                        return get_package(parts[0], parts[1]);
                    });

                    Promise.all(pl)
                    .then(package_choices => package_choices.filter(pc => pc.type == "mod"))
                    .then(package_choices => {
                        this.deps.push({
                            name: dep.name,
                            choices: package_choices.map(p => `${p.author}/${p.name}`)
                        });
                    });
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
                        <th>Installed</th>
                    </tr>
                </thead>
                <tbody v-if="pkg">
                    <tr>
                        <td>{{pkg.name}}</td>
                        <td>
                            <cdb-package-link :pkg="pkg"/>
                        </td>
                        <td></td>
                    </tr>
                </tbody>
                <tbody v-for="dep in deps">
                    <tr>
                        <td>{{dep.name}}</td>
                        <td>
                            <select class="form-control" v-on:change="select_dep(dep.name, $event.target.value)">
                                <option></option>
                                <option v-for="choice in dep.choices" :selected="selected_deps[dep.name] == choice">{{choice}}</option>
                            </select>
                        </td>
                        <td>
                            <i class="fa fa-check" v-if="dep.installed"></i>
                        </td>
                    </tr>
                </tbody>
            </table>
        </default-layout>
    `
};