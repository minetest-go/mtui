import login_store from '../../store/login.js';
import Userprofile from "../Userprofile.js";

export default {
    data: () => login_store,
    components: {
        "user-profile": Userprofile
    },
    template: /*html*/`
        <user-profile :username="claims.username" v-if="claims"/>
    `
};