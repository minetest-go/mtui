
export default {
    props: ["fn", "type", "disabled"],
    data: function() {
        return {
            busy: false,
            error_msg: ""
        };
    },
    computed: {
        classes: function() {
            const c = {
                "btn": true
            };
            if (this.error_msg) {
                c["btn-warning"] = true;
            } else {
                c["btn-" + this.type] = true;
            }
            return c;
        }
    },
    methods: {
        click: function() {
            this.busy = true;
            this.error_msg = "";
            const res = this.fn();
            if (res && typeof(res.then) == "function") {
                res
                .catch(e => this.error_msg = ""+e)
                .finally(() => this.busy = false);
            } else {
                this.busy = false;
            }
        }
    },
    template: /*html*/`
        <button :class="classes" v-on:click="click" :disabled="busy || disabled">
            <slot></slot>
            <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
            <span v-if="error_msg" class="badge text-bg-danger">
                {{error_msg}}
            </span>
        </button>
    `
};