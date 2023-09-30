import { get_mods_by_type, add } from "../../../service/mods.js";
import { search_packages } from "../../../api/cdb.js";

import Preview from "../cdb/Preview.js";

export default {
    components: {
        "package-preview": Preview
    },
    data: function() {
        return {
            packages: [],
            busy: false,
            busy_pkg: null
        };
    },
    created: function() {
        search_packages({
            type: ["game"],
            query: "",
            limit: 20,
            sort: "score",
            order: "desc"
        })
        .then(pkgs => this.packages = pkgs);
    },
    methods: {
        install: function(pkg) {
            this.busy = true;
            this.busy_pkg = pkg;
            add({
                author: this.busy_pkg.author,
                name: this.busy_pkg.name,
                mod_type: "game",
                source_type: "cdb"
            })
            .then(() => this.busy = false);
        }
    },
    computed: {
        game: function() {
            const list = get_mods_by_type("game");
            return list.length == 1 ? list[0] : null;
        },
        complete: function() {
            return this.game;
        }
    },
    template: /*html*/`
    <div>
        <div class="row">
            <div class="col-12">
                <div class="alert alert-primary" v-if="busy">
                    <i class="fa fa-spinner fa-spin"></i> Installing: <b>{{busy_pkg.name}}</b>
                </div>
                <div class="alert alert-success" v-if="complete">
                    <i class="fa fa-check"></i> Game installed: <b>{{game.name}}</b>
                </div>
            </div>
        </div>
        <div v-if="!complete && !busy">
            Select the game to install
            <div style="display: flex; flex-wrap: wrap;">
            <package-preview
                v-for="pkg in packages"
                :pkg="pkg"
                :install_button="true"
                v-on:install="install(pkg)">
            </package-preview>
        </div>
        </div>
    </div>
    `
};