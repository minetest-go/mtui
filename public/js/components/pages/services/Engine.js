import { create, remove, start, stop } from "../../../api/service_engine.js";
import EngineLogs from "./EngineLogs.js";
import { store, update_state } from "../../../service/engine.js";
import EngineStatus from "./EngineStatus.js";

export default {
	components: {
		"engine-logs": EngineLogs,
		"engine-status": EngineStatus
	},
	data: () => store,
	methods: {
		update_state: update_state,
		start: function(){
			store.busy = true;
			start()
			.then(() => update_state());
		},
		stop: function(){
			store.busy = true;
			stop()
			.then(() => update_state());
		},
		remove: function(){
			store.busy = true;
			remove()
			.then(() => update_state());
		},
		create: function(){
			store.busy = true;
			create({version: store.version})
			.then(() => update_state());
		}
	},
	template: /*html*/`
		<div class="row">
			<div class="col-md-12">
				<div class="card">
					<div class="card-header">
						Engine
						<engine-status/>
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
				<engine-logs :running="status && status.running"/>
			</div>
		</div>
	`
};
