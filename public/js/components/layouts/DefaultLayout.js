import Breadcrumb from "../Breadcrumb.js";

export default {
    props: ["title", "icon", "breadcrumb"],
    components: {
        "bread-crumb": Breadcrumb
    },
	template: /*html*/`
        <bread-crumb :items="breadcrumb" v-if="breadcrumb"/>
        <slot></slot>
	`
};
