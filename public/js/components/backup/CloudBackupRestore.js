import { get_backup_restore_info, create_backup_restore_job } from "../../api/backup-restore.js";
import { get_maintenance } from "../../service/stats.js";

import ConfirmationPrompt from "../ConfirmationPrompt.js";

const store = Vue.reactive({
    endpoint: "",
    key_id: "",
    access_key: "",
    bucket: "",
    filename: "world.zip",
    file_key: "",
    handle: null,
    info: null,
    restore_confirm: null
});

export default {
    data: () => store,
    components: {
        "confirmation-prompt": ConfirmationPrompt
    },
    mounted: function() {
        this.update();
        this.handle = setInterval(() => this.update(), 1000);
    },
    unmounted: function() {
        clearInterval(this.handle);
    },
    methods: {
        update: async function() {
            this.info = await get_backup_restore_info();
        },
        backup: async function() {
            await create_backup_restore_job({
                type: "backup",
                endpoint: this.endpoint,
                key_id: this.key_id,
                access_key: this.access_key,
                bucket: this.bucket,
                filename: this.filename,
                file_key: this.file_key
            });
        },
        restore: function() {
            this.restore_confirm = null;
        }
    },
    computed: {
        get_maintenance,
        inputs_valid: function() {
            return (
                this.endpoint &&
                this.key_id &&
                this.access_key &&
                this.bucket &&
                this.filename
            );
        },
        job_in_progress: function() {
            return (this.info && this.info.state == "running");
        }
    },
    template: /*html*/`
    <div>
        <table class="table table-striped">
            <tbody>
                <tr>
                    <td>
                        Endpoint
                    </td>
                    <td>
                        <input class="form-control" type="text" v-model="endpoint" placeholder="http(s) endpoint" :disabled="job_in_progress"/>
                    </td>
                </tr>
                <tr>
                    <td>
                        Key ID
                    </td>
                    <td>
                        <input class="form-control" type="text" v-model="key_id" placeholder="Key identifier" :disabled="job_in_progress"/>
                    </td>
                </tr>
                <tr>
                    <td>
                        Access key
                    </td>
                    <td>
                        <input class="form-control" type="password" v-model="access_key" placeholder="Secret key" :disabled="job_in_progress"/>
                    </td>
                </tr>
                <tr>
                    <td>
                        Bucket
                    </td>
                    <td>
                        <input class="form-control" type="text" v-model="bucket" placeholder="Name of the bucket" :disabled="job_in_progress"/>
                    </td>
                </tr>
                <tr>
                    <td>
                        Filename
                    </td>
                    <td>
                        <input class="form-control" type="text" v-model="filename" placeholder="Filename to use" :disabled="job_in_progress"/>
                    </td>
                </tr>
                <tr>
                    <td>
                        File key (optional)
                    </td>
                    <td>
                        <input class="form-control" type="text" v-model="file_key" placeholder="File key for encryption/decryption" :disabled="job_in_progress"/>
                    </td>
                </tr>
                <tr>
                    <td></td>
                    <td>
                        <span class="btn-group w-100">
                            <button class="btn btn-primary" :disabled="!inputs_valid || job_in_progress" v-on:click="backup">
                                <i class="fa fa-cloud-arrow-up"></i>
                                Backup
                            </button>
                            <button class="btn btn-danger" :disabled="!get_maintenance || !inputs_valid || job_in_progress" v-on:click="restore_confirm = true">
                                <i class="fa fa-cloud-arrow-down"></i>
                                Restore
                            </button>
                        </span>
                        <confirmation-prompt
                            title="Confirm restore"
                            :show="restore_confirm"
                            v-on:abort="restore_confirm = false"
                            v-on:confirm="restore">
                            <b>Warning: </b> All existing world-data will be overwritten by a backup-restore!
                        </confirmation-prompt>
                    </td>
                </tr>
                <tr v-if="info">
                    <td>Status</td>
                    <td>
                        <span>Type: {{info.type}}</span>
                        <span>State: {{info.state}}</span>
                        <span>Message: {{info.message}}</span>
                        <span>Progress: {{info.progress_percent}}</span>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
    `
};