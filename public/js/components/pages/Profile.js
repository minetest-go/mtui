import { get_claims } from '../../service/login.js';
import Userprofile from "../Userprofile.js";

export default {
    computed: {
        get_claims: get_claims
    },
    components: {
        "user-profile": Userprofile
    },
    template: /*html*/`
        <user-profile :username="get_claims.username" v-if="get_claims"/>
    `
};