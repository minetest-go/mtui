<script>

export default {
    props: ["priv"],
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
    }
};
</script>

<template>
  <span v-bind:class="getPrivBadgeClass(priv)">
      {{ priv }}
  </span>
</template>
