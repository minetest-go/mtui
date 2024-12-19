import Userprofile from "../profile/Userprofile.js";
import DefaultLayout from '../layouts/DefaultLayout.js';
import { START, PLAYER_SEARCH } from '../Breadcrumb.js';

export default {
    props: ["name"],
	data: function() {
        return {
            breadcrumb: [START, PLAYER_SEARCH, {
                name: `Player profile for ${this.name}`,
                icon: "user",
                link: `/profile/${this.name}`
            }]
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
