
export const login = (username, password) => fetch("api/login", {
    method: "POST",
    body: JSON.stringify({
        username: username,
        password: password
    })
})
.then(r => r.status == 200);

export const logout = () => fetch("api/login", {
    method: "DELETE"
});

export const get_claims = () => fetch("api/login").then(r => r.status == 200 ? r.json() : null);

export const changepw = (username, oldpw, newpw) => fetch("api/changepw", {
    method: "POST",
    body: JSON.stringify({
        username: username,
        old_password: oldpw,
        new_password: newpw
    })
})
.then(r => r.status == 200);