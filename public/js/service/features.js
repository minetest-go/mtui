import { get_features, set_feature as set } from "../api/features.js";

export const store = Vue.reactive({});

export async function check_features() {
    const f = await get_features()
    Object.assign(store, f);
}

export async function set_feature(name, enabled) {
    await set({ name: name, enabled: enabled })
    await check_features();
};

export const has_feature = name => store[name] && store[name].enabled;