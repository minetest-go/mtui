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
            const min = Math.floor(((this.time_of_day * 24) - this.hour) * 60);
            return min >= 10 ? min : "0" + min;
        }
    },
    template: /*html*/`
        <span v-if="max_lag > 0">
            <i class="fa-solid fa-signal" v-bind:style="{'color': signalColor}"></i>
            {{ Math.floor(max_lag*1000) }} ms
            <i class="fa-solid fa-users"></i>
            {{ player_count }}
            <i class="fa-solid fa-clock"></i>
            {{hour}}:{{minute}}
            <i class="fa-solid fa-sun" style="color: yellow;" v-if="hour >= 6 && hour < 18"></i>
            <i class="fa-solid fa-moon" style="color: lightblue;" v-else></i>
        </span>
    `
};