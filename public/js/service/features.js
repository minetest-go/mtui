import { get_features, set_feature as set } from "../api/features.js";

export const store = Vue.reactive({});

export const check_features = () => {
    return get_features()
    .then(f => Object.assign(store, f));
};

export const set_feature = (name, enabled) => {
    set({ name: name, enabled: enabled })
    .then(() => check_features());
};

export const has_feature = name => store[name] && store[name].enabled;