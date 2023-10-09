import { store, update_state, start, stop, restart, create, remove } from "../../../service/engine.js";

import ServiceStatus from "./ServiceStatus.js";
import ServiceLogs from "./ServiceLogs.js";

export default {
    props: ["servicename"],
	components: {
		"service-logs": ServiceLogs,
		"service-status": ServiceStatus,
	},
	data: () => store,
	methods: {
		update_state: update_state,
		start: start,
		stop: stop,
		restart: restart,
		remove: remove,
		create: create
	},
	template: /*html*/`
    <div class="row">
        <div class="col-md-12">
            <div class="card">
                <div class="card-header">
                    Engine
                    <service-status :status="status"/>
                    <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                </div>
                <div class="card-body">
                    <div class="row">
                        <div v-if="versions && status" class="col-md-4">
                            <select class="form-control" v-model="version" :disabled="!status || status.created">
                                <option v-for="(image, version) in versions" :value="version">{{version}}</option>
                            </select>
                        </div>
                        <div class="col-md-4" v-if="status">
                            <div class="btn-group">
                                <button class="btn btn-secondary" :disabled="busy || status.created || !version" v-on:click="create">
                                    <i class="fa fa-check"></i> Install
                                </button>
                                <button class="btn btn-secondary" :disabled="busy || status.running || !status.created" v-on:click="remove">
                                    <i class="fa fa-times"></i> Uninstall
                                </button>
                                <button class="btn btn-success" :disabled="busy || !status.created || status.running" v-on:click="start">
                                    <i class="fa fa-play"></i> Start
                                </button>
                                <button class="btn btn-warning" :disabled="busy || !status.created || !status.running" v-on:click="stop">
                                    <i class="fa fa-stop"></i> Stop
                                </button>
                                <button class="btn btn-warning" :disabled="busy || !status.created || !status.running" v-on:click="restart">
                                    <i class="fa fa-rotate-right"></i> Restart
                                </button>
                            </div>
                        </div>
                        <div class="col-md-4" v-if="status">
                            Status:
                            <span v-if="!status.created" class="badge bg-secondary">no service installed</span>
                            <span v-if="status.created && !status.running" class="badge bg-primary">service installed ({{version}})</span>
                            <span v-if="status.running" class="badge bg-success">service running ({{version}})</span>
                        </div>
                    </div>
                    <br>
                </div>
            </div>
        </div>
    </div>
    &nbsp;
    <div class="row">
        <div class="col-md-12">
            <service-logs :running="status && status.running" :servicename="servicename"/>
        </div>
    </div>
	`
};
