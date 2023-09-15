import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            breadcrumb: [START]
        };
    },
    template: /*html*/`
        <default-layout title="" icon="" :breadcrumb="breadcrumb">
        </default-layout>
    `
};