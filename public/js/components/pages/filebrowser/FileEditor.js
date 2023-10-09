import { download_text, upload } from "../../../api/filebrowser.js";
import { get_mode_name } from "./common.js";

export default {
    props: ["filename"],
    data: function() {
        return {
            text: "",
            cm: null,
            success: false
        };
    },
    mounted: function() {
        download_text(this.filename)
        .then(t => this.text = t)
        .then(() => {
            this.cm = CodeMirror.fromTextArea(this.$refs.textarea, {
                lineNumbers: true,
                viewportMargin: Infinity,
                mode: {
                    name: get_mode_name(this.filename)
                }
            });
        });
    },
    methods: {
        save: function() {
            upload(this.filename, this.cm.getValue())
            .then(() => this.success = true);
        }
    },
    template: /*html*/`
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
    `
};
