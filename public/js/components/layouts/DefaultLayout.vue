<script>
import Breadcrumb from "../Breadcrumb.vue";

export default {
    props: ["title", "icon", "breadcrumb"],
    components: {
        "bread-crumb": Breadcrumb
    }
};
</script>

<template>
  <bread-crumb :items="breadcrumb" v-if="breadcrumb"/>
  <slot></slot>
</template>
