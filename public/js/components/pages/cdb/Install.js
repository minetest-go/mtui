import DefaultLayout from "../../layouts/DefaultLayout.js";
import { add } from "../../../service/mods.js";
import { resolve_package } from "../../../api/cdb.js";
import { validate } from "../../../api/mods.js";
import { START, ADMINISTRATION, MODS, CDB, CDB_DETAIL } from "../../Breadcrumb.js";
import CDBPackageLink from "../../CDBPackageLink.js";

const DependencyInstallRow = {
    props: ["dep", "selected_dep"],
    computed: {
        no_candidate: function(){
            return this.dep.choices.length == 0 && !this.dep.installed;
        },
        is_installed: function(){
            return this.dep.choices.length == 0 && this.dep.installed;
        },
        has_choices: function() {
            return this.dep.choices.length > 0;
        }
    },
    template: /*html*/`
    <tr v-bind:class="{'table-warning': no_candidate}">
        <td>{{dep.name}}</td>
        <td>
            <select class="form-control" v-on:change="$emit('select_dep', dep.name, $event.target.value)" v-if="has_choices">
                <option v-for="choice in dep.choices" :selected="selected_dep == choice">{{choice}}</option>
            </select>
            <span class="badge bg-danger" v-if="no_candidate">
                <i class="fa-solid fa-triangle-exclamation"></i>
                No installation candidate found!
            </span>
            <span class="badge bg-success" v-if="is_installed">
                <i class="fa fa-check"></i>
                Already installed
            </span>
        </td>
    </tr>
    `
};

export default {
    props: ["author", "name"],
    components: {
        "default-layout": DefaultLayout,
        "cdb-package-link": CDBPackageLink,
        "dependency-install-row": DependencyInstallRow
    },
    data: function() {
        return {
            busy: false,
            missing_dep: false,
            selected_package_map: {}, // modname -> packagename
            installed_mods: [],
            deps: [],
            breadcrumb: [START, ADMINISTRATION, MODS, CDB, CDB_DETAIL(this.author, this.name), {
                name: `Install`,
                icon: "plus",
                link: `/cdb/install/${this.author}/${this.name}`
            }]
        };
    },
    created: function() {
        validate()
        .then(r => this.installed_mods = r.installed)
        .then(() => this.fetch_dependencies());
    },
    methods: {
        fetch_dependencies: function() {
            this.missing_dep = false;
            this.busy = true;

            const selected_packages = [];
            Object.keys(this.selected_package_map)
            .forEach(modname => {
                selected_packages.push(this.selected_package_map[modname]);
            });

            return resolve_package({
                package: `${this.author}/${this.name}`,
                installed_mods: this.installed_mods,
                selected_packages: selected_packages
            })
            .then(deps => this.deps = deps)
            .then(() => {
                this.busy = false;
                this.deps.forEach(d => {
                    if (!d.installed && d.choices.length == 0) {
                        // not installed and no installation candidate
                        this.missing_dep = true;
                    }
                });
            });
        },
        install_next: function() {
            const next_dep = this.deps.find(d => !d.installed);
            if (!next_dep) {
                return;
            }
            const parts = next_dep.selected.split("/");

            add({
                author: parts[0],
				name: parts[1],
				mod_type: "mod",
				source_type: "cdb"
			})
            .then(() => next_dep.installed = true)
            .then(() => this.install_next());
        },
        select_dep: function(modname, dep) {
            this.selected_package_map[modname] = dep;
            this.fetch_dependencies();
        }
    },
    template: /*html*/`
        <default-layout :breadcrumb="breadcrumb" title="Install package" icon="plus">
            <div class="alert alert-primary" v-if="busy">
                <i class="fa fa-spinner fa-spin"></i>
                Calculating dependency-tree...
            </div>
            <table class="table table-striped" v-if="!busy">
                <thead>
                    <tr>
                        <th>Modname</th>
                        <th>Provided by package</th>
                    </tr>
                </thead>
                <tbody v-if="author">
                    <tr>
                        <td>{{name}}</td>
                        <td>
                            <cdb-package-link :author="author" :name="name"/>
                        </td>
                    </tr>
                </tbody>
                <tbody v-for="dep in deps">
                    <dependency-install-row :dep="dep" :selected_dep="dep.selected" v-on:select_dep="select_dep"/>
                </tbody>
            </table>
            <div class="row">
                <div class="col-4"></div>
                <div class="col-4"></div>
                <div class="col-4">
                    <button class="btn btn-success w-100" :disabled="busy || missing_dep" v-on:click="install_next">
                        <i class="fa fa-plus"></i>
                        Install
                    </button>
                </div>
            </div>
        </default-layout>
    `
};