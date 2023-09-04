import { protected_fetch } from "./util.js";

export const fetch_stats = () => protected_fetch("api/stats");
