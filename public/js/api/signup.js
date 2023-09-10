export const signup = (username, password) => fetch("api/onboard", {
    method: "POST",
    body: JSON.stringify({
        username: username,
        password: password
    })
}).then(r => r.text());