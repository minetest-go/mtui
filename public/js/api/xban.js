
export const get_status = () => fetch(`api/xban/status`).then(r => r.json());

export const get_record = playername => fetch(`api/xban/records/${playername}`).then(r => r.json());

export const get_records = () => fetch(`api/xban/records`).then(r => r.json());

export const ban_player = (playername, reason) => fetch(`api/xban/ban`, {
    method: "POST",
    body: JSON.stringify({
        playername: playername,
        reason: reason
    })
}).then(r => r.json());

export const tempban_player = (playername, time, reason) => fetch(`api/xban/tempban`, {
    method: "POST",
    body: JSON.stringify({
        playername: playername,
        time: time,
        reason: reason
    })
}).then(r => r.json());

export const unban_player = playername => fetch(`api/xban/unban`, {
    method: "POST",
    body: JSON.stringify({
        playername: playername
    })
}).then(r => r.json());

export const cleanup = () => fetch(`api/xban/cleanup`, {
    method: "POST",
}).then(r => r.json());
