import { store, ordered_settings, topics, count } from '../../service/mtconfig.js';

const SettingRow = {
    props: ["setting", "is_set"],
    data: function() {
        return {
            work_setting: this.is_set ? this.setting.current : this.setting.default
        };
    },
    methods: {
        save: function() {
            console.log(this.setting, this.work_setting);
        },
        reset: function() {

        },
        remove: function() {

        }
    },
    template: /*html*/`
    <td>
        {{setting.key}}
        <i class="fa fa-lg fa-square-check" style="color: green;" title="this setting is configured/set in the minetest.conf" v-if="is_set"></i>
    </td>
    <td>
        <span class="badge bg-info">{{setting.type}}</span>
    </td>
    <td>{{setting.short_description}}</td>
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
    <td>
        <div class="btn-group">
            <a class="btn btn-success" v-on:click="save">
                <i class="fa fa-floppy-disk"></i>
                Save
            </a>
            <a class="btn btn-primary" v-on:click="reset">
                <i class="fa-solid fa-arrow-rotate-left"></i>
                Reset
            </a>
            <a class="btn btn-danger" v-on:click="remove">
                <i class="fa fa-trash"></i>
                Delete
            </a>
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
            only_configured: true,
            ordered_settings: ordered_settings,
            topics: topics,
            count: count,
            store: store
        };
    },
    template: /*html*/`
        <div>
            <div class="row">
                <div class="col-6">
                    <input type="text" class="form-control" v-model="store.search" placeholder="Search settings"/>
                </div>
                <div class="col-4">
                    <input type="checkbox" class="form-check-input" v-model="store.only_configured"/>
                    <label class="form-check-label">Show only configured settings</label>
                </div>
                <div class="col-2">
                    Found <span class="badge bg-info">{{count}}</span> settings
                </div>
            </div>
            <div v-for="topic in topics">
                <h4>{{topic}}</h4>
                <table class="table table-striped table-condensed">
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th>Type</th>
                            <th>Description</th>
                            <th>Value</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="setting in ordered_settings[topic]" v-key="setting.key">
                            <setting-row :setting="setting" :is_set="setting.is_set"/>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    `
};