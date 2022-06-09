
import store from "../store/features.js";
import { get_features } from "../api/features.js";

export const check_features = () => {
    return get_features()
    .then(f => Object.keys(f).forEach(k => store[k] = f[k]));
};

export const has_feature = name => store[name];