
import store from "../store/error_toast.js";

export default {
    data: () => store,
    methods: {
        close: function() {
            this.message = "";
            this.title = "";
            this.url = "";
            this.status = 0;
        }
    },
    template: /*html*/`
    <div class="toast-container top-50 start-50 translate-middle">
        <div class="toast show" v-if="message">
            <div class="toast-header text-bg-danger">
                <strong class="me-auto">
                    <i class="fa-solid fa-triangle-exclamation"></i>
                    Error
                </strong>
                <small>
                    {{title}}
                    <span v-if="status">Code: {{status}}</span>
                </small>
                <button type="button" class="btn-close" v-on:click="close"></button>
            </div>
            <div class="toast-body">
                <p>
                    <b>Message:</b> {{message}}
                </p>
                <p>
                    <b>URL:</b> {{url}}
                </p>
            </div>
        </div>
    </div>
    `
};