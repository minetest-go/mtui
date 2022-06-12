export const list = () => fetch("api/mail/list").then(r => r.json());
export const contacts = () => fetch("api/mail/contacts").then(r => r.json());
export const check_recipient = r => fetch(`api/mail/checkrecipient/${r}`).then(r => r.text());

export const send = (msg, recipient) => fetch(`api/mail/send/${recipient}`, {
    method: "POST",
    body: JSON.stringify(msg)
});

export const mark_read = msg => fetch(`api/mail/${msg.sender}/${msg.time}/read`, {
    method: "POST"
});

export const remove = msg => fetch(`api/mail/${msg.sender}/${msg.time}`, {
    method: "DELETE"
});