import { has_priv } from "../service/login.js";
import { has_feature } from "../service/features.js";

// sidebar navigation, grouped by audience
//
// item fields:
//   name, icon (fontawesome classes), link
//   match:   path prefixes that mark the item as active (default: [link])
//   exact:   only active on an exact path match
//   visible: function deciding if the item is shown (default: always)
//   badge:   name of a badge provided by the sidebar ("players", "mail")
//
// sections without visible items are hidden
export default [{
	name: "Play",
	items: [{
		name: "Home", icon: "fa-solid fa-house", link: "/", exact: true
	}, {
		name: "Play in browser", icon: "fa-solid fa-play", link: "/play",
		visible: () => has_feature("minetest_web")
	}, {
		name: "Chat", icon: "fa-solid fa-comment", link: "/chat",
		visible: () => has_priv("shout") && has_feature("chat")
	}, {
		name: "Mail", icon: "fa-solid fa-envelope", link: "/mail", badge: "mail",
		visible: () => has_feature("mail")
	}, {
		name: "Skin", icon: "fa-solid fa-user-astronaut", link: "/skin",
		visible: () => has_feature("skinsdb")
	}, {
		name: "Online players", icon: "fa-solid fa-users", link: "/online-players", badge: "players"
	}]
}, {
	name: "Tools",
	items: [{
		name: "Find players", icon: "fa-solid fa-magnifying-glass", link: "/playersearch",
		match: ["/playersearch", "/profile/"],
		visible: () => has_priv("interact")
	}, {
		name: "Commands", icon: "fa-solid fa-terminal", link: "/shell",
		visible: () => has_priv("interact") && has_feature("shell")
	}, {
		name: "Mesecons", icon: "fa-solid fa-microchip", link: "/mesecons",
		visible: () => has_priv("interact") && has_feature("mesecons")
	}]
}, {
	name: "Moderation",
	items: [{
		name: "Bans", icon: "fa-solid fa-ban", link: "/xban",
		visible: () => has_priv("ban") && has_feature("xban")
	}, {
		name: "Logs", icon: "fa-solid fa-list", link: "/log",
		visible: () => has_priv("ban")
	}]
}, {
	name: "Server",
	items: [{
		name: "Mods", icon: "fa-solid fa-cubes", link: "/mods",
		visible: () => has_priv("server") && has_feature("modmanagement")
	}, {
		name: "Server settings", icon: "fa-solid fa-sliders", link: "/minetest-config",
		visible: () => has_priv("server") && has_feature("minetest_config")
	}, {
		name: "Automatic restarts", icon: "fa-solid fa-rotate", link: "/restart-conditions",
		visible: () => has_priv("server")
	}, {
		name: "Files", icon: "fa-solid fa-folder", link: "/filebrowser/",
		match: ["/filebrowser", "/fileedit", "/profiler-view"],
		visible: () => has_priv("server")
	}, {
		name: "Media server", icon: "fa-solid fa-photo-film", link: "/mediaserver",
		visible: () => has_priv("server") && has_feature("mediaserver")
	}, {
		name: "Lua console", icon: "fa-solid fa-code", link: "/lua",
		visible: () => has_priv("server") && has_feature("luashell")
	}, {
		name: "Web UI features", icon: "fa-solid fa-toggle-on", link: "/features",
		visible: () => has_priv("server")
	}, {
		name: "Help", icon: "fa-solid fa-circle-question", link: "/help",
		visible: () => has_priv("server")
	}]
}];

export function is_active(item, path) {
	if (item.exact) {
		return path == item.link;
	}
	return (item.match || [item.link]).some(prefix => path.startsWith(prefix));
}
