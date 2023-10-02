export const get_latest_chat_messages = channel => fetch(`api/chat/${channel}/latest`).then(r => r.json());

export const search_messages = (channel, from, to) => fetch(`api/chat/${channel}/${from}/${to}`).then(r => r.json());

export const send_message = msg => fetch("api/chat", {
    method: "POST",
    body: JSON.stringify(msg)
})
.then(r => r.json());