import { signup, signup_captcha } from "../../api/signup.js";
import { login } from "../../service/login.js";
import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

export default {
    data: function() {
        return {
            username: "",
            password: "",
            password2: "",
            captcha_id: null,
            captcha: "",
            busy: false,
            error_message: "",
            breadcrumb: [START, {
                name: "Signup",
                icon: "user",
                link: "/signup"
            }]
        };
    },
    mounted: function() {
        signup_captcha().then(c => this.captcha_id = c);
    },
    components: {
        "default-layout": DefaultLayout
    },
    computed: {
        validInput: function(){
            return this.username != "" && this.password != "" && (this.password == this.password2) && this.captcha != "";
        }
    },
    methods: {
        create_user: function(){
            this.busy = true;
            this.error_message = "";
            signup({
                username: this.username,
                password: this.password,
                captcha: this.captcha,
                captcha_id: this.captcha_id
            })
            .then(err_msg => {
                this.busy = false;
                if (err_msg) {
                    this.error_message = err_msg;
                } else {
                    login(this.username, this.password)
                    .then(() => this.$router.push("/profile"));
                }
            });
        }
    },
    template: /*html*/`
    <default-layout icon="user" title="Signup" :breadcrumb="breadcrumb">
        <div class="row">
            <div class="col-md-4"></div>
            <div class="col-md-4 card" style="padding: 20px;">
                <h4>Signup</h4>
                <p>
                    Create a new account
                </p>
                <form @submit.prevent="create_user">
                    <input type="text"
                        class="form-control"
                        placeholder="Username"
                        v-model="username"/>
                    <input type="password"
                        class="form-control"
                        placeholder="Password"
                        v-model="password"/>
                    <input type="password"
                        class="form-control"
                        placeholder="Password (repeat)"
                        v-model="password2"/>
                    <img :src="'api/captcha/' + captcha_id + '.png'" v-if="captcha_id"/>
                    <input type="text"
                        class="form-control"
                        placeholder="Captcha"
                        v-model="captcha"/>
                    <button class="btn btn-primary w-100" type="submit" :disabled="!validInput">
                        <i class="fa-solid fa-user"></i>
                        Create account
                        <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                        <span class="badge bg-danger">{{error_message}}</span>
                    </button>
                </form>
            </div>
        </div>
    </default-layout>
    `
};