import { get_all, get_settingtypes } from "../api/mtconfig.js";
import events, { EVENT_LOGGED_IN } from "../events.js";
import { has_priv } from "./login.js";
import { has_feature } from "./features.js";

export const store = Vue.reactive({
    settingtypes: {},
    settings: {},
    search: "",
    only_configured: true
});

export const ordered_settings = Vue.computed(() => {
    // link -> []setting
    const ordered_settings = {};

    Object.keys(store.settingtypes).forEach(key => {
        if (store.only_configured && !store.settings[key]) {
            // not configured, hide
            return;
        }

        const st = store.settingtypes[key];
        st.link = st.topic.join("/");
        st.key = key;
        st.current = store.settings[key];

        if (store.search) {
            // search enabled
            const str = `${key},${st.short_description},${st.long_description}`;
            if (!str.includes(store.search)) {
                return;
            }
        }

        if (!ordered_settings[st.link]) {
            ordered_settings[st.link] = [];
        }
        ordered_settings[st.link].push(st);
    });

    return ordered_settings;
});

export const count = Vue.computed(() => Object
    .keys(ordered_settings.value)
    .map(key => ordered_settings.value[key].length)
    .reduce((a,c) => a + c, 0)
);

export const topics = Vue.computed(() => Object
    .keys(ordered_settings.value)
    .sort((a,b) => a > b)
);

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
    });
});