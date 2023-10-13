import { START, SERVICES } from "../../Breadcrumb.js";

import DefaultLayout from "../../layouts/DefaultLayout.js";
import ServiceView from "./ServiceView.js";

export default {
	components: {
		"default-layout": DefaultLayout,
		"service-view": ServiceView,
	},
	computed: {
		breadcrumb: function() {
			return [START, SERVICES, {
				name: "MTWeb",
				icon: "gear",
				link: "/services/mtweb"
			}];
		}
	},
	template: /*html*/`
	<default-layout icon="gear" title="MTWeb" :breadcrumb="breadcrumb">
        <div class="row">
            <div class="col-12">
		        <service-view servicename="mtweb"/>
            </div>
        </div>
	</default-layout>
	`
};
