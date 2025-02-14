import { protected_fetch } from "./util.js";

export const get_uimod_storage = key => fetch(`api/uimod/storage/${key}`).then(r => {
    if (r.status == 200)
        return r.text();
    else if (r.status == 404)
        return "";
    else
        return Promise.reject("http error: " + r.status);
});

export const set_uimod_storage = (key, value) => fetch(`api/uimod/storage/${key}`, {
    method: "POST",
    body: value
}).then(r => r.text());

export const get_priv_infos = () => protected_fetch(`api/uimod/priv_info`);

export const get_chatcommand_infos = () => protected_fetch(`api/uimod/chatcommand_info`);
