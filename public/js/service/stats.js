import { fetch_stats } from "../api/stats.js";

export const store = Vue.reactive({
    max_lag: null,
    uptime: null,
    time_of_day: null,
    player_count: null,
    players: [],
    maintenance: null
});

var handle;

export const get_stats = () => fetch_stats().then(s => Object.assign(store, s));

export const start_polling = () => handle = setInterval(get_stats, 2000);
export const stop_polling = () => clearInterval(handle);

export const get_player_count = () => store.player_count;