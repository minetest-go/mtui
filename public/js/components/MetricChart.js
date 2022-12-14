
import { get_metric_type, search_metrics } from "../api/metrics";

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
    },
    template: /*html*/`
    <div>
    </div>
    `
};