import DefaultLayout from "../../layouts/DefaultLayout.js";
import { START, ADMINISTRATION } from "../../Breadcrumb.js";

const store = Vue.reactive({
    breadcrumb: [START, ADMINISTRATION, {
        name: "UI Settings",
        icon: "list-check",
        link: "/ui/settings"
    }]
});

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: () => store,
    template: /*html*/`
    <default-layout icon="list-check" title="UI Settings" :breadcrumb="breadcrumb">
        <table class="table table-striped">
            <thead>
                <tr>
                    <th>Setting</th>
                    <th>Value</th>
                    <th>Action</th>
                </tr>
            </thead>
            <tbody>
            </tbody>
        </table>
    </default-layout>
    `
};