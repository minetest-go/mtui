<script>
import { save, unset } from '../../../service/mtconfig.js';

export default {
    props: ["setting"],
    data: function() {
        return {
            old_setting: this.setting.current ? this.setting.current : this.setting.default,
            work_setting: Object.assign({}, this.setting.current ? this.setting.current : this.setting.default),
            busy: false,
            is_set: this.setting.is_set
        };
    },
    computed: {
        is_changed: function() {
            const w = this.work_setting;
            const o = this.old_setting;
            switch (this.setting.type) {
                case "string":
                case "int":
                case "float":
                case "bool":
                case "enum":
                    return w.value != o.value;
                case "v3f":
                    return w.x != o.x || w.y != o.y || w.z != o.z;
                // TODO: flags, noise_params_2d, noise_params_3d
                }
        }
    },
    methods: {
        save: function() {
            this.busy = true;
            // ensure that the value field is of type "string"
            if (this.work_setting.value) {
                this.work_setting.value = "" + this.work_setting.value;
            }
            save(this.setting.key, this.work_setting)
            .then(() => {
                this.is_set = true;
                Object.assign(this.old_setting, this.work_setting);
                this.busy = false;
            });
        },
        reset: function() {
            Object.assign(this.work_setting, this.old_setting);
        },
        unset: function() {
            this.busy = true;
            unset(this.setting.key)
            .then(() => {
                this.is_set = false;
                Object.apply(this.work_setting, this.setting.default);
                this.busy = false;
            });
        }
    }
};
</script>

<template>
  <td>
      {{setting.key}}
      <i class="fa fa-lg fa-square-check" style="color: green;" title="this setting is configured/set in the minetest.conf" v-if="is_set"></i>
  </td>
  <td>
      <span class="badge bg-info">{{setting.type}}</span>
  </td>
  <td>
      <details>
          <summary>{{setting.short_description}}</summary>
          {{setting.long_description}}
      </details>
  </td>
  <td>
      <div v-if="setting.type == 'string'">
          <input type="text" class="form-control" v-model="work_setting.value"/>
      </div>
      <div v-if="setting.type == 'bool'">
          <input type="checkbox" class="form-check-input" v-model="work_setting.value" true-value="true" false-value="false"/>
      </div>
      <div v-if="setting.type == 'int' || setting.type == 'float'">
          <input type="number" class="form-control" v-model="work_setting.value" :min="setting.min" :max="setting.max"/>
      </div>
      <div v-if="setting.type == 'enum'">
          <select class="form-control" v-model="work_setting.value">
              <option v-for="choice in setting.choices">{{choice}}</option>
          </select>
      </div>
      <div v-if="setting.type == 'flags'">
          <details>
              <summary>Flags setting</summary>
              <ul>
                  <li v-for="choice in setting.choices">
                      <input type="checkbox" class="form-check-input"/>
                      {{choice}}
                  </li>
              </ul>
          </details>
      </div>
      <div v-if="setting.type == 'v3f'">
          <details>
              <summary>3D Vector setting</summary>
              <label>X</label>
              <input type="text" class="form-control" v-model.number="work_setting.x"/>
              <label>Y</label>
              <input type="text" class="form-control" v-model.number="work_setting.y"/>
              <label>Z</label>
              <input type="text" class="form-control" v-model.number="work_setting.z"/>
          </details>
      </div>
      <div v-if="setting.type == 'noise_params_2d' || setting.type == 'noise_params_3d'">
          <details>
              <summary>Noise parameter setting</summary>
              <label>Offset</label>
              <input type="number" class="form-control" v-model.number="work_setting.offset"/>
              <label>Scale</label>
              <input type="number" class="form-control" v-model.number="work_setting.scale"/>
              <label>Spread X</label>
              <input type="number" class="form-control" v-model.number="work_setting.spread_x"/>
              <label>Spread Y</label>
              <input type="number" class="form-control" v-model.number="work_setting.spread_y"/>
              <label>Spread Z</label>
              <input type="number" class="form-control" v-model.number="work_setting.spread_z"/>
              <label>Seed</label>
              <input type="text" class="form-control" v-model="work_setting.seed"/>
              <label>Octaves</label>
              <input type="number" class="form-control" v-model.number="work_setting.octaves"/>
              <label>Persistence</label>
              <input type="number" class="form-control" v-model.number="work_setting.persistence"/>
              <label>Lacunarity</label>
              <input type="number" class="form-control" v-model.number="work_setting.lacunarity"/>
          </details>
      </div>
  </td>
  <td class="text-end">
      <i class="fa fa-spinner fa-spin" v-if="busy"></i>
      <div class="btn-group">
          <button class="btn btn-success" v-on:click="save" :disabled="!is_changed">
              <i class="fa fa-floppy-disk"></i>
              Save
          </button>
          <button class="btn btn-primary" v-on:click="reset" :disabled="!is_changed">
              <i class="fa-solid fa-arrow-rotate-left"></i>
              Reset
          </button>
          <button class="btn btn-danger" v-on:click="unset" :disabled="!is_set">
              <i class="fa fa-trash"></i>
              Unset
          </button>
      </div>
  </td>
</template>
