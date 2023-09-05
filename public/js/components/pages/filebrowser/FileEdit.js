import { START, FILEBROWSER } from "../../Breadcrumb.js";
import DefaultLayout from "../../layouts/DefaultLayout.js";

export default {
    components: {
        "default-layout": DefaultLayout
    },
    computed: {
        breadcrumb: function() {
            const bc = [START, FILEBROWSER];
            const parts = this.$route.params.pathMatch.split("/");

            let path = "";
            parts
            .filter(p => p != "")
            .forEach(p => {
                if (path == "") {
                    path = p;
                } else {
                    path = path + "/" + p;
                }

                bc.push({
                    name: p,
                    icon: "folder-open",
                    link: "/filebrowser/" + path
                });
            });

            const lastitem = bc[bc.length-1];
            lastitem.icon = "file";
            lastitem.link = null;

            return bc;
        }
    },
    template: /*html*/`
    <default-layout icon="edit" title="File-edit" :breadcrumb="breadcrumb">
        <div>
            {{$route.params.pathMatch}}
        </div>
    </default-layout>
    `
};