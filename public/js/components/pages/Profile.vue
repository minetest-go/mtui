<script>
import { get_claims } from '../../service/login.js';
import Userprofile from "../profile/Userprofile.vue";
import DefaultLayout from '../layouts/DefaultLayout.vue';
import { START, PLAYER_SEARCH } from '../../util/breadcrumbs.js';

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
    }
};
</script>

<template>
  <default-layout icon="user" title="Current player profile" :breadcrumb="breadcrumb">
      <user-profile :username="get_claims.username" :show_token_link="true" v-if="get_claims"/>
  </default-layout>
</template>
