import DefaultLayout from "../layouts/DefaultLayout.js";
import HelpPopup from '../HelpPopup.js';

import { START } from "../Breadcrumb.js";
import { has_feature } from "../../service/features.js";
import { is_running } from "../../service/engine.js";
import { get_all as get_all_mods } from "../../service/mods.js";

export default {
    components: {
        "default-layout": DefaultLayout,
        "help-popup": HelpPopup
    },
    data: function() {
        return {
            breadcrumb: [START, {
                name: "Help",
                icon: "circle-question",
                link: "/help"
            }]
        };
    },
    methods: {
        has_feature: has_feature,
        is_running: is_running,
        get_all_mods: get_all_mods
    },
    template: /*html*/`
        <default-layout title="Help" icon="circle-question" :breadcrumb="breadcrumb">
            <div class="row">
                <div class="col-2"></div>
                <div class="col-8">
                    <h4>Introduction</h4>
                    <p>
                        Popups for help are marked with a
                        <help-popup title="Intro">
                            Congrats, you opened a help-dialog :)
                        </help-popup>
                        where available and can be opened with a click.
                    </p>
                    <p>
                        You can come back here at any time with the "Help" button on the
                        <router-link to="/">Start</router-link> page.
                    </p>
                    <h4>First steps</h4>
                    <p>
                        You might want to start with
                        <router-link to="/services/engine">setting up a minetest-engine</router-link>
                        first.
                    </p>
                    <p>
                        Afterwards you might want to add a
                        <router-link to="/mods">game and some mods</router-link>
                        and
                        <router-link to="/minetest-config">configure</router-link> them.
                    </p>
                    <h4>Checklist</h4>
                    <ul>
                        <li v-if="has_feature('docker')">
                            <i class="fa-regular fa-square-check" v-if="is_running"></i>
                            <i class="fa-regular fa-square" v-else></i>
                            Create and start the <router-link to="/services/engine">minetest-engine</router-link>
                        </li>
                        <li v-if="has_feature('modmanagement')">
                            <i class="fa-regular fa-square-check" v-if="get_all_mods().find(m => m.mod_type == 'game')"></i>
                            <i class="fa-regular fa-square" v-else></i>
                            Add a <router-link to="/mods">game</router-link>
                        </li>
                        <li v-if="has_feature('modmanagement')">
                            <i class="fa-regular fa-square-check" v-if="get_all_mods().find(m => m.mod_type == 'mod')"></i>
                            <i class="fa-regular fa-square" v-else></i>
                            Add some <router-link to="/mods">mods</router-link>
                        </li>
                    </ul>
                </div>
                <div class="col-2"></div>
            </div>
        </default-layout>
    `
};