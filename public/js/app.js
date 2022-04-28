import NavBar from './components/NavBar.js';
import Footer from './components/Footer.js';

export default {
	components: {
		"nav-bar": NavBar,
		"nav-footer": Footer
	},
	template: /*html*/`
		<div>
			<nav-bar/>
			<div class="container-fluid">
				<br>
				<router-view></router-view>
			</div>
			<nav-footer/>
		</div>
	`
};
