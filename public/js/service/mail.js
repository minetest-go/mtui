import { list_inbox, list_outbox, list_contacts } from "../api/mail.js";
import events, { EVENT_LOGGED_IN } from "../events.js";
import { has_feature } from "./features.js";

const store = Vue.reactive({
    busy: false,
    inbox: [],
    outbox: [],
    contacts: {}
});

export const get_unread_count = () => store.inbox.filter(m => !m.read).length;

export const fetch_mails = () => {
    store.busy = true;
    list_inbox()
    .then(l => l || [])
    .then(l => store.inbox = l.sort((a,b) => a.time < b.time))
    .then(() => list_outbox())
    .then(l => l || [])
    .then(l => store.outbox = l.sort((a,b) => a.time < b.time))
    .finally(store.busy = false);
};

export const fetch_contacts = () =>
    list_contacts()
    .then(c => store.contacts = c);

events.on(EVENT_LOGGED_IN, function() {
    if (has_feature("mail")) {
        fetch_mails();
        fetch_contacts();
    }
});

export const is_busy = () => store.busy;

export const get_mailbox = name => store[name]; // inbox, outbox

export const get_mail = id => store.inbox.concat(store.outbox).find(m => m.id == id);
