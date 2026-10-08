import Start from './components/pages/Start.vue';
import Login from './components/pages/Login.vue';
import PlayerInfo from './components/pages/PlayerInfo.vue';
import AccessToken from './components/pages/AccessToken.vue';
import Profile from './components/pages/Profile.vue';
import Shell from './components/pages/Shell.vue';
import Lua from './components/pages/administration/Lua.vue';
import OnlinePlayers from './components/pages/OnlinePlayers.vue';
import PlayerSearch from './components/pages/PlayerSearch.vue';
import Mail from './components/pages/mail/Mail.vue';
import MailRead from './components/pages/mail/MailRead.vue';
import Compose from './components/pages/mail/Compose.vue';
import Skin from './components/pages/Skin.vue';
import Features from './components/pages/administration/Features.vue';
import Mediaserver from './components/pages/Mediaserver.vue';
import Log from './components/pages/Log.vue';
import Onboard from './components/pages/Onboard.vue';
import Xban from './components/pages/Xban.vue';
import MinetestConfig from './components/pages/administration/MinetestConfig.vue';
import UISettings from './components/pages/administration/UISettings.vue';
import Filebrowser from './components/pages/filebrowser/Filebrowser.vue';
import FileEditPage from './components/pages/filebrowser/FileEditPage.vue';
import Signup from './components/pages/Signup.vue';
import Help from './components/pages/Help.vue';
import ProfilerView from './components/pages/ProfilerView.vue';
import PrivEditor from './components/pages/PrivEditor.vue';

import Mods from './components/pages/mods/Mods.vue';
import ContentBrowse from './components/pages/cdb/Browse.vue';
import ContentdbDetail from './components/pages/cdb/Detail.vue';
import InstallCDB from './components/pages/cdb/Install.vue';
import Chat from './components/pages/Chat.vue';
import Mesecons from './components/pages/Mesecons.vue';
import Luacontroller from './components/pages/Luacontroller.vue';
import Play from './components/pages/Play.vue';
import RestartConditions from './components/pages/administration/RestartConditions.vue';

export default [{
	path: "/", component: Start
}, {
	path: "/restart-conditions", component: RestartConditions,
	meta: { requiredPriv: "server" }
}, {
	path: "/help", component: Help,
	meta: { requiredPriv: "server" }
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
	path:"/ui/settings", component: UISettings,
	meta: { requiredPriv: "server" }
}, {
	path: "/filebrowser/:pathMatch(.*)", component: Filebrowser, props: true,
	meta: { requiredPriv: "server" }
}, {
	path: "/fileedit/:pathMatch(.*)", component: FileEditPage, props: true,
	meta: { requiredPriv: "server" }
}, {
	path: "/profiler-view/:pathMatch(.*)", component: ProfilerView, props: true,
	meta: { requiredPriv: "server" }
}];
