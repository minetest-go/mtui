import NavBar from './components/NavBar.js';
import Footer from './components/Footer.js';
import login_store from './store/login.js';

export default {
	data: () => login_store,
	components: {
		"nav-bar": NavBar,
		"nav-footer": Footer
	},
	template: /*html*/`
		<div>
			<nav-bar v-if="loggedIn"/>
			<div class="container-fluid">
				<br>
				<router-view></router-view>
			</div>
			<nav-footer/>
		</div>
	`
};
