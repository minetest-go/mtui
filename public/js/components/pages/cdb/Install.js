import DefaultLayout from "../../layouts/DefaultLayout.js";
import { add } from "../../../service/mods.js";
import { get_dependencies, get_package } from "../../../service/cdb.js";
import { search_packages, resolve_package } from "../../../api/cdb.js";
import { validate } from "../../../api/mods.js";
import { START, ADMINISTRATION, MODS, CDB, CDB_DETAIL } from "../../Breadcrumb.js";
import CDBPackageLink from "../../CDBPackageLink.js";

const DependencyInstallRow = {
    props: ["dep", "selected_dep"],
    methods: {
        select_dep: function(modname, dep) {
            console.log("selected", modname, dep);
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
            selected_packages: [],
            deps: [],
            breadcrumb: [START, ADMINISTRATION, MODS, CDB, CDB_DETAIL(author, name), {
                name: `Install`,
                icon: "plus",
                link: `/cdb/install/${author}/${name}`
            }]
        };
    },
    created: function() {
        validate()
        .then(r => {
            return resolve_package({
                package: `${this.author}/${this.name}`,
                installed_mods: r.installed,
                selected_packages: this.selected_packages
            });
        })
        .then(deps => {
            this.deps = deps;
            console.log(deps);
        });
    },
    methods: {
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
                <tbody v-if="author">
                    <tr>
                        <td>{{author}}/{{name}}</td>
                        <td>
                            <cdb-package-link :author="author" :name="name"/>
                        </td>
                    </tr>
                </tbody>
                <tbody v-for="dep in deps">
                    <dependency-install-row :dep="dep" selected_dep=""/>
                </tbody>
            </table>
        </default-layout>
    `
};