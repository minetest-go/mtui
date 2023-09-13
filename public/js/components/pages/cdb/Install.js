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
            <select class="form-control" v-on:change="$emit('select_dep', $event.target.value)" v-if="has_choices">
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
            installed_mods: [],
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
        .then(r => this.installed_mods = r.installed)
        .then(() => this.fetch_dependencies());
    },
    methods: {
        fetch_dependencies: function() {
            return resolve_package({
                package: `${this.author}/${this.name}`,
                installed_mods: this.installed_mods,
                selected_packages: this.selected_packages
            })
            .then(deps => this.deps = deps);
        },
        install: function() {
            return add({
				name: this.pkg.name,
                author: this.pkg.author,
				mod_type: this.pkg.type,
				source_type: "cdb"
			});
        },
        select_dep: function(dep) {
            if (!this.installed_mods.find(m => m == dep)) {
                this.installed_mods.push(dep);
            }
            this.fetch_dependencies();
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
        </default-layout>
    `
};