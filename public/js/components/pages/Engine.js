import LiveLogs from '../LiveLogs.js';

export default {
	components: {
		"live-logs": LiveLogs
	},
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
		update_state(){
		},
		start: function(){
		},
		stop: function(){
		}
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
									<option>registry.gitlab.com/minetest/minetest/server:5.2.0</option>
									<option>registry.gitlab.com/minetest/minetest/server:5.3.0</option>
									<option>registry.gitlab.com/minetest/minetest/server:5.4.0</option>
									<option>registry.gitlab.com/minetest/minetest/server:latest</option>
									<option>buckaroobanzay/minetest:5.3.0-r5</option>
									<option>buckaroobanzay/minetest:5.3.0-r7</option>
								</select>
							</div>
							<div class="col-md-4">
								<div class="btn-group">
									<button class="btn btn-secondary" :disabled="state.created">
										<i class="fa fa-check"></i> Initialize
									</button>
									<button class="btn btn-success" :disabled="state.running || !state.created" v-on:click="start">
										<i class="fa fa-play"></i> Start
									</button>
									<button class="btn btn-danger" :disabled="!state.running || !state.created" v-on:click="stop">
										<i class="fa fa-stop"></i> Stop
									</button>
								</div>
							</div>
							<div class="col-md-4">
								Status
							</div>
						</div>
						<br>
						<div class="row">
							<div class="col-md-12">
								<live-logs name="minetest"/>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	`
};
