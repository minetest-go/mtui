export const list = () => fetch("api/mods").then(r => r.json());

export const scan = () => fetch("api/mods/scan", {method: "POST"}).then(r => r.json());
