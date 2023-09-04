import Userprofile from "../Userprofile.js";
import DefaultLayout from '../layouts/DefaultLayout.js';
import { START, PLAYER_SEARCH } from '../Breadcrumb.js';

export default {
	data: function() {
        return {
            breadcrumb: [START, PLAYER_SEARCH, {
                name: `Player profile for ${this.$route.params.name}`,
                icon: "user",
                link: `/profile/${this.$route.params.name}`
            }]
        };
    },
	components: {
        "user-profile": Userprofile,
		"default-layout": DefaultLayout
    },
	template: /*html*/`
	<default-layout icon="user" title="Player profile" :breadcrumb="breadcrumb">
		<user-profile :username="$route.params.name"/>
	</default-layout>
	`
};
