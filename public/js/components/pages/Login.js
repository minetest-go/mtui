import login_store from "../../store/login.js";
import { login, logout } from "../../service/login.js";

export default {
    data: function() {
        return {
            username: "",
            password: "",
            busy: false,
            error_message: "",
            login_store: login_store
        };
    },
    computed: {
        validInput: function(){
            return this.username != "" && this.password != "";
        }
    },
    methods: {
        login: function() {
            this.busy = true;
            this.error_message = "";
            login(this.username, this.password)
            .then(success => {
                this.busy = false;
                if (!success) {
                    this.error_message = "Login failed!";
                } else {
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
                    <button class="btn btn-primary w-100" v-if="!login_store.loggedIn" v-on:click="login" :disabled="!validInput">
                        Login
                        <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                        <span class="badge bg-danger">{{error_message}}</span>
                    </button>
                    <a class="btn btn-secondary w-100" v-if="login_store.loggedIn" v-on:click="logout">
                        Logout
                        <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                    </a>
                </form>
            </div>
            <div class="col-md-4"></div>
        </div>
    `
};