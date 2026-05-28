export const resolve = ip => fetch(`api/geoip/${ip}`).then(r => r.json());
