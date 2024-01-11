import { protected_fetch } from "./util.js";

export const get_controls_metadata = () => protected_fetch(`api/controls/metadata`);
export const get_controls_values = () => protected_fetch(`api/controls/values`);

export const set_control = (name, value) => protected_fetch(`api/mesecons`, {
    method: "POST",
    body: JSON.stringify({
        name,
        value
    })
});
