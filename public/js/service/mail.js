import { list_inbox, list_outbox, list_contacts } from "../api/mail.js";
import store from "../store/mail.js";
import events, { EVENT_LOGGED_IN } from "../events.js";
import { has_feature } from "./features.js";

export const get_unread_count = () => store.inbox.filter(m => !m.read).length;

export const fetch_mails = () =>
    list_inbox()
    .then(l => l || [])
    .then(l => store.inbox = l.sort((a,b) => a.time < b.time))
    .then(() => list_outbox())
    .then(l => l || [])
    .then(l => store.outbox = l.sort((a,b) => a.time < b.time));

export const fetch_contacts = () =>
    list_contacts()
    .then(c => store.contacts = c);

events.on(EVENT_LOGGED_IN, function() {
    if (has_feature("mail")) {
        fetch_mails();
        fetch_contacts();
    }
});