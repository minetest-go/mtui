import { get_themes } from "../../api/themes.js";
import { get_config, set_config } from "../../api/config.js";

const store = Vue.reactive({
    themes: null,
    current_theme: ""
});

export default {
    data: () => store,
    created: function() {
        if (!store.themes) {
            get_themes().then(t => store.themes = t);
            get_config("theme").then(t => store.current_theme = t);
        }
    },
    methods: {
        set_config: set_config
    },
    template: /*html*/`
    <div>
        <h3>
            UI
            <small class="text-muted">settings</small>
        </h3>
        <table class="table table-striped" v-if="themes">
            <tr>
                <th>Setting</th>
                <th>Value</th>
                <th>Action</th>
            </tr>
            <tr>
                <td>
                    Theme
                </td>
                <td>
                    <select class="form-control" v-model="current_theme">
                        <option v-for="theme in themes" :value="theme">{{theme}}</option>
                    </select>
                </td>
                <td>
                    <a class="btn btn-success w-100" v-on:click="set_config('theme', current_theme)">
                        <i class="fa fa-save"></i> Save
                    </a>
                </td>
            </tr>
        </table>
    </div>
    `
};