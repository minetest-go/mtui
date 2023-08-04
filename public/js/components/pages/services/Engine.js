import { create, start, get_status, stop } from "../../../api/service_engine.js";

export default {
	data: function(){
		return {
			state: {
				created: true
			}
		};
	},
	created: function(){
		this.update_state();
	},
	methods: {
		update_state(){},
		start: start,
		stop: stop,
		create: function(){
			create({
				version: "5.7.0"
			});
		},
		get_status: get_status
	},
	template: /*html*/`
		<div class="row" v-if="state">
			<div class="col-md-12">
				<div class="card">
					<div class="card-header">
						Engine
					</div>
					<div class="card-body">
						<div class="row">
							<div class="col-md-4">
								<select class="form-control">
									<option>5.7.0</option>
									<option>5.6.0</option>
								</select>
							</div>
							<div class="col-md-4">
								<div class="btn-group">
									<button class="btn btn-secondary" v-on:click="create">
										<i class="fa fa-check"></i> Create
									</button>
									<button class="btn btn-success" v-on:click="start">
										<i class="fa fa-play"></i> Start
									</button>
									<button class="btn btn-danger" v-on:click="stop">
										<i class="fa fa-stop"></i> Stop
									</button>
									<button class="btn btn-info" v-on:click="get_status">
										<i class="fa fa-info"></i> Info
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
	`
};
