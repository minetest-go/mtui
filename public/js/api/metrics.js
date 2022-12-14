
export const get_metric_types = () => fetch(`api/metric_types`).then(r => r.json());

export const get_metric_type = name => fetch(`api/metric_types/${name}`).then(r => r.json());

export const search_metrics = (q) => fetch(`api/metrics/search`, {
    method: "POST",
    body: JSON.stringify(q)
}).then(r => r.json());

export const count_metrics = (q) => fetch(`api/metrics/count`, {
    method: "POST",
    body: JSON.stringify(q)
})
.then(r => r.text())
.then(c => +c);