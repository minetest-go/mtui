import stats from "../store/stats.js";

export default {
    data: () => stats,
    computed: {
        signalColor: function() {
            if (this.max_lag < 200) return "green";
            if (this.max_lag < 500) return "yellow";
            return "red";
        },
        hour: function() {
            return Math.floor(this.time_of_day * 24);
        },
        minute: function() {
            const min = Math.floor((this.time_of_day % 1000) / 1000 * 60);
            return min >= 10 ? min : "0" + min;
        }
    },
    template: /*html*/`
        <span v-if="max_lag > 0">
            <i class="fa-solid fa-signal" v-bind:style="{'color': signalColor}"></i>
            {{ Math.floor(max_lag*10000)/10 }} ms
            <i class="fa-solid fa-users"></i>
            {{ player_count }}
            <i class="fa-solid fa-clock"></i>
            {{hour}}:{{minute}}
        </span>
    `
};