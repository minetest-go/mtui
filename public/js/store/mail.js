
const store = Vue.reactive({
    busy: false,
    inbox: [],
    outbox: [],
    contacts: {},
    unread_count: Vue.computed(() => store.inbox.filter(m => !m.read).length)
});

export const get_mail = id => store.inbox.concat(store.outbox).find(m => m.id == id);

export default store;