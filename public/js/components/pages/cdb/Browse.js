
import { search_packages } from "../../../api/cdb.js";
import Preview from "./Preview.js";
import DefaultLayout from "../../layouts/DefaultLayout.js";
import { START, ADMINISTRATION, MODS, CDB } from "../../Breadcrumb.js";

const store = Vue.reactive({
    busy: false,
    packages: null,
    query: "",
    type: "mod",
    breadcrumb: [START, ADMINISTRATION, MODS, CDB]
});

export default {
    data: () => store,
    created: function() {
        if (!store.packages) {
            this.search();
        }
    },
    components: {
        "package-preview": Preview,
        "default-layout": DefaultLayout
    },
    watch: {
        "type": function() {
            this.search();
        }
    },
    methods: {
        search: function() {
            this.busy = true;
            this.packages = [];
            search_packages({
                type: [store.type],
                query: store.query,
                limit: 25,
                sort: "score",
                order: "desc"
            })
            .then(pkgs => store.packages = pkgs)
            .finally(() => this.busy = false);
        }
    },
    template: /*html*/`
    <default-layout icon="box-open" title="ContentDB" :breadcrumb="breadcrumb">
        <h3>Browse cdb</h3>
        <div class="row">
            <div class="col-2">
                <label>Type</label>
                <select class="form-control" v-model="type">
                    <option value="mod">Mod</option>
                    <option value="game">Game</option>
                    <option value="txp">Texture pack</option>
                </select>
            </div>
            <div class="col-8">
                <label>Keywords</label>
                <input type="text" class="form-control" v-model="query" v-on:keyup.enter="search"/>
            </div>
            <div class="col-md-2">
                <label>Search</label>
                <button class="btn btn-primary w-100" v-on:click="search">
                    <i class="fa fa-magnifying-glass" v-if="!busy"></i>
                    <i class="fa-solid fa-spinner fa-spin" v-else></i>
                    Search
                    <span class="badge bg-secondary" v-if="packages">{{packages.length}}</span>
                </button>
            </div>
        </div>
        <hr>
        <div style="display: flex; flex-wrap: wrap;">
            <package-preview v-for="pkg in packages" v-if="packages" :pkg="pkg"></package-preview>
        </div>
    </default-layout>
    `
};