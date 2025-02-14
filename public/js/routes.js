import Start from './components/pages/Start.js';
import Login from './components/pages/Login.js';
import PlayerInfo from './components/pages/PlayerInfo.js';
import AccessToken from './components/pages/AccessToken.js';
import Profile from './components/pages/Profile.js';
import Shell from './components/pages/Shell.js';
import Lua from './components/pages/administration/Lua.js';
import OnlinePlayers from './components/pages/OnlinePlayers.js';
import PlayerSearch from './components/pages/PlayerSearch.js';
import Mail from './components/pages/mail/Mail.js';
import MailRead from './components/pages/mail/MailRead.js';
import Compose from './components/pages/mail/Compose.js';
import Skin from './components/pages/Skin.js';
import Features from './components/pages/administration/Features.js';
import Mediaserver from './components/pages/Mediaserver.js';
import Log from './components/pages/Log.js';
import Onboard from './components/pages/Onboard.js';
import Xban from './components/pages/Xban.js';
import MinetestConfig from './components/pages/administration/MinetestConfig.js';
import UISettings from './components/pages/administration/UISettings.js';
import Filebrowser from './components/pages/filebrowser/Filebrowser.js';
import FileEditPage from './components/pages/filebrowser/FileEditPage.js';
import Signup from './components/pages/Signup.js';
import Help from './components/pages/Help.js';
import ProfilerView from './components/pages/ProfilerView.js';
import PrivEditor from './components/pages/PrivEditor.js';

import EngineService from './components/pages/services/Engine.js';
import MatterbridgeService from './components/pages/services/Matterbridge.js';
import MapserverService from './components/pages/services/Mapserver.js';

import Mods from './components/pages/mods/Mods.js';
import ContentBrowse from './components/pages/cdb/Browse.js';
import ContentdbDetail from './components/pages/cdb/Detail.js';
import InstallCDB from './components/pages/cdb/Install.js';
import Wizard from './components/pages/wizard/Wizard.js';
import Chat from './components/pages/Chat.js';
import Mesecons from './components/pages/Mesecons.js';
import Luacontroller from './components/pages/Luacontroller.js';
import Play from './components/pages/Play.js';
import RestartConditions from './components/pages/administration/RestartConditions.js';
import BackupRestore from './components/pages/administration/BackupRestore.js';

export default [{
	path: "/", component: Start,
	meta: { maintenance_page: true }
}, {
	path: "/restart-conditions", component: RestartConditions,
	meta: { requiredPriv: "server" }
}, {
	path: "/help", component: Help,
	meta: { requiredPriv: "server", maintenance_page: true }
}, {
	path: "/login", component: Login
}, {
	path: "/onboard", component: Onboard
}, {
	path: "/signup", component: Signup
}, {
	path: "/play", component: Play
}, {
	path: "/chat", component: Chat,
	meta: { requiredPriv: "shout" }
}, {
	path: "/wizard/:step", component: Wizard, props: true,
	meta: { requiredPriv: "server" }
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
	path: "/profile/:name", component: PlayerInfo, props: true,
}, {
	path: "/profile/:name/priveditor", component: PrivEditor, props: true,
	meta: { requiredPriv: "privs" }
}, {
	path: "/token", component: AccessToken
}, {
	path: "/playersearch", component: PlayerSearch
}, {
	path: "/profile", component: Profile,
	meta: { requiredPriv: "interact" }
}, {
	path: "/shell", component: Shell,
	meta: { requiredPriv: "interact" }
}, {
	path: "/mesecons", component: Mesecons,
	meta: { requiredPriv: "interact" }
}, {
	path: "/mesecons/luacontroller/:x/:y/:z", component: Luacontroller, props: true,
	meta: { requiredPriv: "interact" }
}, {
	path: "/lua", component: Lua,
	meta: { requiredPriv: "server" }
}, {
	path: "/mods", component: Mods,
	meta: { requiredPriv: "server" }
}, {
	path: "/backup", component: BackupRestore,
	meta: { requiredPriv: "server", maintenance_page: true }
}, {
	path: "/cdb/browse", component: ContentBrowse,
	meta: { requiredPriv: "server" }
}, {
	path: "/cdb/detail/:author/:name", component: ContentdbDetail, props: true,
	meta: { requiredPriv: "server" }
}, {
	path: "/cdb/install/:author/:name", component: InstallCDB, props: true,
	meta: { requiredPriv: "server" }
}, {
	path: "/mediaserver", component: Mediaserver,
	meta: { requiredPriv: "server" }
}, {
	path: "/mail", redirect: '/mail/box/inbox'
}, {
	path: "/mail/box/:boxname", component: Mail, props: true,
	meta: { requiredPriv: "interact" }
}, {
	path: "/mail/read/:id", component: MailRead, props: true,
	meta: { requiredPriv: "interact" }
}, {
	path: "/mail/compose", component: Compose,
	meta: { requiredPriv: "interact" }
}, {
	path: "/skin", component: Skin,
	meta: { requiredPriv: "interact" }
}, {
	path: "/minetest-config", component: MinetestConfig,
	meta: { requiredPriv: "server" }
}, {
	path: "/services/engine", component: EngineService,
	meta: { requiredPriv: "server" }
}, {
	path: "/services/matterbridge", component: MatterbridgeService,
	meta: { requiredPriv: "server" }
}, {
	path: "/services/mapserver", component: MapserverService,
	meta: { requiredPriv: "server" }
}, {
	path:"/ui/settings", component: UISettings,
	meta: { requiredPriv: "server" }
}, {
	path: "/filebrowser/:pathMatch(.*)", component: Filebrowser, props: true,
	meta: { requiredPriv: "server", maintenance_page: true }
}, {
	path: "/fileedit/:pathMatch(.*)", component: FileEditPage, props: true,
	meta: { requiredPriv: "server", maintenance_page: true }
}, {
	path: "/profiler-view/:pathMatch(.*)", component: ProfilerView, props: true,
	meta: { requiredPriv: "server" }
}];
