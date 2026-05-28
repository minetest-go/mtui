import { get_service_by_name } from "../../../service/service.js";

import ServiceStatus from "./ServiceStatus.js";
import ServiceStats from "./ServiceStats.js";
import ServiceLogs from "./ServiceLogs.js";

export default {
    props: ["servicename"],
	components: {
		"service-logs": ServiceLogs,
		"service-status": ServiceStatus,
        "service-stats": ServiceStats
	},
	data: function() {
        return {
            service: get_service_by_name(this.servicename)
        };
    },
	methods: {
        get_service() {
            return get_service_by_name(this.servicename);
        }
	},
    computed: {
        running: function() {
            return this.service.store.status && this.service.store.status.running;
        }
    },
	template: /*html*/`
    <div class="row">
        <div class="col-md-12">
            <div class="card">
                <div class="card-header">
                    Service-Container
                    <service-status :status="service.store.status"/>
                    <i class="fa-solid fa-spinner fa-spin" v-if="service.store.busy"></i>
                </div>
                <div class="card-body">
                    <div class="row">
                        <div v-if="service.store.versions && service.store.status" class="col-md-4">
                            <select class="form-control" v-model="service.store.version" :disabled="!service.store.status || service.store.status.created">
                                <option v-for="(image, version) in service.store.versions" :value="version">{{version}}</option>
                            </select>
                        </div>
                        <div class="col-md-4" v-if="service.store.status">
                            <div class="btn-group">
                                <button class="btn btn-secondary" :disabled="service.store.busy || service.store.status.created || !service.store.version" v-on:click="get_service().create()">
                                    <i class="fa fa-check"></i> Install
                                </button>
                                <button class="btn btn-secondary" :disabled="service.store.busy || service.store.status.running || !service.store.status.created" v-on:click="get_service().remove()">
                                    <i class="fa fa-times"></i> Uninstall
                                </button>
                                <button class="btn btn-success" :disabled="service.store.busy || !service.store.status.created || service.store.status.running" v-on:click="get_service().start()">
                                    <i class="fa fa-play"></i> Start
                                </button>
                                <button class="btn btn-warning" :disabled="service.store.busy || !service.store.status.created || !service.store.status.running" v-on:click="get_service().stop()">
                                    <i class="fa fa-stop"></i> Stop
                                </button>
                                <button class="btn btn-warning" :disabled="service.store.busy || !service.store.status.created || !service.store.status.running" v-on:click="get_service().restart()">
                                    <i class="fa fa-rotate-right"></i> Restart
                                </button>
                            </div>
                        </div>
                        <div class="col-md-4" v-if="service.store.status">
                            Status:
                            <span v-if="!service.store.status.created" class="badge bg-secondary">no service installed</span>
                            <span v-if="service.store.status.created && !service.store.status.running" class="badge bg-primary">service installed ({{service.store.version}})</span>
                            <span v-if="service.store.status.running" class="badge bg-success">service running ({{service.store.version}})</span>
                        </div>
                    </div>
                    <br>
                    <service-stats v-if="running" :servicename="servicename"/>
                </div>
            </div>
        </div>
    </div>
    &nbsp;
    <div class="row">
        <div class="col-md-12">
            <service-logs :running="running" :servicename="servicename"/>
        </div>
    </div>
	`
};
