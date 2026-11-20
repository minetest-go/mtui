<script>
import DefaultLayout from "../../layouts/DefaultLayout.vue";
import ConfirmationPrompt from "../../ConfirmationPrompt.vue";

import { browse, get_download_url, mkdir, remove, rename } from "../../../api/filebrowser.js";

import { upload_chunked } from "../../../service/uploader.js";
import format_size from "../../../util/format_size.js";
import format_time from "../../../util/format_time.js";
import { can_edit, filebrowser_breadcrumb } from "./common.js";
import { confirm } from "../../../store/confirm.js";

export default {
    props: ["pathMatch"],
    components: {
        "default-layout": DefaultLayout,
        "confirmation-prompt": ConfirmationPrompt
    },
    data: function() {
        return {
            result: null,
            mkdir_name: "",
            move_name: "",
            move_target: "",
            upload_progress: {}
        };
    },
    methods: {
        format_size,
        format_time,
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
        delete_item: async function(item) {
            const ok = await confirm({
                title: `Delete ${item.name}?`,
                message: item.is_dir ?
                    `The directory "${item.name}" and everything in it will be deleted permanently.` :
                    `The file "${item.name}" will be deleted permanently.`,
                confirm_label: "Delete",
                confirm_icon: "fa fa-trash"
            });
            if (!ok) {
                return;
            }
            remove(this.result.dir + "/" + item.name)
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
            return filebrowser_breadcrumb(this.pathMatch, false);
        }
    },
    mounted: function() {
        this.browse_dir();
    },
    watch: {
        "pathMatch": "browse_dir"
    }
};
</script>

<template>
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
                          <button class="btn btn-sm btn-danger" v-on:click="delete_item(item)" title="Delete">
                              <i class="fa fa-trash"></i>
                          </button>
                      </div>
                  </td>
              </tr>
          </tbody>
      </table>
      <confirmation-prompt
          title="Rename file"
          :show="move_name"
          v-on:confirm="confirm_move()"
          v-on:abort="move_name = null"
          :confirm_disabled="!move_target || move_target == move_name"
          >
          <template v-slot:default>
              Move file <i>{{move_name}}</i>
              <input type="text" class="form-control" placeholder="New name" v-model="move_target"/>
          </template>

          <template v-slot:confirm_button>
              <i class="fa fa-shuffle"></i>
              Confirm rename
          </template>
      </confirmation-prompt>

  </default-layout>
</template>
