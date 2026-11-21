// curated settings for the "Common" tab of the server settings page
//
// label/help replace the engine's technical descriptions,
// "requires" hides a setting unless the given bool setting is enabled.
// Settings unknown to the settingtypes are skipped.
export default [{
	title: "Server",
	icon: "fa-solid fa-server",
	description: "How your server shows up for players.",
	settings: [{
		key: "server_name", label: "Server name",
		help: "Shown in the server list and when players join."
	}, {
		key: "server_description", label: "Description",
		help: "A short text about your server for the server list."
	}, {
		key: "motd", label: "Welcome message",
		help: "Shown in the chat to every player when they join."
	}, {
		key: "server_announce", label: "List in the public server list",
		help: "Lets anyone find your server in the Luanti server browser."
	}, {
		key: "server_address", label: "Public address", requires: "server_announce",
		help: "The domain name or IP address players use to connect."
	}, {
		key: "server_url", label: "Website", requires: "server_announce",
		help: "Shown as a link in the server list."
	}]
}, {
	title: "Players",
	icon: "fa-solid fa-users",
	description: "Who can join and what they are allowed to do.",
	settings: [{
		key: "max_users", label: "Maximum players",
		help: "How many players can be online at the same time."
	}, {
		key: "default_privs", label: "Privileges for new players",
		help: "Comma-separated list, for example: interact, shout. \"interact\" allows building, \"shout\" allows chatting."
	}, {
		key: "disallow_empty_password", label: "Require a password",
		help: "Players can't join or set an empty password."
	}, {
		key: "name", label: "Admin player",
		help: "The player with this name gets all privileges."
	}]
}, {
	title: "Gameplay",
	icon: "fa-solid fa-gamepad",
	description: "How the game plays.",
	settings: [{
		key: "creative_mode", label: "Creative mode",
		help: "Players have unlimited blocks and items (depends on the game)."
	}, {
		key: "enable_damage", label: "Players can be hurt",
		help: "Players lose health from falling, mobs, lava and so on, and can die."
	}, {
		key: "enable_pvp", label: "Players can hurt each other (PvP)",
		help: "Only has an effect when \"Players can be hurt\" is on."
	}, {
		key: "time_speed", label: "Speed of the day/night cycle",
		help: "72 = 20 minutes per day, 360 = 4 minutes, 1 = 24 hours (real time), 0 = time stands still."
	}, {
		key: "static_spawnpoint", label: "Spawn point",
		help: "Position where players (re)spawn as x,y,z, for example 0,10,0. Leave empty for a random spot."
	}]
}, {
	title: "World",
	icon: "fa-solid fa-mountain-sun",
	description: "Settings for the map.",
	settings: [{
		key: "mg_name", label: "Map generator",
		help: "Only used when a new world is created. Changing it has no effect on an existing world."
	}, {
		key: "fixed_map_seed", label: "World seed",
		help: "Only used when a new world is created. Leave empty for a random seed."
	}, {
		key: "enable_rollback_recording", label: "Record block changes",
		help: "Lets moderators undo griefing with /rollback. Uses more disk space."
	}]
}, {
	title: "Messages",
	icon: "fa-solid fa-message",
	description: "What players see when the server goes down.",
	settings: [{
		key: "kick_msg_shutdown", label: "Shutdown message",
		help: "Shown to players when the server shuts down or restarts."
	}, {
		key: "kick_msg_crash", label: "Crash message",
		help: "Shown to players when the server crashes."
	}, {
		key: "ask_reconnect_on_crash", label: "Offer to reconnect after a crash",
		help: "Players get a \"Reconnect\" button after a crash."
	}]
}];
