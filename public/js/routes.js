import Start from './components/pages/Start.js';
import Mods from './components/pages/Mods.js';
import Engine from './components/pages/Engine.js';
import Login from './components/pages/Login.js';
import Profile from './components/pages/Profile.js';
import Shell from './components/pages/Shell.js';

export default [{
	path: "/", component: Start,
	meta: { requiredPriv: "interact" }
}, {
	path: "/login", component: Login
}, {
	path: "/profile", component: Profile,
	meta: { requiredPriv: "interact" }
}, {
	path: "/shell", component: Shell,
	meta: { requiredPriv: "interact" }
}, {
	path: "/mods", component: Mods,
	meta: { requiredPriv: "server" }
}, {
	path: "/engine", component: Engine,
	meta: { requiredPriv: "server" }
}];
