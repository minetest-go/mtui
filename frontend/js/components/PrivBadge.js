
export default {
    props: ["priv"],
    methods: {
        getPrivBadgeClass: function(priv) {
            if (priv == "server" || priv == "privs") {
                return { "badge": true, "bg-danger": true };
            } else if (priv == "ban" || priv == "kick") {
                return { "badge": true, "bg-primary": true };
            } else if (priv == "otp_enabled") {
                return { "badge": true, "bg-success": true };
            } else {
                return { "badge": true, "bg-secondary": true };
            }
        }
    },
    template: /*html*/`
    <span v-bind:class="getPrivBadgeClass(priv)">
        {{ priv }}
    </span>
    `
};