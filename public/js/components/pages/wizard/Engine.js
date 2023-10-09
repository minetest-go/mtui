import EngineSelection from "../services/EngineSelection.js";
import { engine } from "../../../service/service.js";

export default {
    components: {
        "engine-selection": EngineSelection
    },
    data: () => engine.store,
    methods: {
        create: () => engine.create()
    },
    computed: {
        complete: function() {
            return this.status && this.status.created;
        }
    },
    template: /*html*/`
    <div>
        <h4>Engine selection</h4>
        <div class="row" v-if="complete">
            <div class="col-12">
                <div class="alert alert-success">
                    <i class="fa fa-check"></i> Engine installed: <b>{{status.version}}</b>
                </div>
            </div>
        </div>
        Select the minetest engine to install
        (the engine can be controlled in the <router-link to="/services/engine">engine</router-link> page)
        <div class="row">
            <div class="col-8">
                <engine-selection/>
            </div>
            <div class="col-4">
                <button class="btn btn-primary w-100" :disabled="complete || busy || !version" v-on:click="create">
                    <i class="fa fa-check"></i>
                    Install
                    <i class="fa fa-spinner fa-spin" v-if="busy"></i>
                </button>
            </div>
        </div>
    </div>
    `
};