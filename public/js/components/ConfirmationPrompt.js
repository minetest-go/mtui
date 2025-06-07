
export default {
    props: ["title", "show", "confirm_disabled"],
    emits: ["abort", "confirm"],
    template: /*html*/ `
    <div class="modal show" style="display: block;" tabindex="-1" v-show="show">
        <div class="modal-dialog">
            <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5">{{title}}</h1>
                <button type="button" class="btn-close" v-on:click="$emit('abort')"></button>
            </div>
            <div class="modal-body">
                <slot></slot>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-success" v-on:click="$emit('abort')">
                    Abort
                </button>
                <button type="button" class="btn btn-danger" v-on:click="$emit('confirm')" :disabled="confirm_disabled">
                    <slot name="confirm_button">
                        Confirm
                    </slot>
                </button>
            </div>
            </div>
        </div>
    </div>
    <div class="modal-backdrop fade show" v-show="show"></div>
    `
};