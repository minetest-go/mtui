
const store = Vue.reactive({
    busy: false,
    inbox: [],
    outbox: [],
    contacts: {}
});

export const get_mail = id => store.inbox.concat(store.outbox).find(m => m.id == id);

export default store;