import { create, remove, start, get_status, stop, get_versions } from "../../../api/service_engine.js";
import EngineLogs from "./EngineLogs.js";
import events, { EVENT_LOGGED_IN } from "../../../events.js";
import { has_feature } from "../../../service/features.js";
import { has_priv } from "../../../service/login.js";

export const store = Vue.reactive({
	versions: null,
	busy: false,
	status: null,
	version: ""
});

events.on(EVENT_LOGGED_IN, function() {
	if (has_feature("docker") && has_priv("server")) {
		get_versions()
		.then(v => store.versions = v);
	}
});

export default {
	components: {
		"engine-logs": EngineLogs
	},
	data: function(){
		return store;
	},
	created: function(){
		this.update_state();
	},
	methods: {
		update_state(){
			store.busy = true;
			get_status()
			.then(s => store.status = s)
			.then(() => store.version = store.status.version)
			.finally(() => store.busy = false);
		},
		start: function(){
			store.busy = true;
			start()
			.then(() => get_status())
			.then(s => store.status = s)
			.finally(() => store.busy = false);
		},
		stop: function(){
			store.busy = true;
			stop()
			.then(() => get_status())
			.then(s => store.status = s)
			.finally(() => store.busy = false);
		},
		remove: function(){
			store.busy = true;
			remove()
			.then(() => get_status())
			.then(s => store.status = s)
			.finally(() => store.busy = false);
		},
		create: function(){
			store.busy = true;
			create({version: store.version})
			.then(() => get_status())
			.then(s => store.status = s)
			.finally(() => store.busy = false);
		}
	},
	template: /*html*/`
		<div class="row">
			<div class="col-md-12">
				<div class="card">
					<div class="card-header">
						Engine
						<i class="fa fa-play" v-if="status && status.running" style="color: green;"></i>
						<i class="fa fa-stop" v-if="status && !status.running" style="color: red;"></i>
						<i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
					</div>
					<div class="card-body">
						<div class="row">
							<div class="col-md-4" v-if="versions && status">
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
								</div>
							</div>
							<div class="col-md-4" v-if="status">
								Status:
								<span v-if="!status.created" class="badge bg-secondary">no engine installed</span>
								<span v-if="status.created && !status.running" class="badge bg-primary">engine installed ({{version}})</span>
								<span v-if="status.running" class="badge bg-success">engine running ({{version}})</span>
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
				<engine-logs/>
			</div>
		</div>
	`
};
