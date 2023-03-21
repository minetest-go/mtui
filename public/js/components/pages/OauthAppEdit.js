import { get_by_id, remove, save } from "../../api/oauth_app.js";

export default {
    data: function() {
        return {
            app: null,
            busy: false
        };
    },
    computed: {
        input_valid: function() {
            return this.app.name != "" && this.app.redirect_urls != "";
        }
    },
    methods: {
        update: function() {
            get_by_id(this.$route.params.id)
            .then(app => this.app = app);
        },
        save: function() {
            this.busy = true;
            save(this.app)
            .then(() => {
                this.busy = false;
                this.$router.push("/oauth-apps");
            });
        },
        remove: function() {
            this.busy = true;
            remove(this.app.id)
            .then(() => {
                this.busy = false;
                this.$router.push("/oauth-apps");
            });
        }
    },
    mounted: function() {
        this.update();
    },
    template: /*html*/`
        <div v-if="app">
            <div class="row">
                <div class="col-md-2">
                </div>
                <div class="col-md-8">
                    <table class="table">
                        <tbody>
                            <tr>
                                <td>ID</td>
                                <td>{{app.id}}</td>
                            </tr>
                            <tr>
                                <td>Name</td>
                                <td>
                                    <input type="text" class="form-control" placeholder="Application name" v-model="app.name"/>
                                </td>
                            </tr>
                            <tr>
                                <td>Enabled</td>
                                <td>
                                    <input type="checkbox" class="form-check-input" v-model="app.enabled"/>
                                </td>
                            </tr>
                            <tr>
                                <td>Created</td>
                                <td>{{new Date(app.created*1000)}}</td>
                            </tr>
                            <tr>
                                <td>Redirect URLs</td>
                                <td>
                                    <input type="text" class="form-control" placeholder="Redirection urls, comma separated" v-model="app.redirect_urls"/>
                                </td>
                            </tr>
                            <tr>
                                <td>Secret</td>
                                <td>{{app.secret}}</td>
                            </tr>
                            <tr>
                                <td>Allowed privs</td>
                                <td>
                                    <input type="text" class="form-control" placeholder="Allowed privs, comma separated" v-model="app.allowed_privs"/>
                                </td>
                            </tr>
                            <tr>
                                <td>Action</td>
                                <td>
                                    <div class="btn-group">
                                        <button class="btn btn-danger" v-on:click="remove">
                                            <i class="fa fa-trash"></i> Delete
                                        </button>
                                        <button class="btn btn-success" :disabled="!input_valid" v-on:click="save">
                                            <i class="fa fa-save"></i> Save
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="col-md-2">
                </div>
            </div>
        </div>
    `
};