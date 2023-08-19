import { get_all, get_settingtypes } from "../api/mtconfig.js";
import { store } from '../store/mtconfig.js';
import events, { EVENT_LOGGED_IN } from "../events.js";
import { has_priv } from "./login.js";
import { has_feature } from "./features.js";

events.on(EVENT_LOGGED_IN, function() {
    if (!has_priv("server") || !has_feature("minetest_config")){
        return;
    }

    Promise.all([get_all(), get_settingtypes()])
    .then(result => {
        store.settings = result[0];
        store.settingtypes = result[1];
    });
});