
export const execute_lua = code => fetch("api/bridge/lua", {
    method: "POST",
    body: JSON.stringify({ code: code })
})
.then(r => r.json());