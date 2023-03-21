import login_store from "../../store/login.js";
import { login, logout } from "../../service/login.js";
import { get_onboard_status } from "../../api/onboard.js";

export default {
    data: function() {
        return {
            username: "",
            password: "",
            otp_code: "",
            busy: false,
            error_message: "",
            can_oboard: false,
            login_store: login_store
        };
    },
    computed: {
        validInput: function(){
            return this.username != "" && this.password != "";
        }
    },
    mounted: function() {
        get_onboard_status().then(s => this.can_oboard = s);
    },
    methods: {
        login: function() {
            this.busy = true;
            this.error_message = "";
            login(this.username, this.password, this.otp_code)
            .then(success => {
                this.busy = false;
                if (!success) {
                    // no luck
                    this.error_message = "Login failed!";
                } else if (this.$route.query.return_to) {
                    // return url
                    window.location.href = atob(this.$route.query.return_to);

                } else {
                    // go to base page
                    this.$router.push("/");
                }
            });
        },
        logout: function() {
            this.busy = true;
            logout()
            .then(() => this.busy = false);
        }
    },
    template: /*html*/`
        <div class="row">
            <div class="col-md-4"></div>
            <div class="col-md-4 card" style="padding: 20px;">
                <h4>Login</h4>
                <form @submit.prevent="login">
                    <input type="text"
                        class="form-control"
                        placeholder="Username"
                        :disabled="login_store.loggedIn"
                        v-model="username"/>
                    <input type="password"
                        class="form-control"
                        placeholder="Password"
                        :disabled="login_store.loggedIn"
                        v-model="password"/>
                    <input type="text"
                        maxlength="6"
                        class="form-control"
                        placeholder="OTP Code (optional)"
                        :disabled="login_store.loggedIn"
                        v-model="otp_code"/>
                    <button class="btn btn-primary w-100" v-if="!login_store.loggedIn" type="submit" :disabled="!validInput">
                        <i class="fa-solid fa-right-to-bracket"></i>
                        Login
                        <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                        <span class="badge bg-danger">{{error_message}}</span>
                    </button>
                    <a class="btn btn-secondary w-100" v-if="login_store.loggedIn" v-on:click="logout">
                        <i class="fa-solid fa-right-from-bracket"></i>
                        Logout
                        <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                    </a>
                </form>
                &nbsp;
                <div class="alert alert-primary">
                    <b>Note:</b> you can also use <mark>/mtui_tan</mark> ingame to create a temporary password
                </div>
                <div class="alert alert-success" v-if="can_oboard">
                    <b>Onboarding:</b>
                    Please set up an initial admin-user on the
                    <router-link to="/onboard">Onboarding</router-link>
                    page
                </div>
            </div>
            <div class="col-md-4"></div>
        </div>
    `
};