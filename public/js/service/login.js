
import login_store from '../store/login.js';
import { get_claims, login as api_login, logout as api_logout } from '../api/login.js';
import { connect } from '../ws.js';

import { has_feature } from './features.js';
import { fetch_contacts, fetch_mails } from './mail.js';

function update_mails() {
    // check mails if available
    if (has_feature("mail")) {
        fetch_mails();
        fetch_contacts();
    }
}

export const check_login = () => get_claims().then(c => {
    login_store.claims = c;
    if (c) {
        update_mails();
    }
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