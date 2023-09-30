import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            breadcrumb: [START, {
                name: "Setup wizard", icon: "wand-magic-sparkles", link: "/wizard"
            }]
        };
    },
    template: /*html*/`
        <default-layout title="Setup wizard" icon="wand-magic-sparkles" :breadcrumb="breadcrumb">
        </default-layout>
    `
};