import DefaultLayout from "../../layouts/DefaultLayout.js";
import { START, FILEBROWSER } from "../../Breadcrumb.js";
import CDBPackageLink from "../../CDBPackageLink.js";

import { browse, get_download_url, mkdir, remove, rename } from "../../../api/filebrowser.js";

import { upload_chunked } from "../../../service/uploader.js";
import { get_mod_by_name, get_txp_by_name, get_game } from "../../../service/mods.js";
import format_size from "../../../util/format_size.js";
import format_time from "../../../util/format_time.js";
import { can_edit } from "./common.js";

export default {
    props: ["pathMatch"],
    components: {
        "default-layout": DefaultLayout,
        "cdb-package-link": CDBPackageLink
    },
    data: function() {
        return {
            result: null,
            mkdir_name: "",
            move_name: "",
            move_target: "",
            prepare_delete: null,
            upload_progress: {}
        };
    },
    methods: {
        format_size,
        format_time,
        get_mod_by_name,
        get_txp_by_name,
        get_game,
        get_download_url,
        mkdir: function() {
            mkdir(this.result.dir + "/" + this.mkdir_name)
            .then(() => this.browse_dir())
            .then(() => this.mkdir_name = "");
        },
        upload: function() {
            const files = Array.from(this.$refs.input_upload.files);
            files.forEach(file => {
                upload_chunked(this.result.dir, file.name, file, progress => {
                    this.upload_progress[file.name] = {
                        progress,
                        name: file.name,
                        size: file.size
                    };
                }).then(() => {
                    delete this.upload_progress[file.name];
                    this.browse_dir();
                });
            });

            this.$refs.input_upload.value = null;
        },
        confirm_delete: function() {
            remove(this.result.dir + "/" + this.prepare_delete)
            .then(() => this.prepare_delete = null)
            .then(() => this.browse_dir());
        },
        confirm_move: function() {
            rename(this.result.dir + "/" + this.move_name, this.result.dir + "/" + this.move_target)
            .then(() => {
                this.move_name = "";
                this.move_target = "";
                this.browse_dir();
            });
        },
        browse_dir: function() {
            const dir = "/" + this.pathMatch;
            browse(dir)
            .then(r => this.result = r)
            .then(() => {
                if (this.result.dir == "/") {
                    this.result.dir = "";
                }
            });
        },
        can_edit: can_edit,
        is_json_profile: function(filename) {
            return filename.match(/^profile-.*.json$/);
        },
        get_icon: function(item) {
            if (item.is_dir) {
                return "folder";
            } else if (can_edit(item.name)) {
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
            const parts = this.pathMatch.split("/");

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
        "pathMatch": "browse_dir"
    },
    template: /*html*/`
        <default-layout icon="folder" title="Filebrowser" :breadcrumb="breadcrumb">
            <div class="row">
                <div class="col-md-4">
                    <div class="input-group">
                        <input type="text" v-model="mkdir_name" class="form-control" placeholder="Directory name"/>
                        <button class="btn btn-secondary" v-on:click="mkdir" :disabled="!mkdir_name">
                            <i class="fa fa-folder"></i>
                            <i class="fa fa-plus"></i>
                            Create directory
                        </button>
                    </div>
                </div>
                <div class="col-md-8">
                    <div class="input-group">
                        <input ref="input_upload" type="file" class="form-control" multiple/>
                        <button class="btn btn-secondary" v-on:click="upload">
                            <i class="fa fa-upload"></i>
                            Upload file
                        </button>
                    </div>
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
                    <tr v-for="(entry, name) in upload_progress">
                        <td>
                            <i class="fa fa-spinner fa-spin"></i>
                            {{name}} ({{Math.floor(entry.progress * 100)}} % / {{format_size(entry.size)}})
                        </td>
                        <td colspan="3">
                            <div class="progress">
                                <div class="progress-bar progress-bar-striped progress-bar-animated" v-bind:style="{ width: (entry.progress*100)+'%' }">
                                </div>
                            </div>
                        </td>
                    </tr>
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
                            <span v-if="result.dir == '' && item.name == 'game'" v-for="mod in [get_game()]" class="float-end">
                                <cdb-package-link v-if="mod && mod.source_type == 'cdb'" :author="mod.author" :name="mod.name"/>
                            </span>
                            <span v-if="result.dir == '/worldmods' && item.is_dir" v-for="mod in [get_mod_by_name(item.name)]" class="float-end">
                                <cdb-package-link v-if="mod && mod.source_type == 'cdb'" :author="mod.author" :name="mod.name"/>
                            </span>
                            <span v-if="result.dir == '/textures' && item.is_dir" v-for="mod in [get_txp_by_name(item.name)]" class="float-end">
                                <cdb-package-link v-if="mod && mod.source_type == 'cdb'" :author="mod.author" :name="mod.name"/>
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
                                <router-link :to="'/profiler-view/' + result.dir + '/' + item.name" class="btn btn-sm btn-secondary" v-if="is_json_profile(item.name)">
                                    <i class="fa fa-chart-line"></i>
                                </router-link>
                                <router-link :to="'/fileedit/' + result.dir + '/' + item.name" class="btn btn-sm btn-primary" v-bind:class="{disabled:!can_edit(item.name)}">
                                    <i class="fa fa-edit"></i>
                                </router-link>
                                <a class="btn btn-sm btn-secondary" :href="get_download_url(result.dir + '/' + item.name)" v-bind:class="{disabled:item.is_dir}">
                                    <i class="fa fa-download"></i>
                                </a>
                                <button class="btn btn-sm btn-warning" v-on:click="move_name = item.name; move_target = item.name">
                                    <i class="fa fa-shuffle"></i>
                                </button>
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
                        <button type="button" class="btn btn-success" v-on:click="prepare_delete = null">Abort</button>
                        <button type="button" class="btn btn-danger" v-on:click="confirm_delete">
                            <i class="fa fa-trash"></i>
                            Confirm deletion
                        </button>
                    </div>
                    </div>
                </div>
            </div>
            <div class="modal show" style="display: block;" tabindex="-1" v-show="move_name">
                <div class="modal-dialog">
                    <div class="modal-content">
                    <div class="modal-header">
                        <h1 class="modal-title fs-5">Rename file</h1>
                        <button type="button" class="btn-close" v-on:click="move_name = null"></button>
                    </div>
                    <div class="modal-body">
                        Move file <i>{{move_name}}</i>
                        <input type="text" class="form-control" placeholder="New name" v-model="move_target"/>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" v-on:click="move_name = null">Abort</button>
                        <button type="button" class="btn btn-warning" v-on:click="confirm_move" :disabled="!move_target || move_target == move_name">
                            <i class="fa fa-shuffle"></i>
                            Confirm rename
                        </button>
                    </div>
                    </div>
                </div>
            </div>
        </default-layout>
    `
};