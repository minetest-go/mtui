import { reactive } from "vue";
import { START, MAIL } from "../util/breadcrumbs.js";

export const store = reactive({
    body: "",
    subject: "",
    add_recipient_name: "",
    invalid_username: false,
    recipients: [],
    busy: false,
    mail_sent: false,
    breadcrumb: [START, MAIL, {
        name: "Compose mail",
        icon: "pen-to-square",
        link: "/mail/compose"
    }]
});
