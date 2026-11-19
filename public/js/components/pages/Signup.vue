<script>
import { signup } from "../../api/signup.js";
import AccountForm from "../AccountForm.vue";

export default {
    components: {
        "account-form": AccountForm
    },
    methods: {
        submit: function(fields) {
            return signup(fields);
        }
    }
};
</script>

<template>
  <account-form
      icon="user"
      title="Signup"
      link="/signup"
      heading="Signup"
      description="Create a new account"
      button_label="Create account"
      redirect="/profile"
      :captcha_enabled="true"
      :submit="submit"/>
</template>
