import events from "../events.js";

const store = Vue.reactive({
    players: []
});

events.on("player_stats", function(pstats) {
    store.players = pstats;
});

export default store;