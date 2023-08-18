import { get_all } from "../api/mtconfig.js";
import { get_settingtypes } from "../api/mods.js";
import store from '../store/mtconfig.js';
import events, { EVENT_LOGGED_IN } from "../events.js";
import { has_priv } from "./login.js";
import { has_feature } from "./features.js";

events.on(EVENT_LOGGED_IN, function() {
    if (!has_priv("server") || !has_feature("minetest_config")){
        return;
    }

    Promise.all([get_all(), get_settingtypes()])
    .then(result => {
        //const cfg = result[0];
        const settingtypes = result[1];

        // link -> []setting
        const ordered_settings = {};

        settingtypes.forEach(st => {
            st.link = st.topic.join("/");

            if (!ordered_settings[st.link]) {
                ordered_settings[st.link] = [];
            }
            ordered_settings[st.link].push(st);
        });

        const topics = Object
            .keys(ordered_settings)
            .sort((a,b) => a > b);

        store.topics = topics;
        store.ordered_settings = ordered_settings;
    });
});