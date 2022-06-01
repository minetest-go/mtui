import login_store from "../../store/login.js";
import { login, logout } from "../../service/login.js";

const store = Vue.reactive({
    username: "",
    password: "",
    login_store: login_store
});

export default {
    data: () => store,
    methods: {
        login: function() {
            login(this.username, this.password);
        },
        logout: function() {
            logout();
        }
    },
    template: /*html*/`
        <div class="row">
            <div class="col-md-4"></div>
            <div class="col-md-4 card" style="padding: 20px;">
                <h4>Login</h4>
                <input type="text"
                    class="form-control"
                    placeholder="Username"
                    v-model="username"/>
                <input type="password"
                    class="form-control"
                    placeholder="Password"
                    v-model="password"/>
                <button class="btn btn-primary" v-if="!login_store.loggedIn" v-on:click="login">
                    Login
                </button>
                <button class="btn btn-secondary" v-if="login_store.loggedIn" v-on:click="logout">
                    Logout
                </button>
            </div>
            <div class="col-md-4"></div>
        </div>
    `
};