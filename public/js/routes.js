import Start from './components/pages/Start.js';
import Mods from './components/pages/Mods.js';
import Engine from './components/pages/Engine.js';
import Login from './components/pages/login.js';

export default [{
	path: "/", component: Start,
	meta: { requiredPriv: "interact" }
}, {
	path: "/login", component: Login
}, {
	path: "/mods", component: Mods,
	meta: { requiredPriv: "server" }
}, {
	path: "/engine", component: Engine,
	meta: { requiredPriv: "server" }
}];
