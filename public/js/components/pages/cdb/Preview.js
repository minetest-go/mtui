
export default {
    props: ["pkg", "install_button"],
    computed: {
        thumbnail: function() {
            return this.pkg.thumbnail.replaceAll("/thumbnails/1", "/thumbnails/2");
        }
    },
    template: /*html*/`
    <div class="card" style="width: 16rem; margin: 5px; flex: 0 0 auto;">
        <router-link :to="'/cdb/detail/' + pkg.author + '/' + pkg.name">
            <img :src="thumbnail" class="card-img-top" style="height: 180px;"/>
        </router-link>
        <div class="card-body">
            <h5 class="card-title">
                <router-link :to="'/cdb/detail/' + pkg.author + '/' + pkg.name">
                    {{pkg.name}}
                </router-link>
            </h5>
            <h6 class="card-subtitle mb-2 text-body-secondary">by {{pkg.author}}</h6>
            <p class="card-text">{{pkg.short_description}}</p>
            <button class="btn btn-success"
                style="position: absolute; bottom: 15px; right: 15px;"
                v-if="install_button"
                v-on:click="$emit('install', pkg)">
                <i class="fa fa-plus"></i>
                Install
            </button>
        </div>
    </div>
    `
};