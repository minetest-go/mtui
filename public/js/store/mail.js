
const store = Vue.reactive({
    mails: [],
    contacts: {},
    unread_count: Vue.computed(() => store.mails ? store.mails.filter(m => m.unread).length : 0)
});

export default store;