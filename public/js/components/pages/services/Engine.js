import { store, update_state, start, stop, restart, create, remove } from "../../../service/engine.js";
import { START, SERVICES } from "../../Breadcrumb.js";

import EngineStatus from "./EngineStatus.js";
import EngineSelection from "./EngineSelection.js";
import DefaultLayout from "../../layouts/DefaultLayout.js";
import EngineLogs from "./EngineLogs.js";
import HelpPopup from "../../HelpPopup.js";

export default {
	components: {
		"engine-logs": EngineLogs,
		"engine-status": EngineStatus,
		"engine-selection": EngineSelection,
		"default-layout": DefaultLayout,
		"help-popup": HelpPopup
	},
	data: () => store,
	computed: {
		breadcrumb: function() {
			return [START, SERVICES, {
				name: "Minetest engine",
				icon: "gear",
				link: "/services/engine"
			}];
		}
	},
	methods: {
		update_state: update_state,
		start: start,
		stop: stop,
		restart: restart,
		remove: remove,
		create: create
	},
	template: /*html*/`
	<default-layout icon="gear" title="Minetest engine" :breadcrumb="breadcrumb">
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
							<engine-selection class="col-md-4"/>
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
								<help-popup title="Installing a minetest engine">
									<p>
										Here you can choose which minetest-version you want to install.
										The engine can be up- or down-graded at any time (you have to stop it first).
										There is no data stored in the engine, it is just a process running in the background.
									</p>
									<p>
										Setting up a new engine: Select your desired minetest-version in the selection-box.
										Hit the "Install" button and wait for it to be downloaded and installed.
										Start the server with the "Start" button, there might be some errors
										in the console below if no game is installed yet.
									</p>
								</help-popup>
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
	</default-layout>
	`
};
