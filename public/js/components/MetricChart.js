
import { get_metric_type, search_metrics } from "../api/metrics.js";

export default {
    props: ["metric_name"],
    data: function() {
        return {
            metrics: [],
            metric_type: null
        };
    },
    mounted: function() {
        get_metric_type(this.metric_name)
        .then(mt => this.metric_type = mt);

        search_metrics({ name: this.metric_name })
        .then(m => this.metrics = m);

        const ctx = this.$refs.chart.getContext('2d');
        this.chart = new Chart(ctx, {
			type: 'bar',
			data: {
				datasets: [{
					label: "Mylabel",
					borderColor: 'red',
					backgroundColor: 'red',
					data: [],
					yAxisID: 'y1'
				}]
			}
        });
    },
    template: /*html*/`
    <div>
        <canvas ref="chart" style="height: 400px; width: 100%;"></canvas>
    </div>
    `
};