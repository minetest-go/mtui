import store from '../../store/mtconfig.js';
import '../../service/mtconfig.js';

export default {
    data: function() {
        return {
            config: {},
            only_configured: true,
            store: store
        };
    },
    template: /*html*/`
        <div>
            <div class="row">
                <div class="col-6">
                    <input type="text" class="form-control" placeholder="Keywords"/>
                </div>
                <div class="col-4">
                    <input type="checkbox" class="form-check-input" v-model="only_configured"/>
                    <label class="form-check-label">Show only configured settings</label>
                </div>
                <div class="col-2">
                    <a class="btn btn-primary w-100">
                        Search
                    </a>
                </div>
            </div>
            <div v-for="topic in store.topics">
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
                        <tr v-for="setting in store.ordered_settings[topic]" >
                            <td>{{setting.key}}</td>
                            <td>
                                <span class="badge bg-info">{{setting.type}}</span>
                            </td>
                            <td>{{setting.short_description}}</td>
                            <td>
                                <div v-if="setting.type == 'string'">
                                    <input type="text" class="form-control"/>
                                </div>
                                <div v-if="setting.type == 'bool'">
                                    <input type="checkbox" class="form-check-input" :checked="setting.default == 'true'"/>
                                </div>
                                <div v-if="setting.type == 'int' || setting.type == 'float'">
                                    <input type="number" class="form-control" :value="setting.default" :min="setting.min" :max="setting.max"/>
                                </div>
                                <div v-if="setting.type == 'enum'">
                                    <select class="form-control">
                                        <option v-for="choice in setting.choices" :selected="choice == setting.default">{{choice}}</option>
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
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    `
};