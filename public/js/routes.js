import Start from './components/pages/Start.js';
import Login from './components/pages/Login.js';
import Profile from './components/pages/Profile.js';
import Shell from './components/pages/Shell.js';
import OnlinePlayers from './components/pages/OnlinePlayers.js';
import Mail from './components/pages/Mail.js';
import MailRead from './components/pages/MailRead.js';

export default [{
	path: "/", component: Start,
	meta: { requiredPriv: "interact" }
}, {
	path: "/login", component: Login
}, {
	path: "/online-players", component: OnlinePlayers,
	meta: { requiredPriv: "ban" },
}, {
	path: "/profile", component: Profile,
	meta: { requiredPriv: "interact" }
}, {
	path: "/shell", component: Shell,
	meta: { requiredPriv: "interact" }
}, {
	path: "/mail", component: Mail
}, {
	path: "/mail/read/:sender/:time", component: MailRead
}];
