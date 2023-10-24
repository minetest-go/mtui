import DefaultLayout from "../../layouts/DefaultLayout.js";
import { browse, get_zip_url, get_download_url, mkdir, remove, upload, upload_zip, rename } from "../../../api/filebrowser.js";
import format_size from "../../../util/format_size.js";
import format_time from "../../../util/format_time.js";
import { START, FILEBROWSER } from "../../Breadcrumb.js";
import { get_maintenance } from "../../../service/stats.js";
import { can_edit } from "./common.js";

export default {
    props: ["pathMatch"],
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            result: null,
            mkfile_name: "",
            move_name: "",
            move_target: "",
            upload_busy: false,
            upload_zip_busy: false,
            prepare_delete: null
        };
    },
    methods: {
        get_maintenance: get_maintenance,
        format_size: format_size,
        format_time: format_time,
        get_zip_url: get_zip_url,
        get_download_url: get_download_url,
        mkdir: function() {
            mkdir(this.result.dir + "/" + this.mkfile_name)
            .then(() => this.browse_dir())
            .then(() => this.mkfile_name = "");
        },
        mkfile: function() {
            upload(this.result.dir + "/" + this.mkfile_name, "")
            .then(() => this.browse_dir())
            .then(() => this.mkfile_name = "");
        },
        upload: function() {
            const files = Array.from(this.$refs.input_upload.files);
            this.upload_busy = true;

            const promises = files.map(file => {
                return file.arrayBuffer()
                .then(buf => upload(this.result.dir + "/" + file.name, buf));
            });

            Promise.all(promises).then(() => {
                this.$refs.input_upload.value = null;
                this.upload_busy = false;
                this.browse_dir();
            });
        },
        upload_zip: function() {
            if (this.$refs.input_upload_zip.files.length == 0) {
                return;
            }
            this.upload_zip_busy = true;

            const file = this.$refs.input_upload_zip.files[0];
            file.arrayBuffer()
            .then(buf => upload_zip(this.result.dir, buf))
            .then(() => {
                this.$refs.input_upload_zip.value = null;
                this.upload_zip_busy = false;
                this.browse_dir();
            });
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
        is_database: function(filename) {
            return filename.match(/.*(sqlite|sqlite-shm|sqlite-wal)$/i);
        },
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
                <div class="col-4">
                    <div class="input-group">
                        <input type="text" v-model="mkfile_name" class="form-control" placeholder="Directory name"/>
                        <button class="btn btn-secondary" v-on:click="mkdir" :disabled="!mkfile_name">
                            <i class="fa fa-folder"></i>
                            <i class="fa fa-plus"></i>
                            Create directory
                        </button>
                        <button class="btn btn-secondary" v-on:click="mkfile" :disabled="!mkfile_name">
                            <i class="fa fa-file"></i>
                            <i class="fa fa-plus"></i>
                            Create file
                        </button>
                    </div>
                </div>
                <div class="col-3">
                    <div class="input-group">
                        <input ref="input_upload" type="file" class="form-control" multiple/>
                        <button class="btn btn-secondary" v-on:click="upload">
                            <i class="fa fa-upload"></i>
                            Upload file
                            <i class="fa fa-spinner fa-spin" v-if="upload_busy"></i>
                        </button>
                    </div>
                </div>
                <div class="col-3">
                    <div class="input-group">
                        <input ref="input_upload_zip" type="file" class="form-control" accept=".zip"/>
                        <button class="btn btn-secondary" v-on:click="upload_zip">
                            <i class="fa fa-upload"></i>
                            Upload zip
                            <i class="fa-solid fa-triangle-exclamation" style="color: orange;" title="The contents of the zip-file will overwrite files with the same name!"></i>
                            <i class="fa fa-spinner fa-spin" v-if="upload_zip_busy"></i>
                        </button>
                    </div>
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
                    <tr v-for="item in result.items" v-bind:class="{'table-warning': is_database(item.name) && !get_maintenance()}">
                        <td>
                            <router-link :to="'/filebrowser' + result.dir + '/' + item.name" v-if="item.is_dir">
                                <i v-bind:class="get_icon_class(item)"></i>
                                {{item.name}}
                            </router-link>
                            <span v-if="!item.is_dir">
                                <i v-bind:class="get_icon_class(item)"></i>
                                {{item.name}}
                                <i class="fa-solid fa-triangle-exclamation" v-if="is_database(item.name) && !get_maintenance()" title="Database might be inconsistent when downloading without maintenance-mode enabled!"></i>
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