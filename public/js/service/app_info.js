import app_info from "../store/app_info.js";
import { get_appinfo } from "../api/app_info.js";

export const fetch_info = () => {
    get_appinfo().then(i => {
        Object.keys(i).forEach(k => app_info[k] = i[k]);
    });

    if (app_info.servername) {
        document.title = `${document.title} [${app_info.servername}]`;
    }
};