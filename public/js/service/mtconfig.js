import { reactive } from "vue";
import { get_all, get_settingtypes, set, remove } from "../api/mtconfig.js";
import events, { EVENT_LOGGED_IN } from "../events.js";
import { has_priv } from "./login.js";
import { has_feature } from "./features.js";

export const store = reactive({
    settingtypes: {},
    settings: {},
    filtered_settings: {},
    filtered_count: 0,
    filtered_topics: [],
    // unsaved edits: key -> setting
    pending: {},
    ready: false
});

export function save(key, setting) {
    return set(key, setting)
    .then(() => store.settings[key] = setting);
}

export const get_setting = key => store.settings[key];

export const is_ready = () => store.ready;

export function unset(key) {
    return remove(key)
    .then(() => {
        delete store.settings[key];
        delete store.pending[key];
    });
}

// fields that make up the value of a setting, per type
const value_fields = {
    v3f: ["x", "y", "z"],
    noise_params_2d: ["offset", "scale", "spread_x", "spread_y", "spread_z", "seed", "octaves", "persistence", "lacunarity"],
    noise_params_3d: ["offset", "scale", "spread_x", "spread_y", "spread_z", "seed", "octaves", "persistence", "lacunarity"]
};

export function same_setting(type, a, b) {
    return (value_fields[type] || ["value"])
        .every(f => String((a && a[f]) ?? "") == String((b && b[f]) ?? ""));
}

export const is_set = key => store.settings[key] != undefined;

// the saved value, or the default if not configured
export function saved_setting(key) {
    const st = store.settingtypes[key];
    return store.settings[key] || (st && st.default) || { value: "" };
}

// the value to show in the editor: unsaved edit or saved value
export const working_setting = key => store.pending[key] || saved_setting(key);

export const is_pending = key => store.pending[key] != undefined;

export function set_pending(key, setting) {
    const st = store.settingtypes[key];
    if (same_setting(st ? st.type : "string", setting, saved_setting(key))) {
        delete store.pending[key];
    } else {
        store.pending[key] = setting;
    }
}

export const discard_pending = key => delete store.pending[key];

export function discard_all_pending() {
    Object.keys(store.pending).forEach(key => delete store.pending[key]);
}

export const pending_keys = () => Object.keys(store.pending);

// saves all unsaved edits one after another (the backend rewrites the whole config file on each save)
export async function save_all_pending() {
    const keys = pending_keys();
    for (const key of keys) {
        const setting = Object.assign({}, store.pending[key]);
        if (setting.value != null) {
            setting.value = "" + setting.value;
        }
        await save(key, setting);
        delete store.pending[key];
    }
    return keys.length;
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

export function update_settings() {
    return Promise.all([get_all(), get_settingtypes()])
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
        store.ready = true;
    });
}

events.on(EVENT_LOGGED_IN, function() {
    if (!has_priv("server") || !has_feature("minetest_config")){
        return;
    }

    update_settings();
});