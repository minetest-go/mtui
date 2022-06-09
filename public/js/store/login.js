
const store = Vue.reactive({
    claims: null,
    loggedIn: Vue.computed(() => store.claims != null)
});

export default store;
