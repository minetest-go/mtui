import { get_all } from "../../api/mtconfig.js";

export default {
    data: function() {
        return {
            config: {}
        };
    },
    mounted: function() {
        get_all().then(c => this.config = c);
    },
    template: /*html*/`
        <div>
            <table class="table table-striped table-condensed">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Description</th>
                        <th>Value</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="(value, key) in config" >
                        <td>{{key}}</td>
                        <td></td>
                        <td>{{value}}</td>
                        <td>
                            <router-link :to="'/minetest-config/' + key" class="btn btn-secondary">
                                <i class="fa fa-edit"></i> Edit
                            </router-link>
                        </td>
                    </tr>
                    <tr>
                        <td></td>
                        <td></td>
                        <td></td>
                        <td>
                            <button class="btn btn-success" :disabled="!new_setting_valid" v-on:click="create">
                                <i class="fa fa-plus"></i> Create
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    `
};