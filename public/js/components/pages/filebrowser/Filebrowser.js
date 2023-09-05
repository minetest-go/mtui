import DefaultLayout from "../../layouts/DefaultLayout.js";
import { browse, get_zip_url, get_download_url, mkdir, remove, upload } from "../../../api/filebrowser.js";
import format_size from "../../../util/format_size.js";
import format_time from "../../../util/format_time.js";
import { START, FILEBROWSER } from "../../Breadcrumb.js";

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            result: null,
            mkdir_name: "",
            mkfile_name: "",
            prepare_delete: null
        };
    },
    methods: {
        format_size: format_size,
        format_time: format_time,
        get_zip_url: get_zip_url,
        get_download_url: get_download_url,
        mkdir: function() {
            mkdir(this.result.dir + "/" + this.mkdir_name)
            .then(() => this.browse_dir())
            .then(() => this.mkdir_name = "");
        },
        mkfile: function() {
            upload(this.result.dir + "/" + this.mkfile_name, "")
            .then(() => this.browse_dir())
            .then(() => this.mkfile_name = "");
        },
        upload: function() {
            const files = this.$refs.input_upload.files;
            const promises = [];
            for (let i=0; i<files.length; i++) {
                const file = files[i];
                const p = file.arrayBuffer()
                .then(buf => upload(this.result.dir + "/" + file.name, buf))
                promises.push(p);
            }
            Promise.all(promises).then(() => {
                this.$refs.input_upload.value = null;
                this.browse_dir();
            });
        },
        confirm_delete: function() {
            remove(this.result.dir + "/" + this.prepare_delete)
            .then(() => this.prepare_delete = null)
            .then(() => this.browse_dir());
        },
        browse_dir: function() {
            const dir = "/" + this.$route.params.pathMatch;
            browse(dir)
            .then(r => this.result = r)
            .then(() => {
                if (this.result.dir == "/") {
                    this.result.dir = "";
                }
            });
        },
        can_edit: function(filename) {
            return filename.match(/.*(js|lua|txt|conf|cfg|json|md)$/i);
        },
        get_icon: function(item) {
            if (item.is_dir) {
                return "folder";
            } else if (item.name.match(/.*(txt|conf|cfg|md)$/i)) {
                return "file-lines";
            } else if (item.name.match(/.*(js|lua|json)$/i)) {
                return "file-code";
            } else if (item.name.match(/.*(sqlite)$/i)) {
                return "database";
            } else {
                return "file";
            }
        },
        get_icon_class: function(item) {
            const icon = this.get_icon(item);
            return {fa:true, ["fa-"+icon]: true};
        }
    },
    computed: {
        breadcrumb: function() {
            const bc = [START, FILEBROWSER];
            const parts = this.$route.params.pathMatch.split("/");

            let path = "";
            parts
            .filter(p => p != "")
            .forEach(p => {
                if (path == "") {
                    path = p;
                } else {
                    path = path + "/" + p;
                }

                bc.push({
                    name: p,
                    icon: "folder-open",
                    link: "/filebrowser/" + path
                });
            });

            return bc;
        }
    },
    mounted: function() {
        this.browse_dir();
    },
    watch: {
        "$route.params.pathMatch": "browse_dir"
    },
    template: /*html*/`
        <default-layout icon="folder" title="Filebrowser" :breadcrumb="breadcrumb">
            <div class="row">
                <div class="col-2">
                    <div class="input-group">
                        <input type="text" v-model="mkdir_name" class="form-control" placeholder="Directory name"/>
                        <button class="btn btn-secondary" v-on:click="mkdir" :disabled="!mkdir_name">
                            <i class="fa fa-folder"></i>
                            <i class="fa fa-plus"></i>
                            Create directory
                        </button>
                    </div>
                </div>
                <div class="col-2">
                    <div class="input-group">
                        <input type="text" v-model="mkfile_name" class="form-control" placeholder="Filename"/>
                        <button class="btn btn-secondary" v-on:click="mkfile" :disabled="!mkfile_name">
                            <i class="fa fa-file"></i>
                            <i class="fa fa-plus"></i>
                            Create file
                        </button>
                    </div>
                </div>
                <div class="col-2">
                    <div class="input-group">
                        <input ref="input_upload" type="file" class="form-control" multiple/>
                        <button class="btn btn-secondary" v-on:click="upload">
                            <i class="fa fa-upload"></i>
                            Upload file
                        </button>
                    </div>
                </div>
                <div class="col-2">
                </div>
                <div class="col-2">
                </div>
                <div class="col-2" v-if="result">
                    <a class="btn btn-secondary w-100" :href="get_zip_url(result.dir)">
                        <i class="fa fa-download"></i>
                        Download as zip
                    </a>
                </div>
            </div>
            <table class="table table-striped table-condensed" v-if="result">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Size</th>
                        <th>Modification time</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="result.parent_dir">
                        <td>
                            <router-link :to="'/filebrowser' + result.parent_dir">
                                <i class="fa fa-folder-open"></i>
                                Parent dir
                            </router-link>
                        </td>
                        <td></td>
                        <td></td>
                        <td></td>
                    </tr>
                    <tr v-for="item in result.items">
                        <td>
                            <router-link :to="'/filebrowser' + result.dir + '/' + item.name" v-if="item.is_dir">
                                <i v-bind:class="get_icon_class(item)"></i>
                                {{item.name}}
                            </router-link>
                            <span v-if="!item.is_dir">
                                <i v-bind:class="get_icon_class(item)"></i>
                                {{item.name}}
                            </span>
                        </td>
                        <td>
                            <span v-if="!item.is_dir">
                                {{format_size(item.size)}}
                            </span>
                        </td>
                        <td>
                            <span v-if="!item.is_dir">
                                {{format_time(item.mtime)}}
                            </span>
                        </td>
                        <td>
                            <div class="btn-group">
                                <router-link :to="'/fileedit/' + result.dir + '/' + item.name" class="btn btn-sm btn-primary" v-bind:class="{disabled:!can_edit(item.name)}">
                                    <i class="fa fa-edit"></i>
                                </router-link>
                                <a class="btn btn-sm btn-secondary" :href="get_download_url(result.dir + '/' + item.name)" v-bind:class="{disabled:item.is_dir}">
                                    <i class="fa fa-download"></i>
                                </a>
                                <button class="btn btn-sm btn-danger" v-on:click="prepare_delete = item.name">
                                    <i class="fa fa-trash"></i>
                                </button>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
            <div class="modal show" style="display: block;" tabindex="-1" v-show="prepare_delete">
                <div class="modal-dialog">
                    <div class="modal-content">
                    <div class="modal-header">
                        <h1 class="modal-title fs-5">Confirm deletion</h1>
                        <button type="button" class="btn-close" v-on:click="prepare_delete = null"></button>
                    </div>
                    <div class="modal-body">
                        Confirm deletion of <i>{{prepare_delete}}</i>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" v-on:click="prepare_delete = null">Close</button>
                        <button type="button" class="btn btn-danger" v-on:click="confirm_delete">
                            <i class="fa fa-trash"></i>
                            Confirm deletion
                        </button>
                    </div>
                    </div>
                </div>
            </div>
        </default-layout>
    `
};