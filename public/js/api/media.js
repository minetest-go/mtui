
export const scan = () => fetch("api/media/scan", { method: "POST" }).then(r => r.json());
export const stats = () => fetch("api/media/stats").then(r => r.json());