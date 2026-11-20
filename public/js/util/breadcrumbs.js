export const START = { name: "Home", icon: "home", link: "/" };
export const PLAYER_SEARCH = { name: "Find players", icon: "magnifying-glass", link: "/playersearch" };
export const PLAYER_INFO = name => {
    return { name: `Player profile for ${name}`, icon: "user", link: `/profile/${name}` };
};
export const MAIL = { icon: "envelope", name: "Mail", link: "/mail" };
export const MODERATION = { icon: "hammer", name: "Moderation" };
export const SERVICES = { icon: "gears", name: "Services" };
export const ADMINISTRATION = { icon: "screwdriver-wrench", name: "Server" };
export const MODS = { name: "Mods", icon: "cubes", link: "/mods" };
export const FILEBROWSER = { name: "Files", icon: "folder", link: "/filebrowser/" };
export const RESTART_CONDITIONS = { name: "Automatic restarts", icon: "refresh", link: "/restart-conditions/" };
