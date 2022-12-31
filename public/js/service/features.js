
import store from "../store/features.js";
import { get_features, set_feature as set } from "../api/features.js";

export const check_features = () => {
    return get_features()
    .then(f => Object.keys(f).forEach(k => store[k] = f[k]));
};

export const set_feature = (name, enabled) => {
    set({ name: name, enabled: enabled })
    .then(() => check_features());
};

export const has_feature = name => store[name] && store[name].enabled;