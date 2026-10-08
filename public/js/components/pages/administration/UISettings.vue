<script>
import { reactive } from "vue";
import DefaultLayout from "../../layouts/DefaultLayout.vue";
import { START, ADMINISTRATION } from "../../../util/breadcrumbs.js";

const store = reactive({
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
    data: () => store
};
</script>

<template>
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
</template>
