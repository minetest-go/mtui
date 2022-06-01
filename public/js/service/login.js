
import login_store from '../store/login.js';
import { get_claims, login as api_login, logout as api_logout } from '../api/login.js';

export const check_login = () => get_claims().then(c => login_store.claims = c);

export const login = (username, password) => {
    return api_login(username, password)
    .then(success => {
        if (success) {
            return get_claims()
            .then(c => {
                login_store.claims = c;
                login_store.loggedIn = c != null;
            });
        }
    });
};

export const logout = () => {
    return api_logout()
    .then(() => {
        login_store.claims = null;
        login_store.loggedIn = false;
    });
};