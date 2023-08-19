
export const store = Vue.reactive({
    settingtypes: {},
    settings: {}
});

export const ordered_settings = Vue.computed(() => {
    // link -> []setting
    const ordered_settings = {};

    Object.keys(store.settingtypes).forEach(key => {
        const st = store.settingtypes[key];
        st.link = st.topic.join("/");
        st.key = key;
        st.current = store.settings[key];

        if (!ordered_settings[st.link]) {
            ordered_settings[st.link] = [];
        }
        ordered_settings[st.link].push(st);
    });

    return ordered_settings;
});

export const topics = Vue.computed(() => {
    return Object
    .keys(ordered_settings.value)
    .sort((a,b) => a > b);
});