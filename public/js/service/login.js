
import { get_claims as fetch_claims, login as api_login, logout as api_logout } from '../api/login.js';
import { connect } from '../ws.js';
import events, { EVENT_LOGGED_IN } from '../events.js';

const store = Vue.reactive({
    claims: null
});

export const is_logged_in = () => store.claims != null;
export const get_claims = () => store.claims;

export const check_login = () => fetch_claims().then(c => {
    store.claims = c;
    if (c) {
        events.emit(EVENT_LOGGED_IN, c);
    }
    return c;
});

export const has_priv = priv => store.claims && store.claims.privileges.find(e => e == priv);

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
        store.claims = null;
    });
};