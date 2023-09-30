import { get_setting, save, is_ready } from "../../../service/mtconfig.js";

const SettingSuggestion = {
    props: ["type", "name", "title", "description", "default_value"],
    data: function() {
        const setting = get_setting(this.name);
        return {
            busy: false,
            setting: setting,
            value: setting ? setting.value : this.default_value
        };
    },
    methods: {
        set: function(value) {
            this.busy = true;
            save(this.name, { value: ""+value })
            .then(() => this.busy = false);
        },
    },
    template: /*html*/`
    <div class="card w-100" style="padding: 10px;">
        <div class="row">
            <div class="col-3">
                <h4>{{title}}</h4>
            </div>
            <div class="col-5">
                {{description}}
            </div>
            <div class="col-4" v-if="type == 'bool'">
                <button class="btn btn-success w-100" v-if="setting && setting.value == 'true'" v-on:click="set('false')" :disabled="busy">
                    <i class="fa fa-check"></i>
                    Enabled
                    <i class="fa fa-spinner fa-spin" v-if="busy"></i>
                </button>
                <button class="btn btn-secondary w-100" v-else v-on:click="set('true')" :disabled="busy">
                    <i class="fa fa-close"></i>
                    Disabled
                    <i class="fa fa-spinner fa-spin" v-if="busy"></i>
                </button>
            </div>
            <div class="col-4" v-if="type == 'int'">
                <div class="input-group w-100">
                    <input class="form-control" type="number" v-model="value"/>
                    <button class="btn btn-success" v-on:click="set(this.value)" :disabled="busy">
                        <i class="fa fa-floppy-disk"></i>
                        Save
                        <i class="fa fa-spinner fa-spin" v-if="busy"></i>
                    </button>
                </div>
            </div>
            <div class="col-4" v-if="type == 'string'">
                <div class="input-group w-100">
                    <input class="form-control" type="text" v-model="value"/>
                    <button class="btn btn-success" v-on:click="set(this.value)" :disabled="busy">
                        <i class="fa fa-floppy-disk"></i>
                        Save
                        <i class="fa fa-spinner fa-spin" v-if="busy"></i>
                    </button>
                </div>
            </div>
        </div>
    </div>
    `
};

export default {
    components: {
        "setting-suggestion": SettingSuggestion
    },
    computed: {
        ready: is_ready
    },
    template: /*html*/`
    <div v-if="ready">
        <h4>Common settings</h4>
        This is just a subset of the available settings, see the <router-link to="/minetest-config">minetest config</router-link> page
        <setting-suggestion name="server_name" type="string" title="Server name" description="The name of the server"/>
        <setting-suggestion name="max_users" type="int" title="Max players" description="Maximum players allowed on the server" :default_value="15"/>
        <setting-suggestion name="server_announce" type="bool" title="Announce" description="Announce the server in the public serverlist"/>
    </div>
    `
};