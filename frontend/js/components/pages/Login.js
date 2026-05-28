import { login, logout, is_logged_in } from "../../service/login.js";
import { get_onboard_status } from "../../api/onboard.js";
import { has_feature } from "../../service/features.js";
import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

export default {
    data: function() {
        return {
            username: "",
            password: "",
            otp_code: "",
            busy: false,
            error_message: "",
            can_oboard: false,
            breadcrumb: [START, {
                name: "Login",
                icon: "user",
                link: "/login"
            }]
        };
    },
    components: {
        "default-layout": DefaultLayout
    },
    computed: {
        is_logged_in: is_logged_in,
        validInput: function(){
            return this.username != "" && this.password != "";
        }
    },
    mounted: function() {
        get_onboard_status().then(s => this.can_oboard = s);
    },
    methods: {
        has_feature: has_feature,
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
    <default-layout icon="user" title="Login" :breadcrumb="breadcrumb">
        <div class="row">
            <div class="col-md-4"></div>
            <div class="col-md-4 card" style="padding: 20px;">
                <h4>Login</h4>
                <form @submit.prevent="login">
                    <input type="text"
                        class="form-control"
                        placeholder="Username"
                        :disabled="is_logged_in"
                        v-model="username"/>
                    <input type="password"
                        class="form-control"
                        placeholder="Password"
                        :disabled="is_logged_in"
                        v-model="password"/>
                    <input type="text"
                        maxlength="6"
                        class="form-control"
                        placeholder="OTP Code (optional)"
                        :disabled="is_logged_in"
                        v-if="has_feature('otp')"
                        v-model="otp_code"/>
                    <button class="btn btn-primary w-100" v-if="!is_logged_in" type="submit" :disabled="!validInput">
                        <i class="fa-solid fa-right-to-bracket"></i>
                        Login
                        <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                        <span class="badge bg-danger">{{error_message}}</span>
                    </button>
                    <a class="btn btn-secondary w-100" v-if="is_logged_in" v-on:click="logout">
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
                    <i class="fa fa-user"></i>
                    <b>Onboarding:</b>
                    Please set up an initial admin-user on the
                    <router-link to="/onboard">Onboarding</router-link>
                    page
                </div>
                <div class="alert alert-success" v-if="has_feature('signup') && !can_oboard">
                    <i class="fa fa-user"></i>
                    <b>Signup:</b>
                    You can create a new account on the
                    <router-link to="/signup">signup</router-link>
                    page
                </div>
            </div>
            <div class="col-md-4"></div>
        </div>
    </default-layout>
    `
};