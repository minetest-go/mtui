import { START, FILEBROWSER } from "../../Breadcrumb.js";

import DefaultLayout from "../../layouts/DefaultLayout.js";
import FileEditor from "./FileEditor.js";

export default {
    props: ["pathMatch"],
    components: {
        "default-layout": DefaultLayout,
        "file-editor": FileEditor
    },
    computed: {
        breadcrumb: function() {
            const bc = [START, FILEBROWSER];
            const parts = this.pathMatch.split("/");

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
        <file-editor :filename="pathMatch"/>
    </default-layout>
    `
};
