
import { get_package } from "../../api/cdb.js";

export default {
    data: function() {
        return {
            author: this.$route.params.author,
            name: this.$route.params.name,
            pkg: null
        };
    },
    created: function() {
        get_package(this.author, this.name)
        .then(p => this.pkg = p);
    },
    template: /*html*/`
    <div v-if="pkg">
        <h4>
            {{pkg.title}}
            <small class="text-muted">by {{pkg.author}}</small>
        </h4>
        <div class="row">
            <div class="col-2">
                <div class="card">
                    <div class="card-header">
                        Featured
                    </div>
                    <div class="card-body">
                        <h5 class="card-title">Special title treatment</h5>
                        <p class="card-text">With supporting text below as a natural lead-in to additional content.</p>
                        <a href="#" class="btn btn-primary">Go somewhere</a>
                    </div>
                </div>
            </div>
            <div class="col-10">
                <div class="card">
                    <div class="card-header">
                        Featured
                    </div>
                    <div class="card-body">
                        <h5 class="card-title">Special title treatment</h5>
                        <p class="card-text">With supporting text below as a natural lead-in to additional content.</p>
                        <a href="#" class="btn btn-primary">Go somewhere</a>
                    </div>
                </div>
            </div>
        </div>
    </div>
    `
};