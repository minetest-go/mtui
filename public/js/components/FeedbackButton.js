
export default {
    props: ["fn", "classes", "content"],
    data: function() {
        return {
            busy: false
        };
    },
    methods: {
        click: function() {
            this.busy = true;
            const res = this.fn();
            if (res && typeof(res.then) == "function") {
                res
                .catch(e => console.log(e))
                .finally(() => this.busy = false);
            } else {
                this.busy = false;
            }
        }
    },
    template: /*html*/`
        <button :class="classes" v-on:click="click" :disabled="busy">
            <slot></slot>
        </button>
    `
}