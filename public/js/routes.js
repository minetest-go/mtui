import Start from './components/pages/Start.js';
import Login from './components/pages/Login.js';
import PlayerInfo from './components/pages/PlayerInfo.js';
import Profile from './components/pages/Profile.js';
import Shell from './components/pages/Shell.js';
import Lua from './components/pages/Lua.js';
import OnlinePlayers from './components/pages/OnlinePlayers.js';
import PlayerSearch from './components/pages/PlayerSearch.js';
import Mail from './components/pages/mail/Mail.js';
import MailRead from './components/pages/mail/MailRead.js';
import Mods from './components/pages/Mods.js';
import Compose from './components/pages/Compose.js';
import Skin from './components/pages/Skin.js';
import Features from './components/pages/Features.js';
import Mediaserver from './components/pages/Mediaserver.js';
import Log from './components/pages/Log.js';
import Onboard from './components/pages/Onboard.js';
import Xban from './components/pages/Xban.js';
import OauthApps from './components/pages/OauthApps.js';
import OauthAppEdit from './components/pages/OauthAppEdit.js';

export default [{
	path: "/", component: Start,
	meta: { requiredPriv: "interact" }
}, {
	path: "/login", component: Login
}, {
	path: "/onboard", component: Onboard
}, {
	path: "/xban", component: Xban,
	meta: { requiredPriv: "ban" }
}, {
	path: "/features", component: Features,
	meta: { requiredPriv: "server" }
}, {
	path: "/log", component: Log,
	meta: { requiredPriv: "ban" }
}, {
	path: "/online-players", component: OnlinePlayers
}, {
	path: "/profile/:name", component: PlayerInfo
}, {
	path: "/playersearch", component: PlayerSearch
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
	path: "/mediaserver", component: Mediaserver,
	meta: { requiredPriv: "server" }
}, {
	path: "/mail", redirect: '/mail/box/inbox'
}, {
	path: "/mail/box/:boxname", component: Mail,
	meta: { requiredPriv: "interact" }
}, {
	path: "/mail/read/:id", component: MailRead,
	meta: { requiredPriv: "interact" }
}, {
	path: "/mail/compose", component: Compose,
	meta: { requiredPriv: "interact" }
}, {
	path: "/skin", component: Skin,
	meta: { requiredPriv: "interact" }
}, {
	path: "/oauth-apps", component: OauthApps,
	meta: { requiredPriv: "server" }
}, {
	path: "/oauth-apps/:id", component: OauthAppEdit,
	meta: { requiredPriv: "server" }
}];
