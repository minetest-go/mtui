<script>
import {
    store,
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
import sections from './common_settings.js';

export default {
    components: {
        "setting-input": SettingInput
    },
    computed: {
        sections: function() {
            return sections.map(section => Object.assign({}, section, {
                settings: section.settings
                    .filter(s => store.settingtypes[s.key])
                    .filter(s => !s.requires || working_setting(s.requires).value == "true")
                    .map(s => Object.assign({ settingtype: store.settingtypes[s.key] }, s))
            }));
        }
    },
    methods: {
        is_set,
        is_pending,
        working_setting,
        set_pending,
        discard_pending,
        reset: async function(s) {
            const ok = await confirm({
                title: `Reset "${s.label}"?`,
                message: "The setting goes back to the default value of the game engine.",
                confirm_label: "Reset to default",
                confirm_icon: "fa-solid fa-arrow-rotate-left"
            });
            if (!ok) {
                return;
            }
            await unset(s.key);
            success(`"${s.label}" reset to its default value`);
        }
    }
};
</script>

<template>
  <div class="card mb-3" v-for="section in sections" :key="section.title">
      <div class="card-header">
          <div class="fw-semibold">
              <i :class="section.icon" class="fa-fw me-1"></i>
              {{section.title}}
          </div>
          <div class="small text-body-secondary">{{section.description}}</div>
      </div>
      <ul class="list-group list-group-flush">
          <li
              class="list-group-item"
              :class="{'list-group-item-warning': is_pending(s.key)}"
              v-for="s in section.settings"
              :key="s.key">
              <div class="row g-2 align-items-center">
                  <div class="col-md-5">
                      <label class="fw-medium" :for="'common-' + s.key">{{s.label}}</label>
                      <span class="badge text-bg-warning ms-1" v-if="is_pending(s.key)">unsaved</span>
                      <div class="small text-body-secondary">{{s.help}}</div>
                  </div>
                  <div class="col-md-5">
                      <setting-input
                          :id="'common-' + s.key"
                          :settingtype="s.settingtype"
                          :model-value="working_setting(s.key)"
                          :is_default="!is_set(s.key) && !is_pending(s.key)"
                          v-on:update:model-value="v => set_pending(s.key, v)"/>
                  </div>
                  <div class="col-md-2 text-md-end small text-nowrap">
                      <button
                          type="button"
                          class="btn btn-link btn-sm p-0"
                          v-if="is_pending(s.key)"
                          v-on:click="discard_pending(s.key)">
                          Undo
                      </button>
                      <button
                          type="button"
                          class="btn btn-link btn-sm p-0 link-secondary"
                          v-else-if="is_set(s.key)"
                          v-on:click="reset(s)"
                          title="Remove from minetest.conf and use the default">
                          Reset to default
                      </button>
                      <span class="text-body-secondary" v-else>Default</span>
                  </div>
              </div>
          </li>
      </ul>
  </div>
</template>
