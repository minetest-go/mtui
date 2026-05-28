import Userprofile from "../profile/Userprofile.js";
import DefaultLayout from '../layouts/DefaultLayout.js';
import { START, PLAYER_SEARCH, PLAYER_INFO } from '../Breadcrumb.js';

export default {
    props: ["name"],
	data: function() {
        return {
            breadcrumb: [START, PLAYER_SEARCH, PLAYER_INFO(this.name)]
        };
    },
	components: {
        "user-profile": Userprofile,
		"default-layout": DefaultLayout
    },
	template: /*html*/`
	<default-layout icon="user" title="Player profile" :breadcrumb="breadcrumb">
		<user-profile :username="name"/>
	</default-layout>
	`
};
