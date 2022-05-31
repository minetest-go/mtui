import Start from './components/pages/Start.js';
import Mods from './components/pages/Mods.js';
import Engine from './components/pages/Engine.js';
import Login from './components/pages/login.js';

export default [{
	path: "/", component: Start
}, {
	path: "/login", component: Login
}, {
	path: "/mods", component: Mods
}, {
	path: "/engine", component: Engine
}];
