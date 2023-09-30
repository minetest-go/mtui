import DefaultLayout from "../layouts/DefaultLayout.js";
import HelpPopup from '../HelpPopup.js';

import { START } from "../Breadcrumb.js";
import { has_feature } from "../../service/features.js";

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
        has_feature
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
                        You might want to start with the
                        <router-link to="/wizard/1" class="btn btn-primary" v-if="has_feature('docker') && has_feature('minetest_config') && has_feature('modmanagement')">
                            <i class="fa-solid fa-wand-magic-sparkles"></i> Setup wizard
                        </router-link>
                    </p>

                </div>
                <div class="col-2"></div>
            </div>
        </default-layout>
    `
};