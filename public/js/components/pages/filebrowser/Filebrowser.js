import DefaultLayout from "../../layouts/DefaultLayout.js";
import { browse, get_zip_url, get_download_url } from "../../../api/filebrowser.js";
import format_size from "../../../util/format_size.js";

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            result: null
        };
    },
    methods: {
        format_size: format_size,
        get_zip_url: get_zip_url,
        get_download_url: get_download_url,
        browse_dir: function() {
            const dir = "/" + this.$route.params.pathMatch;
            browse(dir)
            .then(r => this.result = r)
            .then(() => {
                if (this.result.dir == "/") {
                    this.result.dir = "";
                }
            });
        }
    },
    mounted: function() {
        this.browse_dir();
    },
    watch: {
        "$route.params.pathMatch": "browse_dir"
    },
    template: /*html*/`
        <default-layout icon="folder" title="Filebrowser">
            <div class="row">
                <div class="col-4">
                </div>
                <div class="col-4">
                </div>
                <div class="col-4" v-if="result">
                    <a class="btn btn-sm btn-secondary" :href="get_zip_url(result.dir)">
                        <i class="fa fa-download"></i>
                        Download as zip
                    </a>
                </div>
            </div>
            <table class="table" v-if="result">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Size</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="item in result.items">
                        <td>
                            <router-link :to="'/filebrowser' + result.dir + '/' + item.name" v-if="item.is_dir">
                                <i class="fa fa-folder"></i>
                                {{item.name}}
                            </router-link>
                            <span v-if="!item.is_dir">
                                <i class="fa fa-file"></i>
                                {{item.name}}
                            </span>
                        </td>
                        <td>
                            <span v-if="!item.is_dir">
                                {{format_size(item.size)}}
                            </span>
                        </td>
                        <td>
                            <div class="btn-group">
                                <a class="btn btn-sm btn-secondary" :href="get_download_url(result.dir + '/' + item.name)">
                                    <i class="fa fa-download"></i>
                                </a>
                                <button class="btn btn-sm btn-secondary">
                                    <i class="fa fa-edit"></i>
                                </button>
                                <button class="btn btn-sm btn-danger">
                                    <i class="fa fa-trash"></i>
                                </button>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </default-layout>
    `
};