
export default {
    props: ["author", "name", "pkg"],
    computed: {
        link: function() {
            if (this.pkg) {
                return `/cdb/detail/${this.pkg.author}/${this.pkg.name}`;
            } else {
                return `/cdb/detail/${this.author}/${this.name}`;
            }
        },
        text: function() {
            if (this.pkg) {
                return `${this.pkg.author}/${this.pkg.name}`;
            } else {
                return `${this.author}/${this.name}`;
            }
        }
    },
    template: /*html*/`
    <router-link :to="link">
        <i class="fa fa-cubes"></i>
        {{text}}
    </router-link>
    `
};