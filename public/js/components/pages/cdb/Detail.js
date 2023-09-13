import { get_dependencies, get_package } from "../../../service/cdb.js";
import { get_cdb_mod, get_game } from "../../../service/mods.js";
import FeedbackButton from "../../FeedbackButton.js";
import DefaultLayout from "../../layouts/DefaultLayout.js";
import { START, ADMINISTRATION, MODS, CDB, CDB_DETAIL } from "../../Breadcrumb.js";

export default {
    props: ["author", "name"],
    components: {
        "feedback-button": FeedbackButton,
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            pkg: null,
            deps: null,
            breadcrumb: [START, ADMINISTRATION, MODS, CDB, CDB_DETAIL(this.author, this.name)]
        };
    },
    created: function() {
        get_package(this.author, this.name)
        .then(p => this.pkg = p);

        get_dependencies(this.author, this.name)
        .then(d => this.deps = d[`${this.author}/${this.name}`]);
    },
    computed: {
        thumbnails: function() {
            return this.pkg.screenshots.map(s => s.replaceAll("/uploads/", "/thumbnails/2/"));
        },
        cdb_link: function() {
            return `https://content.minetest.net/packages/${this.pkg.author}/${this.pkg.name}/`;
        },
        install_disabled: function() {
            return get_cdb_mod(this.author, this.name) || (this.pkg && this.pkg.type == "game" && get_game());
        }
    },
    template: /*html*/`
    <default-layout v-if="pkg" icon="box-open" title="Package detail" :breadcrumb="breadcrumb">
        <h4>
            {{pkg.title}}
            <small class="text-muted">by {{pkg.author}}</small>
        </h4>
        <div class="row">
            <div class="col-10">
                <div class="card">
                    <div class="card-header">
                        Package details
                    </div>
                    <div class="card-body">
                        <div>
                            <img v-for="screenshot in thumbnails" :src="screenshot" style="margin: 5px;"/>
                        </div>
                        <span v-for="tag in pkg.tags" style="margin: 2px;" class="badge bg-success">{{tag}}</span>
                        <hr>
                        <h4>Description</h4>
                        <pre>{{pkg.long_description || pkg.short_description}}</pre>
                        <h4>Dependencies</h4>
                        <ul v-if="deps">
                            <li v-for="dep in deps">
                                {{dep.name}}
                                <span class="badge bg-info" v-if="dep.is_optional">optional</span>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
            <div class="col-2">
                <div class="card">
                    <div class="card-header">
                        Actions
                    </div>
                    <div class="card-body">
                        <a :href="cdb_link" class="btn btn-secondary" target="_blank">
                            <i class="fa-solid fa-box-open"></i>
                            View on ContentDB
                        </a>
                        <hr>
                        <a :href="pkg.repo" class="btn btn-secondary" target="_blank" v-if="pkg.repo">
                            <i class="fa-brands fa-git-alt"></i>
                            View source
                        </a>
                        <hr>
                        <router-link class="btn btn-success" :to="'/cdb/install/' + pkg.author + '/' + pkg.name" v-if="!install_disabled">
                            <i class="fa-solid fa-plus"></i>
                            Install
                        </router-link>
                    </div>
                </div>
            </div>
        </div>
    </default-layout>
    `
};