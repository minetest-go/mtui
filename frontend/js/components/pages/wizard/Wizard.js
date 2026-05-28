import DefaultLayout from "../../layouts/DefaultLayout.js";
import Engine from "./Engine.js";
import Game from "./Game.js";
import Mods from "./Mods.js";
import Settings from "./Settings.js";
import Done from "./Done.js";

import { START } from "../../Breadcrumb.js";

export default {
    props: ["step"],
    components: {
        "default-layout": DefaultLayout,
        "engine-step": Engine,
        "game-step": Game,
        "mods-step": Mods,
        "settings-step": Settings,
        "done-step": Done
    },
    data: function() {
        return {
            breadcrumb: [START, {
                name: "Setup wizard", icon: "wand-magic-sparkles", link: `/wizard/${this.step}`
            }],
            max_steps: 5
        };
    },
    methods: {
        next: function() {
            if (this.step < this.max_steps) {
                this.$router.push(`/wizard/${+this.step+1}`);
            }
        },
        previous: function() {
            if (this.step > 1) {
                this.$router.push(`/wizard/${+this.step-1}`);
            }
        }
    },
    template: /*html*/`
        <default-layout title="Setup wizard" icon="wand-magic-sparkles" :breadcrumb="breadcrumb">
            <div class="row">
                <div class="col-5">
                    <button class="btn btn-outline-secondary w-100" v-on:click="previous" v-if="step > 1">
                        <i class="fa-solid fa-chevron-left"></i>
                        Previous
                    </button>
                </div>
                <div class="col-2 text-center">
                    Step
                    <span class="badge bg-primary">{{step}}</span>
                    /
                    <span class="badge bg-primary">{{max_steps}}</span>
                </div>
                <div class="col-5">
                    <button class="btn btn-outline-secondary w-100" v-on:click="next" v-if="step < max_steps">
                        Next
                        <i class="fa-solid fa-chevron-right"></i>
                    </button>
                </div>
            </div>
            <div class="row">
                <div class="col-2"></div>
                <div class="col-8">
                    &nbsp;
                    <engine-step v-if="step == 1"/>
                    <game-step v-if="step == 2"/>
                    <mods-step v-if="step == 3"/>
                    <settings-step v-if="step == 4"/>
                    <done-step v-if="step == 5"/>
                </div>
                <div class="col-2"></div>
            </div>
        </default-layout>
    `
};