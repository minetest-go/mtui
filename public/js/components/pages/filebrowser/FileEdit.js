import { START, FILEBROWSER } from "../../Breadcrumb.js";
import DefaultLayout from "../../layouts/DefaultLayout.js";
import { download_text, upload } from "../../../api/filebrowser.js";

export default {
    props: ["pathMatch"],
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            text: "",
            cm: null,
            success: false
        };
    },
    mounted: function() {
        download_text(this.pathMatch)
        .then(t => this.text = t)
        .then(() => {

            const mode = {};
            if (this.pathMatch.match(/.*(lua)$/i)) {
                mode.name = "lua";
            } else if (this.pathMatch.match(/.*(js|json)$/i)) {
                mode.name = "javascript";
            }

            this.cm = CodeMirror.fromTextArea(this.$refs.textarea, {
                lineNumbers: true,
                viewportMargin: Infinity,
                mode: mode
            });
        });
    },
    methods: {
        save: function() {
            upload(this.pathMatch, this.cm.getValue())
            .then(() => this.success = true);
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

            const lastitem = bc[bc.length-1];
            lastitem.icon = "file";
            lastitem.link = null;

            return bc;
        }
    },
    template: /*html*/`
    <default-layout icon="edit" title="File-edit" :breadcrumb="breadcrumb">
        <div class="row">
            <div class="col-2">
                <button class="btn btn-success w-100" v-on:click="save">
                    <i class="fa fa-floppy-disk"></i>
                    Save
                    <i class="fa fa-check" v-if="success"></i>
                </button>
            </div>
        </div>
        <hr>
        <textarea ref="textarea" class="w-100" style="height: 800px;" v-model="text"></textarea>
    </default-layout>
    `
};