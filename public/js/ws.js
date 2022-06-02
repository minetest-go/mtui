
import events from "./events.js";

export default function() {
    const wsurl = window.location.protocol.replace("http", "ws") +
        "//" + window.location.host +
        window.location.pathname.substring(0, window.location.pathname.lastIndexOf("/")) +
        "/api/ws";
    
    const ws = new WebSocket(wsurl);
    ws.onmessage = function(e) {
        const cmd = JSON.parse(e.data);
        events.emit(cmd.type, cmd.data);
    };
}