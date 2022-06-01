import login_store from '../../store/login.js';

export default {
    data: () => login_store,
    methods: {
        getPrivBadgeClass: function(priv) {
            if (priv == "server" || priv == "privs") {
                return { "badge": true, "bg-danger": true };
            } else if (priv == "ban" || priv == "kick") {
                return { "badge": true, "bg-primary": true };
            } else {
                return { "badge": true, "bg-secondary": true };
            }
        }
    },
    template: /*html*/`
        <div>
            <h3>
                Profile for
                <small class="text-muted">
                    {{ claims.username }}
                </small>
            </h3>
            <hr>
            <h5>Privileges</h5>
            <ul v-for="priv in claims.privileges">
                <li>
                    <span v-bind:class="getPrivBadgeClass(priv)">{{ priv }}</span>
                </li>
            </ul>
        </div>
    `
};