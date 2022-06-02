
export const execute_chatcommand = (playername, command) => fetch("api/bridge/execute_chatcommand", {
    method: "POST",
    body: JSON.stringify({
        playername: playername,
        command: command
    })
})
.then(r => r.json());