import { create, remove, start, get_status, stop, get_versions, get_logs } from "../../../api/service_engine.js";

const store = Vue.reactive({
	versions: null,
	busy: false,
	status: null,
	version: "",
	logs: "",
	logs_since: Date.now() - (1000*60*60) //since an hour
});

function update_logs(){
	if (!store.status || !store.status.created) {
		// skip log-fetching if no container created
		return;
	}

	// fetch and shift window
	const now = Date.now();
	get_logs(store.logs_since, now)
	.then(l => {
		if (l.out){
			store.logs += l.out;
		}
		if (l.err){
			store.logs += l.err;
		}
		store.logs_since = now + 1;
	});
}

export default {
	data: function(){
		return store;
	},
	created: function(){
		if (!store.status){
			this.update_state();
		}
		if (!store.versions) {
			get_versions()
			.then(v => store.versions = v);
		}
		this.log_update_handle = setInterval(update_logs, 1000);
	},
	unmounted: function() {
		clearInterval(this.log_update_handle);
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
								<select class="form-control" v-model="version">
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
							<div class="col-md-4">
								Status
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
				<div class="card">
					<div class="card-header">
						Logs
					</div>
					<div class="card-body">
						<pre style="height: 400px; background: grey;">{{logs}}</pre>
					</div>
				</div>
			</div>
		</div>
	`
};
