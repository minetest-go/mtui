export const get_features = () => fetch("api/features").then(r => r.json());
