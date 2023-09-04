import { get_claims, is_logged_in } from '../service/login.js';
import { get_maintenance } from '../service/stats.js';

const LoginPath = { path: "/login" };

export default function(router) {
    router.beforeEach((to) => {
        if (get_maintenance()) {
            // maintenance mode enabled, only start and maintenance page available
            if (is_logged_in() && to.meta && to.meta.maintenance_page) {
                return;
            } else {
                return { path: "/" };
            }
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