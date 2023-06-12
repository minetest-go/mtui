export const list_inbox = () => fetch("api/mail/folder/inbox").then(r => r.json());
export const list_outbox = () => fetch("api/mail/folder/outbox").then(r => r.json());
export const list_contacts = () => fetch("api/mail/contacts").then(r => r.json());

export const check_recipient = r => fetch(`api/mail/checkrecipient/${r}`).then(r => r.text());

export const send = msg => fetch(`api/mail`, {
    method: "POST",
    body: JSON.stringify(msg)
});

export const mark_read = msg => fetch(`api/mail/${msg.id}/read`, {
    method: "POST"
});

export const mark_unread = msg => fetch(`api/mail/${msg.id}/unread`, {
    method: "POST"
});

export const remove = msg => fetch(`api/mail/${msg.id}`, {
    method: "DELETE"
});