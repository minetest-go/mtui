import { engine } from "../../../service/service.js";

export default {
    data: () => engine.store,
    template: /*html*/`
    <div v-if="versions && status">
        <select class="form-control" v-model="version" :disabled="!status || status.created">
            <option v-for="(image, version) in versions" :value="version">{{version}}</option>
        </select>
    </div>
    `
};