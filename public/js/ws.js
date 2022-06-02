
import events from "./events.js";

const wsurl = window.location.protocol.replace("http", "ws") +
    "//" + window.location.host +
    window.location.pathname.substring(0, window.location.pathname.lastIndexOf("/")) +
    "/api/ws";

var ws;

export function connect() {
    if (ws) {
        ws.onerror = null;
        ws.close();
    }
    
    ws = new WebSocket(wsurl);
    ws.onmessage = function(e) {
        const cmd = JSON.parse(e.data);
        events.emit(cmd.type, cmd.data);
    };

    // reconnect on error
    ws.onerror = function() {
        setTimeout(connect, 1000);
    };
}
