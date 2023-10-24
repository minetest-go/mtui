import { download_text, upload } from "../../../api/filebrowser.js";
import { get_mode_name } from "./common.js";
import CodeEditor from "../../CodeEditor.js";

export default {
    props: ["filename"],
    components: {
        "code-editor": CodeEditor
    },
    data: function() {
        return {
            text: "",
            success: false,
            busy: false,
            mode: get_mode_name(this.filename)
        };
    },
    mounted: function() {
        this.busy = true;
        download_text(this.filename)
        .then(t => this.text = t)
        .finally(() => this.busy = false);
    },
    methods: {
        save: function() {
            this.busy = true;
            upload(this.filename, this.text)
            .then(() => this.success = true)
            .finally(() => this.busy = false);
        }
    },
    template: /*html*/`
    <div class="row">
        <div class="col-2">
            <button class="btn btn-success w-100" v-on:click="save">
                <i class="fa fa-floppy-disk"></i>
                Save
                <i class="fa fa-check" v-if="success"></i>
                <i class="fa fa-spinner fa-spin" v-if="busy"></i>
            </button>
        </div>
    </div>
    <hr>
    <code-editor v-model="text" class="w-100" style="height: 800px;" :mode="mode"/>
    `
};
