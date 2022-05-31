
export const login = (username, password) => fetch("api/login", {
    method: "POST",
    body: JSON.stringify({
        username: username,
        password: password
    })
});

export const logout = () => fetch("api/login", {
    method: "DELETE"
});

export const get_claims = () => fetch("api/login").then(r => r.status == 200 ? r.json() : null);