import { START, SERVICES } from "../../Breadcrumb.js";

import DefaultLayout from "../../layouts/DefaultLayout.js";
import ServiceView from "./ServiceView.js";

export default {
	components: {
		"default-layout": DefaultLayout,
		"service-view": ServiceView
	},
	computed: {
		breadcrumb: function() {
			return [START, SERVICES, {
				name: "Minetest engine",
				icon: "gear",
				link: "/services/engine"
			}];
		}
	},
	template: /*html*/`
	<default-layout icon="gear" title="Minetest engine" :breadcrumb="breadcrumb">
		<service-view servicename="engine"/>
	</default-layout>
	`
};
