export const get_join_password = () => fetch(`api/wasm/joinpassword`).then(r => r.json());
