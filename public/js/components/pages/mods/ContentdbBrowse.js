
import { get_packages } from "../../../api/cdb.js";

const store = Vue.reactive({
    packages: null
});

export default {
    data: function() {
        return store;
    },
    created: function() {
        if (!store.packages) {
            get_packages("game")
            .then(pkgs => store.packages = pkgs);
        }
    },
    template: /*html*/`
    <div>
        <h3>Browse mods</h3>
        <table class="table table-condensed">
            <thead>
                <tr>
                    <th>Author</th>
                    <th>Name</th>
                    <th>Action</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="pkg in packages" v-if="packages">
                    <td>{{pkg.author}}</td>
                    <td>{{pkg.name}}</td>
                    <td>
                        <router-link :to="'/mods/cdb/install/' + pkg.author + '/' + pkg.name" class="btn btn-success">
                            <i class="fa-solid fa-plus"></i>
                            Install
                        </router-link>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
    `
};