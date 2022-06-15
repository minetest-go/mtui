
import { get, upload } from "../../api/skin.js";

export default {
    data: function() {
        return {
            busy: false,
            current_image: null,
            upload_data: null,
            success: false,
        };
    },
    mounted: function() {
        this.update();
    },
    methods: {
        update: function() {
            get().then(b => {
                if (b && b.size > 0){
                    this.current_image = URL.createObjectURL(b);
                }
            });
        },
        upload: function() {
            this.upload_data.arrayBuffer()
            .then(buf => upload(buf))
            .then(() => this.update())
            .then(() => this.success = true);
        },
        prepare_upload: function(e){
            this.upload_data = e.target.files[0];
        }
    },
    template: /*html*/`
        <div>
            <h3>
                Skin
                <small class="text-muted">manager</small>
                <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
            </h3>
            <div class="row">
                <div class="col-md-6">
                    <div class="card">
                        <div class="card-header">
                            Current skin
                        </div>
                        <div class="card-body">
                            <div class="alert alert-primary" v-if="!current_image">
                                No skin uploaded yet
                            </div>
                            <img :src="current_image" v-if="current_image"/>
                        </div>
                    </div>
                </div>
                <div class="col-md-6">
                    <div class="card">
                        <div class="card-header">
                            Upload new skin
                        </div>
                        <div class="card-body">
                            <form @submit.prevent="upload">
                                <input type="file" class="form-control" v-on:change="prepare_upload"/>
                                <button class="btn btn-primary" type="submit" :disabled="!upload_data">
                                    Upload
                                    <i class="fa-solid fa-check" v-if="success" style="color: green;"></i>
                                </button>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `
};