
import { get_package, get_dependencies } from "../../api/cdb.js";
import { add } from "../../service/mods.js";
import store from '../../store/mods.js';

export default {
    data: function() {
        return {
            author: this.$route.params.author,
            name: this.$route.params.name,
            pkg: null,
            deps: null
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
        is_installed: function() {
            return store.list && store.list.find(m => m.name == this.name && m.author == this.author);
        }
    },
    methods: {
        install: function() {
            add({
				name: this.pkg.name,
                author: this.pkg.author,
				mod_type: this.pkg.type,
				source_type: "cdb"
			});
        }
    },
    template: /*html*/`
    <div v-if="pkg">
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
                            <img v-for="screenshot in thumbnails" class="img-thumbnail" :src="screenshot"/>
                        </div>
                        <span v-for="tag in pkg.tags" style="margin: 2px;" class="badge bg-success">{{tag}}</span>
                        <hr>
                        <h4>Description</h4>
                        <pre>{{pkg.long_description}}</pre>
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
                        <a :href="cdb_link" class="btn btn-secondary" target="new">
                            <i class="fa-solid fa-box-open"></i>
                            View on ContentDB
                        </a>
                        <hr>
                        <a :href="pkg.repo" class="btn btn-secondary" target="new" v-if="pkg.repo">
                            <i class="fa-brands fa-git-alt"></i>
                            View source
                        </a>
                        <hr>
                        <button class="btn btn-success" v-on:click="install" :disabled="is_installed">
                            <i class="fa-solid fa-plus"></i>
                            Install
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
    `
};