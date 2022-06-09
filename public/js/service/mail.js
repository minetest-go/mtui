import { list, contacts } from "../api/mail.js";
import store from "../store/mail.js";

export const fetch_mails = () => list().then(l => store.mails = l);
export const fetch_contacts = () => contacts().then(c => store.contacts = c);
