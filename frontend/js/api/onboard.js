
export const get_onboard_status = () => fetch("api/onboard").then(r => r.json());

export const create_initial_user = (username, password) => fetch("api/onboard", {
    method: "POST",
    body: JSON.stringify({
        username: username,
        password: password
    })
}).then(r => r.text());