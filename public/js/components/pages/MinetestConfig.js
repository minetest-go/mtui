import { store, ordered_settings, topics, count } from '../../service/mtconfig.js';

const SettingRow = {
    props: ["setting"],
    template: /*html*/`
    <td>{{setting.key}}</td>
    <td>
        <span class="badge bg-info">{{setting.type}}</span>
    </td>
    <td>{{setting.short_description}}</td>
    <td>
        <div v-if="setting.type == 'string'">
            <input type="text" class="form-control" :value="setting.default.value"/>
        </div>
        <div v-if="setting.type == 'bool'">
            <input type="checkbox" class="form-check-input" :checked="setting.default.value == 'true'"/>
        </div>
        <div v-if="setting.type == 'int' || setting.type == 'float'">
            <input type="number" class="form-control" :value="setting.default.value" :min="setting.min" :max="setting.max"/>
        </div>
        <div v-if="setting.type == 'enum'">
            <select class="form-control">
                <option v-for="choice in setting.choices" :selected="choice == setting.default.value">{{choice}}</option>
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
    </td>
    <td>
        <div class="btn-group" v-if="setting.type == 'string'">
            <a class="btn btn-success">
                <i class="fa fa-floppy-disk"></i>
                Save
            </a>
            <a class="btn btn-primary">
                <i class="fa-solid fa-arrow-rotate-left"></i>
                Reset
            </a>
            <a class="btn btn-danger">
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
                            <setting-row :setting="setting"/>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    `
};