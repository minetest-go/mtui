import { changepw } from "../api/login.js";

export default {
    props: ["username"],
    data: function() {
        return {
            oldpassword: "",
            newpassword: "",
            newpassword2: "",
            busy: false,
            error: false
        };
    },
    methods: {
        changepw: function() {
            this.busy = true;
            this.error = false;

            changepw(this.username, this.oldpassword, this.newpassword)
            .then(success => {
                this.busy = false;
                this.error = !success;
                if (success) {
                    this.oldpassword = "";
                    this.newpassword = "";
                    this.newpassword2 = "";
                }
            });
        }
    },
    template: /*html*/`
        <form @submit.prevent="changepw">
            <input type="password" class="form-control" placeholder="Old password" v-model="oldpassword"/>
            <input type="password" class="form-control" placeholder="New password" v-model="newpassword"/>
            <input type="password" class="form-control" placeholder="New password (again)" v-model="newpassword2"/>
            <button class="btn btn-primary w-100" type="submit"
                :disabled="newpassword == '' || oldpassword == '' || newpassword != newpassword2">
                Change password
                <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                <span class="badge bg-danger" v-if="error">Change failed</span>
            </button>
        </form>
    `
};