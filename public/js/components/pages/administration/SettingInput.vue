<script>
// editor for a single setting value, emits a new setting object on every change
const noise_fields = [
    { key: "offset", label: "Offset", number: true },
    { key: "scale", label: "Scale", number: true },
    { key: "spread_x", label: "Spread X", number: true },
    { key: "spread_y", label: "Spread Y", number: true },
    { key: "spread_z", label: "Spread Z", number: true },
    { key: "seed", label: "Seed", number: false },
    { key: "octaves", label: "Octaves", number: true },
    { key: "persistence", label: "Persistence", number: true },
    { key: "lacunarity", label: "Lacunarity", number: true }
];

export default {
    // is_default: the setting is not configured, text inputs then show the default as placeholder only
    props: ["settingtype", "modelValue", "id", "is_default"],
    emits: ["update:modelValue"],
    data: function() {
        return { noise_fields };
    },
    computed: {
        type: function() {
            return this.settingtype.type;
        },
        value: function() {
            return this.modelValue.value;
        },
        placeholder: function() {
            const d = this.settingtype.default;
            return d && d.value ? d.value : "";
        },
        flags: function() {
            return (this.value || "").split(",").map(f => f.trim()).filter(f => f != "");
        }
    },
    methods: {
        emit_field: function(field, value) {
            this.$emit("update:modelValue", Object.assign({}, this.modelValue, { [field]: value }));
        },
        toggle_flag: function(flag, enabled) {
            const flags = this.flags.filter(f => f != flag);
            if (enabled) {
                flags.push(flag);
            }
            this.emit_field("value", flags.join(","));
        }
    }
};
</script>

<template>
  <input
      v-if="type == 'string'"
      :id="id"
      type="text"
      class="form-control"
      :value="is_default ? '' : value"
      :placeholder="placeholder"
      v-on:input="emit_field('value', $event.target.value)"/>

  <div class="form-check form-switch mb-0" v-else-if="type == 'bool'">
      <input
          :id="id"
          type="checkbox"
          role="switch"
          class="form-check-input"
          :checked="value == 'true'"
          v-on:change="emit_field('value', $event.target.checked ? 'true' : 'false')"/>
      <label class="form-check-label small text-body-secondary" :for="id">{{ value == 'true' ? 'On' : 'Off' }}</label>
  </div>

  <input
      v-else-if="type == 'int' || type == 'float'"
      :id="id"
      type="number"
      class="form-control"
      :step="type == 'float' ? 'any' : 1"
      :min="settingtype.min"
      :max="settingtype.max"
      :value="value"
      :placeholder="placeholder"
      v-on:input="emit_field('value', $event.target.value)"/>

  <select
      v-else-if="type == 'enum'"
      :id="id"
      class="form-select"
      :value="value"
      v-on:change="emit_field('value', $event.target.value)">
      <option v-for="choice in settingtype.choices" :key="choice" :value="choice">{{choice}}</option>
  </select>

  <div v-else-if="type == 'flags'">
      <div class="form-check form-check-inline" v-for="choice in settingtype.choices" :key="choice">
          <input
              :id="id + '-' + choice"
              type="checkbox"
              class="form-check-input"
              :checked="flags.includes(choice)"
              v-on:change="toggle_flag(choice, $event.target.checked)"/>
          <label class="form-check-label" :for="id + '-' + choice">{{choice}}</label>
      </div>
  </div>

  <div class="row g-1" v-else-if="type == 'v3f'">
      <div class="col" v-for="axis in ['x', 'y', 'z']" :key="axis">
          <div class="input-group input-group-sm">
              <span class="input-group-text">{{axis.toUpperCase()}}</span>
              <input
                  type="number"
                  step="any"
                  class="form-control"
                  :value="modelValue[axis]"
                  v-on:input="emit_field(axis, Number($event.target.value))"/>
          </div>
      </div>
  </div>

  <details v-else-if="type == 'noise_params_2d' || type == 'noise_params_3d'">
      <summary>Noise parameters</summary>
      <div class="row g-1 mt-1">
          <div class="col-6 col-md-4" v-for="f in noise_fields" :key="f.key">
              <label class="form-label small mb-0" :for="id + '-' + f.key">{{f.label}}</label>
              <input
                  :id="id + '-' + f.key"
                  :type="f.number ? 'number' : 'text'"
                  step="any"
                  class="form-control form-control-sm"
                  :value="modelValue[f.key]"
                  v-on:input="emit_field(f.key, f.number ? Number($event.target.value) : $event.target.value)"/>
          </div>
      </div>
  </details>

  <input
      v-else
      :id="id"
      type="text"
      class="form-control"
      :value="value"
      v-on:input="emit_field('value', $event.target.value)"/>
</template>
