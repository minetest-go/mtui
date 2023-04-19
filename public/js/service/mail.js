import { list_inbox } from "../api/mail.js";
import store from "../store/mail.js";

export const fetch_mails = () => list_inbox()
    .then(l => l || [])
    .then(l => store.mails = l.sort((a,b) => a.time < b.time));
