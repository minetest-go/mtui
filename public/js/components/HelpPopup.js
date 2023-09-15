
export default {
    props: ["title"],
    data: function() {
        return {
            show: false
        };
    },
    template: /*html*/`
    <span style="margin: 5px;">
        <i class="fa-solid fa-circle-question" v-on:click="show = true" title="Click for help"></i>
    </span>
    <div class="modal show" style="display: block;" tabindex="-1" v-show="show">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h1 class="modal-title fs-5">{{title}}</h1>
                    <button type="button" class="btn-close" v-on:click="show = false"></button>
                </div>
                <div class="modal-body">
                    <slot></slot>
                </div>
            </div>
        </div>
    </div>
    `
};