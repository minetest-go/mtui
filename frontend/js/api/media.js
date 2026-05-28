import { protected_fetch } from "./util.js";

export const scan = () => protected_fetch("api/media/scan", { method: "POST" });
export const stats = () => protected_fetch("api/media/stats");