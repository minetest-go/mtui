import { get_stats } from "../../../api/service.js";
import format_size from "../../../util/format_size.js";

export default {
    props: ["servicename"],
    data: function() {
        return {
            live: true,
            stats: null
        };
    },
    mounted: function() {
        this.update();
    },
    unmounted: function() {
        this.live = false;
    },
    methods: {
        format_size,
        update: function() {
            const start = Date.now();
            get_stats(this.servicename)
            .then(s => {                
                if (!this.live) {
                    return;
                }
                this.stats = s;

                // re-schedule in one second interval at most
                const delay = Math.max(0, 1000 - (Date.now() - start));
                setTimeout(() => this.update(), delay);
            });
        }
    },
    computed: {
        memory_percent: function() {
            return this.stats.memory_usage / this.stats.memory_max * 100;
        }
    },
    template: /*html*/`
        <div v-if="stats">
            CPU:
            <div class="progress">
                <div class="progress-bar" style="width: 25%" v-bind:style="{width: stats.cpu_percent + '%'}">
                    {{ Math.floor(stats.cpu_percent*10)/10}} %
                </div>
            </div>
            Memory:
            <div class="progress">
                <div class="progress-bar" style="width: 25%" v-bind:style="{width: memory_percent + '%'}">
                    {{ format_size(stats.memory_usage) }}
                </div>
            </div>
        </div>
    `
};