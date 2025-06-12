import DefaultLayout from "../../layouts/DefaultLayout.js";
import CDBPackageLink from "../../CDBPackageLink.js";

import { add, get_cdb_mod } from "../../../service/mods.js";
import { update_settings } from "../../../service/mtconfig.js";
import { resolve_package, get_package } from "../../../api/cdb.js";
import { validate } from "../../../api/mods.js";
import { START, ADMINISTRATION, MODS, CDB, CDB_DETAIL } from "../../Breadcrumb.js";

const DependencyInstallRow = {
    props: ["dep", "selected_dep"],
    computed: {
        no_candidate: function(){
            return this.dep.choices.length == 0 && !this.dep.installed;
        },
        is_installed: function(){
            return this.dep.installed;
        },
        has_choices: function() {
            return this.dep.choices.length > 0 && !this.dep.installed;
        }
    },
    template: /*html*/`
    <tr v-bind:class="{'table-warning': no_candidate}">
        <td>{{dep.name}}</td>
        <td>
            <select class="form-control" v-on:change="$emit('select_dep', dep.name, $event.target.value)" v-if="has_choices">
                <option v-for="choice in dep.choices" :selected="selected_dep == choice">{{choice}}</option>
            </select>
            <i class="fa fa-spinner fa-spin" v-if="dep.busy"></i>
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
            install_busy: false,
            busy: false,
            missing_dep: false,
            selected_package_map: {}, // modname -> packagename
            installed_mods: [],
            pkg: null,
            deps: [],
            breadcrumb: [START, ADMINISTRATION, MODS, CDB, CDB_DETAIL(this.author, this.name), {
                name: `Install`,
                icon: "plus",
                link: `/cdb/install/${this.author}/${this.name}`
            }]
        };
    },
    created: async function() {
        const pkg = await get_package(this.author, this.name);
        const mod = get_cdb_mod(pkg.author, pkg.name);
        pkg.installed = !!mod;
        if (pkg.type == "mod") {
            // resolve mod dependencies
            const r = await validate();
            this.installed_mods = r.installed;
            await this.fetch_dependencies();
        }
        this.pkg = pkg;
    },
    methods: {
        fetch_dependencies: async function() {
            this.missing_dep = false;
            this.busy = true;

            const selected_packages = [];
            Object.keys(this.selected_package_map)
            .forEach(modname => {
                selected_packages.push(this.selected_package_map[modname]);
            });

            const deps = await resolve_package({
                package: `${this.author}/${this.name}`,
                installed_mods: this.installed_mods,
                selected_packages: selected_packages
            });
            deps.forEach(d => {
                d.busy = false;
                if (!d.installed && d.choices.length == 0) {
                    // not installed and no installation candidate
                    this.missing_dep = true;
                }
            });
            this.deps = deps;
            this.busy = false;
        },
        install_next: async function() {
            this.install_busy = true;
            const next_dep = this.deps.find(d => !d.installed);

            if (!next_dep && !this.pkg.installed) {
                // install root dep
                await add({
                    author: this.pkg.author,
                    name: this.pkg.name,
                    mod_type: this.pkg.type,
                    source_type: "cdb"
                });
                this.install_busy = false;
                this.pkg.installed = true;

                update_settings();
                // go back to cdb browser
                this.$router.push("/cdb/browse");
                return;
            }

            if (!next_dep) {
                // no more packages
                this.install_busy = false;
                update_settings();
                // go back to cdb browser
                this.$router.push("/cdb/browse");
                return;
            }

            const parts = next_dep.selected.split("/");
            next_dep.busy = true;

            await add({
                author: parts[0],
				name: parts[1],
				mod_type: "mod",
				source_type: "cdb"
			});
            next_dep.installed = true;
            next_dep.busy = false;
            await this.install_next();
        },
        select_dep: async function(modname, dep) {
            this.selected_package_map[modname] = dep;
            await this.fetch_dependencies();
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
                <tbody>
                    <tr v-if="pkg">
                        <td>{{name}}</td>
                        <td>
                            <cdb-package-link :author="author" :name="name"/>
                            <i class="fa fa-check" v-if="pkg.installed"></i>
                        </td>
                    </tr>
                    <dependency-install-row v-for="dep in deps" :dep="dep" :key="dep.name" :selected_dep="dep.selected" v-on:select_dep="select_dep"/>
                </tbody>
            </table>
            <div class="row">
                <div class="col-4"></div>
                <div class="col-4"></div>
                <div class="col-4">
                    <button class="btn btn-success w-100" :disabled="busy || install_busy || missing_dep" v-on:click="install_next">
                        <i class="fa fa-plus"></i>
                        Install
                        <i class="fa fa-spinner fa-spin" v-if="install_busy"></i>
                    </button>
                </div>
            </div>
        </default-layout>
    `
};