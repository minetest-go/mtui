import Start from './components/pages/Start.js';
import Mods from './components/pages/Mods.js';
import Engine from './components/pages/Engine.js';

export default [{
	path: "/", component: Start
}, {
	path: "/mods", component: Mods
}, {
	path: "/engine", component: Engine
}];
