import { reactive } from "vue";

// state of the global confirmation dialog (rendered once by ConfirmDialog.vue in App.vue)
export const store = reactive({
    show: false,
    title: "",
    message: "",
    confirm_label: "",
    confirm_icon: "",
    danger: true,
    resolve: null
});

// asks the user for confirmation, resolves to true if confirmed, false otherwise
//
// usage: if (!await confirm({ title: "Remove mod?", message: "...", confirm_label: "Remove" })) return;
export function confirm(opts) {
    // a dialog that is still open counts as cancelled
    close(false);

    return new Promise(resolve => {
        Object.assign(store, {
            show: true,
            title: opts.title || "Are you sure?",
            message: opts.message || "",
            confirm_label: opts.confirm_label || "Confirm",
            confirm_icon: opts.confirm_icon || "",
            danger: opts.danger !== false,
            resolve
        });
    });
}

export function close(result) {
    const resolve = store.resolve;
    store.show = false;
    store.resolve = null;
    if (resolve) {
        resolve(result);
    }
}
