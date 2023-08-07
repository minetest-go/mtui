
export default {
    created: function(){

    },
    data: function() {
        return {
            author: this.$route.params.author,
            name: this.$route.params.name,
        };
    },
    template: /*html*/`
    <div>
        <h3>
            Install mod
            <small class="text-muted">{{author}}/{{name}}</small>
        </h3>
    </div>
    `
};