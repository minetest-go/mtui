export const get_latest_chat_messages = channel => fetch(`api/chat/${channel}`).then(r => r.json());

export const set_feature = msg => fetch("api/chat", {
    method: "POST",
    body: JSON.stringify(msg)
})
.then(r => r.json());