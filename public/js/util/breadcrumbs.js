export const START = { name: "Start", icon: "home", link: "/" };
export const PLAYER_SEARCH = { name: "Player search", icon: "magnifying-glass", link: "/playersearch" };
export const PLAYER_INFO = name => {
    return { name: `Player profile for ${name}`, icon: "user", link: `/profile/${name}` };
};
export const MAIL = { icon: "envelope", name: "Mail", link: "/mail" };
export const MODERATION = { icon: "hammer", name: "Moderation" };
export const SERVICES = { icon: "gears", name: "Services" };
export const ADMINISTRATION = { icon: "screwdriver-wrench", name: "Administration" };
export const MODS = { name: "Mods", icon: "cubes", link: "/mods" };
export const CDB = { name: "ContentDB", icon: "box-open", link: "/cdb/browse" };
export const FILEBROWSER = { name: "Filebrowser", icon: "folder", link: "/filebrowser/" };
export const RESTART_CONDITIONS = { name: "Restart conditions", icon: "refresh", link: "/restart-conditions/" };
export const CDB_DETAIL = (author, name) => {
    return { name: `'${author}/${name}'`, icon: "box-open", link: `/cdb/detail/${author}/${name}` };
};
