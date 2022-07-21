import events from "../events.js";
import { has_priv } from "../service/login.js";

const store = Vue.reactive({
    players: []
});

// plain player info
events.on("player_stats", function(pstats) {
    if (!has_priv("ban")) {
        store.players = pstats;
    }
});

// additional infos
events.on("player_stats_extra", function(pstats) {
    if (has_priv("ban")) {
        store.players = pstats;
    }
});

export default store;