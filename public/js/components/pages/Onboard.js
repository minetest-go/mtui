import { create_initial_user } from "../../api/onboard.js";
import { login } from "../../service/login.js";
import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

export default {
    data: function() {
        return {
            username: "",
            password: "",
            busy: false,
            error_message: "",
            breadcrumb: [START, {
                name: "Onboard",
                icon: "plus",
                link: "/onboard"
            }]
        };
    },
    components: {
        "default-layout": DefaultLayout
    },
    computed: {
        validInput: function(){
            return this.username != "" && this.password != "";
        }
    },
    methods: {
        create_user: function(){
            this.busy = true;
            this.error_message = "";
            create_initial_user(this.username, this.password).then(err_msg => {
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
    <default-layout icon="plus" title="Onboard" :breadcrumb="breadcrumb">
        <div class="row">
            <div class="col-md-4"></div>
            <div class="col-md-4 card" style="padding: 20px;">
                <h4>Onboarding</h4>
                <p>
                    Create an initial admin account
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
                    <button class="btn btn-primary w-100" type="submit" :disabled="!validInput">
                        <i class="fa-solid fa-user"></i>
                        Create admin account
                        <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                        <span class="badge bg-danger">{{error_message}}</span>
                    </button>
                </form>

            </div>
        </div>
    </default-layout>
    `
};