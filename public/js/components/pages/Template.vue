<script>
import DefaultLayout from "../layouts/DefaultLayout.vue";
import { START } from "../../util/breadcrumbs.js";

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            breadcrumb: [START]
        };
    }
};
</script>

<template>
  <default-layout title="" icon="" :breadcrumb="breadcrumb">
  </default-layout>
</template>
