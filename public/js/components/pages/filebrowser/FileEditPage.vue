<script>
import DefaultLayout from "../../layouts/DefaultLayout.vue";
import FileEditor from "./FileEditor.vue";
import { filebrowser_breadcrumb } from "./common.js";

export default {
    props: ["pathMatch"],
    components: {
        "default-layout": DefaultLayout,
        "file-editor": FileEditor
    },
    computed: {
        breadcrumb: function() {
            return filebrowser_breadcrumb(this.pathMatch, true);
        }
    }
};
</script>

<template>
  <default-layout icon="edit" title="File-edit" :breadcrumb="breadcrumb">
      <file-editor :filename="pathMatch"/>
  </default-layout>
</template>
