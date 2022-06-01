import login_store from '../store/login.js';

const LoginPath = { path: "/login" };

export default function(router) {
    router.beforeEach((to) => {
        if (to.meta.loggedIn && !login_store.loggedIn) {
            return { path: '/login' };
        }

        if (to.meta.requiredPriv) {
            if (!login_store.loggedIn) {
                // quick login check
                return LoginPath;
            }

            if (!login_store.claims.privileges.find(p => p == to.meta.requiredPriv)){
                // check required priv
                return LoginPath;
            }
        }
    });   
}