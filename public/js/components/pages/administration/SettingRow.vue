<script>
import {
    unset,
    is_set,
    is_pending,
    working_setting,
    set_pending,
    discard_pending
} from '../../../service/mtconfig.js';
import { confirm } from '../../../store/confirm.js';
import { success } from '../../../store/toast.js';
import SettingInput from './SettingInput.vue';

export default {
    props: ["setting"],
    components: {
        "setting-input": SettingInput
    },
    data: function() {
        return {
            busy: false
        };
    },
    computed: {
        key: function() {
            return this.setting.key;
        },
        is_set: function() {
            return is_set(this.key);
        },
        is_pending: function() {
            return is_pending(this.key);
        },
        working: function() {
            return working_setting(this.key);
        }
    },
    methods: {
        update: function(v) {
            set_pending(this.key, v);
        },
        undo: function() {
            discard_pending(this.key);
        },
        unset: async function() {
            const ok = await confirm({
                title: `Unset ${this.key}?`,
                message: "The setting is removed from minetest.conf and falls back to its default value.",
                confirm_label: "Unset",
                confirm_icon: "fa fa-trash"
            });
            if (!ok) {
                return;
            }
            this.busy = true;
            unset(this.key)
            .then(() => success(`${this.key} reset to its default value`))
            .finally(() => this.busy = false);
        }
    }
};
</script>

<template>
  <td :class="{'table-warning': is_pending}">
      <label :for="'setting-' + key" class="text-break">{{key}}</label>
      <span class="badge text-bg-success ms-1" v-if="is_set" title="Set in minetest.conf">set</span>
      <span class="badge text-bg-warning ms-1" v-if="is_pending">unsaved</span>
  </td>
  <td :class="{'table-warning': is_pending}">
      <span class="badge text-bg-info">{{setting.type}}</span>
  </td>
  <td :class="{'table-warning': is_pending}">
      <details>
          <summary>{{setting.short_description}}</summary>
          {{setting.long_description}}
      </details>
  </td>
  <td :class="{'table-warning': is_pending}">
      <setting-input
          :id="'setting-' + key"
          :settingtype="setting"
          :model-value="working"
          :is_default="!is_set && !is_pending"
          v-on:update:model-value="update"/>
  </td>
  <td class="text-end text-nowrap" :class="{'table-warning': is_pending}">
      <i class="fa fa-spinner fa-spin me-1" v-if="busy"></i>
      <button class="btn btn-sm btn-outline-secondary me-1" v-on:click="undo" v-if="is_pending" title="Discard the unsaved change">
          <i class="fa-solid fa-arrow-rotate-left"></i>
          Undo
      </button>
      <button class="btn btn-sm btn-outline-danger" v-on:click="unset" :disabled="!is_set || busy" title="Remove from minetest.conf, use the default">
          <i class="fa fa-trash"></i>
          Unset
      </button>
  </td>
</template>
