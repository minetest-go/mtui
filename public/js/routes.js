import Start from './components/pages/Start.js';
import Login from './components/pages/Login.js';
import PlayerInfo from './components/pages/PlayerInfo.js';
import Profile from './components/pages/Profile.js';
import Shell from './components/pages/Shell.js';
import Lua from './components/pages/Lua.js';
import OnlinePlayers from './components/pages/OnlinePlayers.js';
import Mail from './components/pages/Mail.js';
import MailRead from './components/pages/MailRead.js';
import Mods from './components/pages/Mods.js';
import Compose from './components/pages/Compose.js';
import Skin from './components/pages/Skin.js';
import Features from './components/pages/Features.js';

export default [{
	path: "/", component: Start,
	meta: { requiredPriv: "interact" }
}, {
	path: "/login", component: Login
}, {
	path: "/features", component: Features,
	meta: { requiredPriv: "server" }
}, {
	path: "/online-players", component: OnlinePlayers
}, {
	path: "/player/:name", component: PlayerInfo
}, {
	path: "/profile", component: Profile,
	meta: { requiredPriv: "interact" }
}, {
	path: "/shell", component: Shell,
	meta: { requiredPriv: "interact" }
}, {
	path: "/lua", component: Lua,
	meta: { requiredPriv: "server" }
}, {
	path: "/mods", component: Mods,
	meta: { requiredPriv: "server" }
}, {
	path: "/mail", component: Mail,
	meta: { requiredPriv: "interact" }
}, {
	path: "/mail/read/:sender/:time", component: MailRead,
	meta: { requiredPriv: "interact" }
}, {
	path: "/mail/compose", component: Compose,
	meta: { requiredPriv: "interact" }
}, {
	path: "/skin", component: Skin,
	meta: { requiredPriv: "interact" }
}];
