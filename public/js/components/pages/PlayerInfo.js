import { get } from "../../api/playerinfo.js";

export default {
	data: function() {
		return {
			info: null
		};
	},
	mounted: function() {
		const name = this.$route.params.name;
		get(name).then(info => this.info = info);
	},
	methods: {
		getPrivBadgeClass: function(priv) {
			if (priv == "server" || priv == "privs") {
				return { "badge": true, "bg-danger": true };
			} else if (priv == "ban" || priv == "kick") {
				return { "badge": true, "bg-primary": true };
			} else {
				return { "badge": true, "bg-secondary": true };
			}
		}
	},
	template: /*html*/`
	<div v-if="info">
		<h3>
			Profile for
			<small class="text-muted">
				{{ info.name }}
			</small>
		</h3>
		<div class="row">
			<div class="col-md-4">
				<div class="card">
					<div class="card-header">
						Privileges
					</div>
					<div class="card-body">
						<ul>
							<li v-for="priv in info.privs">
								<span v-bind:class="getPrivBadgeClass(priv)">{{ priv }}</span>
							</li>
						</ul>
					</div>
				</div>
			</div>
			<div class="col-md-4">
				<div class="card">
					<div class="card-header">
						Health stats
					</div>
					<div class="card-body">
						{{ info.health }} / {{ info.breath }}
					</div>
				</div>
			</div>
			<div class="col-md-4">
				<div class="card">
					<div class="card-header">
						Login stats
					</div>
					<div class="card-body">
						{{ info.last_login }} / {{ info.first_login }}
					</div>
				</div>
			</div>
		</div>
	</div>
	`
};
