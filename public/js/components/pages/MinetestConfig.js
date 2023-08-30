import { store, apply_filter } from '../../service/mtconfig.js';

const SettingRow = {
    props: ["setting"],
    data: function() {
        return {
            work_setting: this.setting.current ? this.setting.current : this.setting.default
        };
    },
    computed: {
        is_changed: function() {
            return this.work_setting.value != this.setting.default.value;
        }
    },
    methods: {
        save: function() {
        },
        reset: function() {
        },
        unset: function() {
        }
    },
    template: /*html*/`
    <td>
        {{setting.key}}
        <i class="fa fa-lg fa-square-check" style="color: green;" title="this setting is configured/set in the minetest.conf" v-if="setting.is_set"></i>
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
            <input type="text" class="form-control" :value="work_setting.value"/>
        </div>
        <div v-if="setting.type == 'bool'">
            <input type="checkbox" class="form-check-input" :checked="work_setting.value == 'true'"/>
        </div>
        <div v-if="setting.type == 'int' || setting.type == 'float'">
            <input type="number" class="form-control" :value="work_setting.value" :min="setting.min" :max="setting.max"/>
        </div>
        <div v-if="setting.type == 'enum'">
            <select class="form-control">
                <option v-for="choice in setting.choices" :selected="choice == work_setting.value">{{choice}}</option>
            </select>
        </div>
        <div v-if="setting.type == 'flags'">
            <ul>
                <li v-for="choice in setting.choices">
                    <input type="checkbox" class="form-check-input"/>
                    {{choice}}
                </li>
            </ul>
        </div>
        <div v-if="setting.type == 'v3f'">
            <details>
                <summary>3D Vector setting</summary>
                <label>X</label>
                <input type="text" class="form-control" :value="work_setting.x"/>
                <label>Y</label>
                <input type="text" class="form-control" :value="work_setting.y"/>
                <label>Z</label>
                <input type="text" class="form-control" :value="work_setting.z"/>
            </details>
        </div>
        <div v-if="setting.type == 'noise_params_2d' || setting.type == 'noise_params_3d'">
            <details>
                <summary>Noise parameter setting</summary>
                <label>Offset</label>
                <input type="number" class="form-control" :value="work_setting.offset"/>
                <label>Scale</label>
                <input type="number" class="form-control" :value="work_setting.scale"/>
                <label>Spread X</label>
                <input type="number" class="form-control" :value="work_setting.spread_x"/>
                <label>Spread Y</label>
                <input type="number" class="form-control" :value="work_setting.spread_y"/>
                <label>Spread Z</label>
                <input type="number" class="form-control" :value="work_setting.spread_z"/>
                <label>Seed</label>
                <input type="text" class="form-control" :value="work_setting.seed"/>
                <label>Octaves</label>
                <input type="number" class="form-control" :value="work_setting.octaves"/>
                <label>Persistence</label>
                <input type="number" class="form-control" :value="work_setting.persistence"/>
                <label>Lacunarity</label>
                <input type="number" class="form-control" :value="work_setting.lacunarity"/>
            </details>
        </div>
    </td>
    <td class="text-end">
        <div class="btn-group">
            <button class="btn btn-success" v-on:click="save" :disabled="!is_changed">
                <i class="fa fa-floppy-disk"></i>
                Save
            </button>
            <button class="btn btn-primary" v-on:click="reset" :disabled="!is_changed">
                <i class="fa-solid fa-arrow-rotate-left"></i>
                Reset
            </button>
            <button class="btn btn-danger" v-on:click="unset" :disabled="!setting.is_set">
                <i class="fa fa-trash"></i>
                Unset
            </button>
        </div>
    </td>
    `
};

export default {
    components: {
        "setting-row": SettingRow
    },
    data: function() {
        return {
            store: store,
            search: "",
            only_configured: true
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
    },
    template: /*html*/`
        <div>
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
        </div>
    `
};