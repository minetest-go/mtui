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
            <div class="row">
                <div class="col-md-4">
                    <div class="card">
                        <div class="card-header">
                            Privileges
                        </div>
                        <div class="card-body">
                            <ul v-for="priv in claims.privileges">
                                <li>
                                    <span v-bind:class="getPrivBadgeClass(priv)">{{ priv }}</span>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
                <div class="col-md-4">
                    <div class="card">
                        <div class="card-header">
                            Login stats
                        </div>
                        <div class="card-body">
                        </div>
                    </div>
                </div>
                <div class="col-md-4">
                    <div class="card">
                        <div class="card-header">
                            Password
                        </div>
                        <div class="card-body">
                            <input type="password" class="form-control" placeholder="Old password"/>
                            <input type="password" class="form-control" placeholder="New password"/>
                            <input type="password" class="form-control" placeholder="New password (again)"/>
                            <a class="btn btn-primary w-100">
                                Change password
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `
};