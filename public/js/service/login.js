
import login_store from '../store/login.js';
import { get_claims, login as api_login, logout as api_logout } from '../api/login.js';
import { connect } from '../ws.js';

export const check_login = () => get_claims().then(c => {
    login_store.claims = c;
    return c;
});

export const has_priv = priv => login_store.claims && login_store.claims.privileges.find(e => e == priv);

export const login = (username, password) => {
    return api_login(username, password)
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