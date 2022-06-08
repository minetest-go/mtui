
export default {
	template: /*html*/`
	<div>
		<div class="text-center">
			<h4>Start page</h4>
			<hr/>
			<router-link to="/shell" class="btn btn-primary">
				<i class="fa-solid fa-terminal"></i> Shell
			</router-link>
			&nbsp;
			<router-link to="/profile" class="btn btn-primary">
				<i class="fa fa-user"></i> Profile
			</router-link>
			&nbsp;
			<a class="btn btn-secondary" href="https://github.com/minetest-go/mtui" target="new">
				<i class="fa-brands fa-github"></i> Source
			</a>
		</div>
	</div>
	`
};
