<script>
import { store, apply_filter, update_settings } from '../../../service/mtconfig.js';
import DefaultLayout from '../../layouts/DefaultLayout.vue';
import { START, ADMINISTRATION } from '../../../util/breadcrumbs.js';
import SettingRow from "./SettingRow.vue";

export default {
    components: {
        "setting-row": SettingRow,
        "default-layout": DefaultLayout
    },
    created: function() {
        update_settings();
    },
    data: function() {
        return {
            store: store,
            search: "",
            only_configured: true,
            breadcrumb: [START, ADMINISTRATION, {
                name: "Minetest config",
                icon: "cog",
                link: "/minetest-config"
            }]
        };
    },
    methods: {
        apply_filter: function() {
            apply_filter({
                search: this.search,
                only_configured: this.only_configured
            });
        }
    },
    watch: {
        "search": "apply_filter",
        "only_configured": "apply_filter"
    }
};
</script>

<template>
  <default-layout icon="cog" title="Minetest config" :breadcrumb="breadcrumb">
      <div class="row">
          <div class="col-6">
              <input type="text" class="form-control" v-model="search" placeholder="Search settings"/>
          </div>
          <div class="col-4">
              <input type="checkbox" class="form-check-input" v-model="only_configured"/>
              <label class="form-check-label">Show only configured settings</label>
          </div>
          <div class="col-2">
              Found <span class="badge bg-info">{{store.filtered_count}}</span> settings
          </div>
      </div>
      <div v-for="topic in store.filtered_topics">
          <h4>{{topic}}</h4>
          <table class="table table-striped table-sm">
              <thead>
                  <tr>
                      <th style="width: 20%;">Name</th>
                      <th style="width: 5%;">Type</th>
                      <th style="width: 25%;">Description</th>
                      <th style="width: 35%;">Value</th>
                      <th style="width: 15%;" class="text-end">Actions</th>
                  </tr>
              </thead>
              <tbody>
                  <tr v-for="setting in store.filtered_settings[topic]" :key="setting.key">
                      <setting-row :setting="setting"/>
                  </tr>
              </tbody>
          </table>
      </div>
  </default-layout>
</template>
