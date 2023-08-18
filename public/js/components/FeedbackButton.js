
export default {
    props: ["fn", "type", "disabled"],
    data: function() {
        return {
            busy: false
        };
    },
    computed: {
        classes: function() {
            const c = {
                "btn": true
            };
            c["btn-" + this.type] = true;
            return c;
        }
    },
    methods: {
        click: function() {
            this.busy = true;
            const res = this.fn();
            if (res && typeof(res.then) == "function") {
                res
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
        </button>
    `
};