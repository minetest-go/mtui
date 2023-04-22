import { list_inbox, list_contacts } from "../api/mail.js";
import store from "../store/mail.js";

export const fetch_mails = () =>
    list_inbox()
    .then(l => l || [])
    .then(l => store.mails = l.sort((a,b) => a.time < b.time));

export const fetch_contacts = () =>
    list_contacts()
    .then(c => store.contacts = c);