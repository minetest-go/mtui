import { get_claims, is_logged_in } from '../service/login.js';

const LoginPath = { path: "/login" };

export default function(router) {
    router.beforeEach((to) => {
        if (!to.meta) {
            return;
        }

        if (to.meta.loggedIn && !is_logged_in()) {
                return LoginPath;
        }

        if (to.meta.requiredPriv) {
            if (!is_logged_in()) {
                // quick login check
                return LoginPath;
            }

            if (!get_claims().privileges.find(p => p == to.meta.requiredPriv)){
                // check required priv
                return LoginPath;
            }
        }
    });   
}