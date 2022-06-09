export const list = () => fetch("api/mail/list").then(r => r.json());
export const contacts = () => fetch("api/mail/contacts").then(r => r.json());
export const check_recipient = r => fetch(`api/mail/checkrecipient/${r}`).then(r => r.text());

export const send = msg => fetch("api/mail/send", {
    method: "POST",
    body: JSON.stringify(msg)
}).then(r => r.text());

export const mark_read = msg => fetch(`api/mail/markread/${msg.sender}/${msg.time}`, {
    method: "POST"
}).then(r => r.text());