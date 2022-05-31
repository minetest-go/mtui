
import { login, logout, get_claims } from "../../api/login.js";

export default {
    methods: {
        login: function() {
            login("test", "enter").then(r => console.log(r.status));
        },
        claims: function() {
            get_claims()
            .then(c => console.log(c))
            .catch(e => console.log(e));
        },
        logout: function() {
            logout()
        }
    },
    template: /*html*/`
        <div>
        <button v-on:click="login">Login</button>
        <button v-on:click="claims">Claims</button>
        <button v-on:click="logout">Logout</button>
        </div>
    `
};