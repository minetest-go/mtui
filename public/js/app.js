import NavBar from './components/NavBar.js';

export default {
	components: {
		"nav-bar": NavBar
	},
	template: /*html*/`
		<div>
			<nav-bar/>
			<div class="container-fluid">
				<br>
				<router-view></router-view>
			</div>
		</div>
	`
};
