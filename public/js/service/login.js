
import login_store from '../store/login.js';
import { get_claims, login as api_login, logout as api_logout } from '../api/login.js';
import { connect } from '../ws.js';
import events, { EVENT_LOGGED_IN } from '../events.js';

export const check_login = () => get_claims().then(c => {
    login_store.claims = c;
    if (c) {
        events.emit(EVENT_LOGGED_IN, c);
    }
    return c;
});

export const has_priv = priv => login_store.claims && login_store.claims.privileges.find(e => e == priv);

export const login = (username, password, otp_code) => {
    return api_login(username, password, otp_code)
    .then(success => {
        if (success) {
            // reconnect websocket connection
            connect();
            return check_login();
        }
    });
};

export const logout = () => {
    return api_logout()
    .then(() => {
        // reconnect websocket connection
        connect();
        login_store.claims = null;
    });
};