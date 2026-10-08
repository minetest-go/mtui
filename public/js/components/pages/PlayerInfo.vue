<script>
import Userprofile from "../profile/Userprofile.vue";
import DefaultLayout from '../layouts/DefaultLayout.vue';
import { START, PLAYER_SEARCH, PLAYER_INFO } from '../../util/breadcrumbs.js';

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
    }
};
</script>

<template>
  <default-layout icon="user" title="Player profile" :breadcrumb="breadcrumb">
  	<user-profile :username="name"/>
  </default-layout>
</template>
