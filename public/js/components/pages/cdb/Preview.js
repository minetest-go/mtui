
export default {
    props: ["pkg"],
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
        </div>
    </div>
    `
};