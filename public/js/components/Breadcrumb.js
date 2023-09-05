
export default {
    props: ["items"],
    methods: {
        get_icon_class: function(item) {
            const cl = {};
            if (item.icon){
                cl.fa = true;
                cl["fa-" + item.icon] = true;
            }
            return cl;
        }
    },
    template: /*html*/`
    <nav>
        <ol class="breadcrumb text-bg-secondary" style="padding: 5px; border-radius: 5px;">
            <li class="breadcrumb-item" v-for="item in items">
                <router-link class="link-light" :to="item.link" v-if="item.link">
                    <i v-bind:class="get_icon_class(item)" v-if="item.icon"></i>
                    {{item.name}}
                </router-link>
                <span v-else>
                    <i v-bind:class="get_icon_class(item)" v-if="item.icon"></i>
                    {{item.name}}
                </span>
            </li>
        </ol>
    </nav>
    `
};

export const START = { name: "Start", icon: "home", link: "/" };
export const PLAYER_SEARCH = { name: "Player search", icon: "magnifying-glass", link: "/playersearch" };
export const MAIL = { icon: "envelope", name: "Mail", link: "/mail" };
export const MODERATION = { icon: "hammer", name: "Moderation" };
export const SERVICES = { icon: "gears", name: "Services" };
export const ADMINISTRATION = { icon: "screwdriver-wrench", name: "Administration" };
export const OAUTH_APPS = { icon: "passport", name: "OAuth apps", link: "/oauth-apps" };
export const MODS = { name: "Mods", icon: "cubes", link: "/mods" };
export const CDB = { name: "ContentDB", icon: "box-open", link: "/cdb/browse" };
export const FILEBROWSER = { name: "Filebrowser", icon: "folder", link: "/filebrowser/" };