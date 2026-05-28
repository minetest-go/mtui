export const get_features = () => fetch("api/features").then(r => r.json());

export const set_feature = feature => fetch("api/feature", {
    method: "POST",
    body: JSON.stringify(feature)
})
.then(r => r.json());