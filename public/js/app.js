import NavBar from './components/NavBar.js';
import ErrorToast from './components/ErrorToast.js';

export default {
	components: {
		"nav-bar": NavBar,
		"error-toast": ErrorToast
	},
	template: /*html*/`
		<div>
			<nav-bar/>
			<error-toast/>
			<div class="container-fluid">
				<br>
				<router-view></router-view>
			</div>
		</div>
	`
};
