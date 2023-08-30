import { get_all, get_settingtypes } from "../api/mtconfig.js";
import events, { EVENT_LOGGED_IN } from "../events.js";
import { has_priv } from "./login.js";
import { has_feature } from "./features.js";

export const store = Vue.reactive({
    settingtypes: {},
    settings: {},
    filtered_settings: {},
    filtered_count: 0,
    filtered_topics: []
});

export function save() {
    //TODO
}

export function apply_filter(filter) {
    // link -> []setting
    const filtered_settings = {};

    Object.keys(store.settingtypes).forEach(key => {
        if (filter.only_configured && !store.settings[key]) {
            // not configured, hide
            return;
        }

        const st = store.settingtypes[key];
        st.link = st.topic.join("/");
        st.key = key;
        st.current = store.settings[key];
        st.is_set = st.current != undefined;

        if (filter.search) {
            // search enabled
            const str = `${key},${st.short_description},${st.long_description}`;
            if (!str.toLowerCase().includes(filter.search.toLowerCase())) {
                return;
            }
        }

        if (!filtered_settings[st.link]) {
            filtered_settings[st.link] = [];
        }
        filtered_settings[st.link].push(st);
    });

    store.filtered_settings = filtered_settings;
    store.filtered_count = Object
        .keys(filtered_settings)
        .map(key => filtered_settings[key].length)
        .reduce((a,c) => a + c, 0);
    store.filtered_topics = Object
        .keys(filtered_settings)
        .sort((a,b) => a > b);

}

events.on(EVENT_LOGGED_IN, function() {
    if (!has_priv("server") || !has_feature("minetest_config")){
        return;
    }

    Promise.all([get_all(), get_settingtypes()])
    .then(result => {
        store.settings = result[0];
        store.settingtypes = result[1];

        // add dummy settingtypes for unknown/custom settings
        Object
        .keys(store.settings)
        .filter(s => !store.settingtypes[s])
        .forEach(key => {
            store.settingtypes[key] = {
                topic: ["custom"],
                type: "string",
                default: {
                    value: ""
                }
            };
        });

        apply_filter({ only_configured: true });
    });
});