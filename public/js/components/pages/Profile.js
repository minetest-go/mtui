import { get_claims } from '../../service/login.js';
import Userprofile from "../Userprofile.js";
import DefaultLayout from '../layouts/DefaultLayout.js';
import { START, PLAYER_SEARCH } from '../Breadcrumb.js';

export default {
    data: function() {
        return {
            breadcrumb: [START, PLAYER_SEARCH, {
                name: "Current player profile",
                icon: "user",
                link: "/profile"
            }]
        };
    },
    computed: {
        get_claims: get_claims
    },
    components: {
        "user-profile": Userprofile,
        "default-layout": DefaultLayout
    },
    template: /*html*/`
    <default-layout icon="user" title="Current player profile" :breadcrumb="breadcrumb">
        <user-profile :username="get_claims.username" v-if="get_claims"/>
    </default-layout>
    `
};