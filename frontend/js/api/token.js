
export const generate_token = (expiry, privs) => fetch("api/token", {
    method: "POST",
    body: JSON.stringify({
        expiry,
        privs
    })
})
.then(r => r.text());