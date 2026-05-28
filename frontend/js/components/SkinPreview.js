import { get_skin } from "../api/playerinfo.js";
import renderskin from "../util/renderskin.js";

export default {
    props: ["playername"],
    data: function() {
        return {
            image_src: null
        };
    },
    mounted: function() {
        get_skin(this.playername)
        .then(b => renderskin(URL.createObjectURL(b)))
        .then(u => this.image_src = u);
    },
    template: /*html*/`
        <img :src="image_src" v-if="image_src" height="32" width="16"/>
        <i class="fas fa-user" v-else></i>
    `
};
