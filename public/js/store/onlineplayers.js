import events, { EVENT_PLAYER_STATS, EVENT_PLAYER_STATS_EXTRA } from "../events.js";
import { has_priv } from "../service/login.js";

const store = Vue.reactive({
    players: []
});

// plain player info
events.on(EVENT_PLAYER_STATS, function(pstats) {
    if (!has_priv("ban")) {
        store.players = pstats;
    }
});

// additional infos
events.on(EVENT_PLAYER_STATS_EXTRA, function(pstats) {
    if (has_priv("ban")) {
        store.players = pstats;
    }
});

export default store;