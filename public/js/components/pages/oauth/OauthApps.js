import { get_all, save } from "../../../api/oauth_app.js";
import { START, ADMINISTRATION, OAUTH_APPS } from "../../Breadcrumb.js";
import DefaultLayout from "../../layouts/DefaultLayout.js";

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            apps: [],
            new_app: {
                name: "",
                domain: ""
            },
            breadcrumb: [START, ADMINISTRATION, OAUTH_APPS]
        };
    },
    computed: {
        new_app_valid: function() {
            return this.new_app.name != "" && this.new_app.domain != "";
        }
    },
    methods: {
        update: function() {
            get_all().then(apps => this.apps = apps);
        },
        create: function() {
            save(this.new_app).then(app => {
                this.new_app.name = "";
                this.new_app.domain = "";
                this.$router.push("/oauth-apps/" + app.id);
            });
        }
    },
    mounted: function() {
        this.update();
    },
    template: /*html*/`
        <default-layout icon="passport" title="OAuth apps" :breadcrumb="breadcrumb">
            <table class="table table-condensed table-striped">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Name</th>
                        <th>Enabled</th>
                        <th>Created</th>
                        <th>Redirect URLs</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="app in apps">
                        <td>
                            <router-link :to="'/oauth-apps/' + app.id">
                                {{app.id}}
                            </router-link>
                        </td>
                        <td>{{app.name}}</td>
                        <td>{{app.enabled}}</td>
                        <td>{{new Date(app.created*1000)}}</td>
                        <td>{{app.domain}}</td>
                        <td>
                            <router-link :to="'/oauth-apps/' + app.id" class="btn btn-secondary">
                                <i class="fa fa-edit"></i> Edit
                            </router-link>
                        </td>
                    </tr>
                    <tr>
                        <td></td>
                        <td>
                            <input type="text" class="form-control" placeholder="Application name" v-model="new_app.name"/>
                        </td>
                        <td></td>
                        <td></td>
                        <td>
                            <input type="text" class="form-control" placeholder="Redirection url" v-model="new_app.domain"/>
                        </td>
                        <td>
                            <button class="btn btn-success" :disabled="!new_app_valid" v-on:click="create">
                                <i class="fa fa-plus"></i> Create
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </default-layout>
    `
};