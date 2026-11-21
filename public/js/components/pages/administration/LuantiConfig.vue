<script>
import {
    store,
    apply_filter,
    update_settings,
    pending_keys,
    discard_all_pending,
    save_all_pending
} from '../../../service/mtconfig.js';
import DefaultLayout from '../../layouts/DefaultLayout.vue';
import { START, ADMINISTRATION } from '../../../util/breadcrumbs.js';
import { confirm } from '../../../store/confirm.js';
import { success } from '../../../store/toast.js';
import SettingRow from "./SettingRow.vue";
import CommonSettings from "./CommonSettings.vue";

function warn_unsaved(e) {
    if (pending_keys().length > 0) {
        e.preventDefault();
        e.returnValue = "";
    }
}

export default {
    components: {
        "setting-row": SettingRow,
        "common-settings": CommonSettings,
        "default-layout": DefaultLayout
    },
    created: function() {
        update_settings();
    },
    mounted: function() {
        window.addEventListener("beforeunload", warn_unsaved);
    },
    unmounted: function() {
        window.removeEventListener("beforeunload", warn_unsaved);
    },
    beforeRouteLeave: async function() {
        const count = pending_keys().length;
        if (count == 0) {
            return true;
        }
        return await confirm({
            title: "Discard unsaved changes?",
            message: `You have ${count} unsaved change(s) that will be lost.`,
            confirm_label: "Discard changes"
        }).then(ok => {
            if (ok) {
                discard_all_pending();
            }
            return ok;
        });
    },
    data: function() {
        return {
            store: store,
            search: "",
            only_configured: true,
            saving: false,
            breadcrumb: [START, ADMINISTRATION, {
                name: "Server settings",
                icon: "sliders",
                link: "/minetest-config"
            }]
        };
    },
    computed: {
        tab: function() {
            return this.$route.query.tab == "all" ? "all" : "common";
        },
        pending: function() {
            return pending_keys();
        }
    },
    methods: {
        set_tab: function(tab) {
            this.$router.replace({ query: tab == "all" ? { tab: "all" } : {} });
        },
        apply_filter: function() {
            apply_filter({
                search: this.search,
                only_configured: this.only_configured
            });
        },
        save: async function() {
            this.saving = true;
            try {
                const count = await save_all_pending();
                success(`Saved ${count} setting(s). Some changes only take effect after a server restart.`);
            } finally {
                this.saving = false;
            }
        },
        discard: async function() {
            const ok = await confirm({
                title: "Discard unsaved changes?",
                message: `${this.pending.length} change(s) will be reverted.`,
                confirm_label: "Discard changes"
            });
            if (ok) {
                discard_all_pending();
            }
        }
    },
    watch: {
        "search": "apply_filter",
        "only_configured": "apply_filter"
    }
};
</script>

<template>
  <default-layout icon="sliders" title="Server settings" :breadcrumb="breadcrumb">
      <ul class="nav nav-tabs mb-3" role="tablist">
          <li class="nav-item" role="presentation">
              <button type="button" role="tab" class="nav-link" :class="{active: tab == 'common'}" :aria-selected="tab == 'common'" v-on:click="set_tab('common')">
                  <i class="fa-solid fa-star"></i>
                  Common
              </button>
          </li>
          <li class="nav-item" role="presentation">
              <button type="button" role="tab" class="nav-link" :class="{active: tab == 'all'}" :aria-selected="tab == 'all'" v-on:click="set_tab('all')">
                  <i class="fa-solid fa-list"></i>
                  All settings
              </button>
          </li>
      </ul>

      <div v-if="!store.ready" class="text-body-secondary">
          <i class="fa-solid fa-spinner fa-spin"></i>
          Loading settings
      </div>

      <div v-else-if="tab == 'common'" class="settings-common">
          <common-settings/>
      </div>

      <div v-else>
          <div class="row g-2 align-items-center mb-2">
              <div class="col-md-6">
                  <input type="search" class="form-control" v-model="search" placeholder="Search settings by name or description" aria-label="Search settings"/>
              </div>
              <div class="col-md-4">
                  <div class="form-check form-switch mb-0">
                      <input id="only-configured" type="checkbox" role="switch" class="form-check-input" v-model="only_configured"/>
                      <label class="form-check-label" for="only-configured">Show only settings that are set</label>
                  </div>
              </div>
              <div class="col-md-2 text-md-end text-body-secondary">
                  {{store.filtered_count}} settings
              </div>
          </div>
          <div class="text-body-secondary" v-if="store.filtered_count == 0">
              No settings found.
              <span v-if="only_configured">Turn off "Show only settings that are set" to see all available settings.</span>
          </div>
          <div v-for="topic in store.filtered_topics" :key="topic">
              <h5 class="mt-3">{{topic}}</h5>
              <div class="table-responsive">
                  <table class="table table-sm align-middle">
                      <thead>
                          <tr>
                              <th style="width: 20%;">Name</th>
                              <th style="width: 5%;">Type</th>
                              <th style="width: 25%;">Description</th>
                              <th style="width: 35%;">Value</th>
                              <th style="width: 15%;"></th>
                          </tr>
                      </thead>
                      <tbody>
                          <tr v-for="setting in store.filtered_settings[topic]" :key="setting.key">
                              <setting-row :setting="setting"/>
                          </tr>
                      </tbody>
                  </table>
              </div>
          </div>
      </div>

      <!-- save bar, shown while there are unsaved changes -->
      <div class="settings-savebar card shadow border-warning" v-if="pending.length > 0">
          <div class="card-body d-flex flex-wrap align-items-center gap-2 py-2">
              <i class="fa-solid fa-pen-to-square text-warning"></i>
              <span class="me-auto">
                  <b>{{pending.length}}</b> unsaved change(s):
                  <span class="text-body-secondary small">{{pending.join(", ")}}</span>
              </span>
              <button type="button" class="btn btn-outline-secondary" v-on:click="discard" :disabled="saving">
                  Discard
              </button>
              <button type="button" class="btn btn-primary" v-on:click="save" :disabled="saving">
                  <i class="fa-solid fa-spinner fa-spin" v-if="saving"></i>
                  <i class="fa fa-floppy-disk" v-else></i>
                  Save {{pending.length}} change(s)
              </button>
          </div>
      </div>
  </default-layout>
</template>

<style scoped>
.settings-common {
    max-width: 60rem;
}

.settings-savebar {
    position: sticky;
    bottom: 1rem;
    z-index: 10;
    margin-top: 1rem;
}
</style>
