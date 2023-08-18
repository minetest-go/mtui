import { protected_fetch } from "./util.js";

export const get_status = () => protected_fetch(`api/xban/status`);

export const get_record = playername => protected_fetch(`api/xban/records/${playername}`);

export const get_records = () => protected_fetch(`api/xban/records`);

export const ban_player = (playername, reason) => protected_fetch(`api/xban/ban`, {
    method: "POST",
    body: JSON.stringify({
        playername: playername,
        reason: reason
    })
});

export const tempban_player = (playername, time, reason) => protected_fetch(`api/xban/tempban`, {
    method: "POST",
    body: JSON.stringify({
        playername: playername,
        time: time,
        reason: reason
    })
});

export const unban_player = playername => protected_fetch(`api/xban/unban`, {
    method: "POST",
    body: JSON.stringify({
        playername: playername
    })
});

export const cleanup = () => protected_fetch(`api/xban/cleanup`, {
    method: "POST",
});
