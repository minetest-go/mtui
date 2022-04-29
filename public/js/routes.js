import Start from './components/Start.js';
import Mods from './components/Mods.js';
import Engine from './components/Engine.js';

export default [{
	path: "/", component: Start
}, {
	path: "/mods", component: Mods
}, {
	path: "/engine", component: Engine
}];
