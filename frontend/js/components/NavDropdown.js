
export default {
    props: ["icon", "name"],
    data: function() {
        return {
            open: false
        };
    },
    computed: {
        icon_classes: function() {
            return {
                fa: true,
                [`fa-${this.icon}`]: true
            };
        }
    },
    watch: {
        "$route": function() {
            this.open = false;
        }
    },
    template: /*html*/`
    <li class="nav-item dropdown" v-on:mouseleave="open = false">
        <a class="nav-link dropdown-toggle" v-on:click="open = true" v-on:mouseover="open = true">
            <i v-bind:class="icon_classes" v-if="icon"></i>
            {{name}}
        </a>		
        <ul class="dropdown-menu" v-bind:class="{'show': open}">
            <slot></slot>
        </ul>
    </li>
    `
};
