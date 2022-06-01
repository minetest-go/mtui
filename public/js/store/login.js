
const store = Vue.reactive({
    claims: null,
    loggedIn: Vue.computed(() => store.claims != null)
});

export default store;

export const has_priv = priv => store.claims && store.claims.privileges.find(e => e == priv);