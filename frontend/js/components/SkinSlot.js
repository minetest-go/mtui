
import { get, upload, remove } from "../api/skin.js";
import renderskin from "../util/renderskin.js";

export default {
    props: ["skin_id"],
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
            get(this.skin_id).then(b => {
                if (b && b.size > 0){
                    renderskin(URL.createObjectURL(b))
                    .then(u => this.current_image = u);
                }
            });
        },
        upload: function(e) {
            const upload_data = e.target.files[0];
            upload_data.arrayBuffer()
            .then(buf => upload(this.skin_id, buf))
            .then(() => this.update())
            .then(() => this.success = true);
        },
        remove: function(){
            remove(this.skin_id)
            .then(() => this.current_image = null);
        }
    },
    template: /*html*/`
    <div class="card">
        <div class="card-header">
            Skin #{{skin_id}}
            <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
        </div>
        <div class="card-body">
            <div class="alert alert-primary" v-if="!current_image">
                Empty
            </div>
            <img
                style="image-rendering: pixelated; height: 128px; width: 64px;"
                :src="current_image"
                v-if="current_image"
            />
            <hr>
            <form @submit.prevent="upload">
                <input type="file" class="form-control" v-on:change="upload"/>
            </form>
            <a class="btn btn-warning" v-on:click="remove" v-if="current_image">
                <i class="fa fa-trash"></i>
                Remove
            </a>
        </div>
    </div>
    `
};