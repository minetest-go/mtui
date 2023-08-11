
export const get_themes = () => fetch(`api/themes`).then(r => r.json());