
export const enable_maintenance = () => fetch("api/maintenance", {
    method: "PUT"
});

export const disable_maintenance = () => fetch("api/maintenance", {
    method: "DELETE"
});

export const get_maintenance = () => fetch("api/maintenance").then(r => r.json());