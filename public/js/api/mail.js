export const list = () => fetch("api/mail/list").then(r => r.json());
export const contacts = () => fetch("api/mail/contacts").then(r => r.json());