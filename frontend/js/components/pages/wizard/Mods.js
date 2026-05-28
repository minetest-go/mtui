import { get_game, get_cdb_mod, add } from "../../../service/mods.js";
import { set_feature } from "../../../service/features.js";

import CDBPackageLink from "../../CDBPackageLink.js";

const ModSuggestion = {
    props: ["name", "author", "description", "enables_feature"],
    components: {
        "cdb-link": CDBPackageLink
    },
    data: function() {
        return {
            busy: false
        };
    },
    methods: {
        add: function() {
            this.busy = true;
            add({
                author: this.author,
                name: this.name,
                mod_type: "mod",
                source_type: "cdb"
            })
            .then(() => {
                this.busy = false;
                if (this.enables_feature) {
                    set_feature(this.enables_feature, true);
                }
            });
        }
    },
    computed: {
        added: function() {
            return get_cdb_mod(this.author, this.name);
        }
    },
    template: /*html*/`
    <div class="card w-100" style="padding: 10px;">
        <div class="row">
            <div class="col-3">
                <h4>
                    {{name}}
                    <small class="text-muted">by {{author}}</small>
                </h4>
            </div>
            <div class="col-3">
                <cdb-link :name="name" :author="author"/>
            </div>
            <div class="col-4">
                {{description}}
                <span class="badge bg-success" v-if="enables_feature">
                    UI-Integration
                </span>
            </div>
            <div class="col-2">
                <button class="btn btn-success w-100" v-if="!added" v-on:click="add">
                    <i class="fa fa-plus"></i>
                    Add
                    <i class="fa fa-spinner fa-spin" v-if="busy"></i>
                </button>
                <button class="btn btn-secondary w-100" disabled="true" v-if="added">
                    <i class="fa fa-check"></i>
                    Added
                </button>
            </div>
        </div>
    </div>
    `
};


export default {
    components: {
        "mod-suggestion": ModSuggestion,
    },
    data: function() {
        return {
            
        };
    },
    computed: {
        game: get_game
    },
    template: /*html*/`
    <div>
        <h4>Suggested mods</h4>
        General suggestions
        (more mods and other installation methods are available in the <router-link to="/mods">mods</router-link> page)
        <mod-suggestion name="xban2" author="Kaeza" description="Enhanced player ban system" enables_feature="xban"/>
        <mod-suggestion name="areas" author="ShadowNinja" description="Area protection system"/>
        <mod-suggestion name="mail" author="mt-mods" description="Ingame mail system" enables_feature="mail"/>

        Game specific suggestions
        <div v-if="game && game.name == 'minetest_game'">
            <mod-suggestion name="skinsdb" author="bell07" description="Skin database" enables_feature="skinsdb"/>
            <mod-suggestion name="mesecons" author="Jeija" description="Mesecons digital circuitry" enables_feature="mesecons"/>
        </div>
        <div v-if="game && game.name == 'nodecore'">
            <mod-suggestion name="nc_cats" author="Warr1024" description="Nodecore cats!"/>
        </div>
    </div>
    `
};