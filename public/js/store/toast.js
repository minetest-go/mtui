import { reactive } from "vue";

// notification stack, rendered by Toasts.vue in App.vue
export const store = reactive({
    list: []
});

const MAX_TOASTS = 5;
let next_id = 1;

export function dismiss(id) {
    const i = store.list.findIndex(t => t.id == id);
    if (i >= 0) {
        store.list.splice(i, 1);
    }
}

// adds a notification
// opts: type ("success", "danger", "info", "warning"), title, details (shown on demand), timeout (ms, 0 = sticky)
export function notify(message, opts = {}) {
    const type = opts.type || "success";

    // collapse repeated identical messages (polling errors for example) into one entry
    const existing = store.list.find(t => t.type == type && t.message == message && t.title == opts.title);
    if (existing) {
        existing.count++;
        return existing.id;
    }

    const toast = {
        id: next_id++,
        type,
        message,
        title: opts.title || "",
        details: opts.details || "",
        count: 1
    };
    store.list.push(toast);
    if (store.list.length > MAX_TOASTS) {
        store.list.shift();
    }

    const timeout = opts.timeout != null ? opts.timeout : (type == "danger" ? 0 : 3500);
    if (timeout > 0) {
        setTimeout(() => dismiss(toast.id), timeout);
    }
    return toast.id;
}

export const success = message => notify(message, { type: "success" });
export const error = (title, message, details) => notify(message, { type: "danger", title, details });
