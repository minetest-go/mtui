export const signup = (username, password) => fetch("api/signup", {
    method: "POST",
    body: JSON.stringify({
        username: username,
        password: password
    })
}).then(r => r.text());