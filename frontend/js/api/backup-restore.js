
export const get_backup_restore_info = () => fetch("api/backup-restore").then(r => {
    if (r.status == 200) {
        return r.json();
    } else {
        return Promise.resolve(null);
    }
});

export const create_backup_restore_job = data => fetch("api/backup-restore/create", {
    method: "POST",
    body: JSON.stringify(data)
}).then(r => r.json());