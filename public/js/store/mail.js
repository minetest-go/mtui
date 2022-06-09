
const store = Vue.reactive({
    mails: [],
    contacts: {},
    unread_count: Vue.computed(() => store.mails.filter(m => m.unread).length)
});

export default store;