<script>
import { create_initial_user } from "../../api/onboard.js";
import AccountForm from "../AccountForm.vue";

export default {
    components: {
        "account-form": AccountForm
    },
    methods: {
        submit: function({ username, password }) {
            return create_initial_user(username, password);
        }
    }
};
</script>

<template>
  <account-form
      icon="plus"
      title="Onboard"
      link="/onboard"
      heading="Onboarding"
      description="Create an initial admin account"
      button_label="Create admin account"
      redirect="/"
      :submit="submit"/>
</template>
