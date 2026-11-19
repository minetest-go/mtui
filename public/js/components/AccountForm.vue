<script>
import { signup_captcha } from "../api/signup.js";
import { login } from "../service/login.js";
import DefaultLayout from "./layouts/DefaultLayout.vue";
import { START } from "../util/breadcrumbs.js";

// shared form for the onboarding and signup pages
export default {
    props: {
        icon: String,
        title: String,
        heading: String,
        description: String,
        button_label: String,
        link: String,
        redirect: String,
        // optional captcha support
        captcha_enabled: { type: Boolean, default: false },
        // function({ username, password, captcha, captcha_id }) returning a promise with an error message or nothing
        submit: Function
    },
    data: function() {
        return {
            username: "",
            password: "",
            password2: "",
            captcha_id: null,
            captcha: "",
            busy: false,
            error_message: ""
        };
    },
    mounted: function() {
        if (this.captcha_enabled) {
            this.update_captcha();
        }
    },
    components: {
        "default-layout": DefaultLayout
    },
    computed: {
        breadcrumb: function() {
            return [START, { name: this.title, icon: this.icon, link: this.link }];
        },
        validInput: function(){
            return this.username != "" && this.password != "" && (this.password == this.password2)
                && (!this.captcha_enabled || this.captcha != "");
        }
    },
    methods: {
        update_captcha: function() {
            signup_captcha().then(c => this.captcha_id = c);
        },
        create_user: function(){
            this.busy = true;
            this.error_message = "";
            this.submit({
                username: this.username,
                password: this.password,
                captcha: this.captcha,
                captcha_id: this.captcha_id
            })
            .then(err_msg => {
                this.busy = false;
                if (err_msg) {
                    this.error_message = err_msg;
                    if (this.captcha_enabled) {
                        this.update_captcha();
                    }
                } else {
                    login(this.username, this.password)
                    .then(() => this.$router.push(this.redirect));
                }
            });
        }
    }
};
</script>

<template>
  <default-layout :icon="icon" :title="title" :breadcrumb="breadcrumb">
      <div class="row">
          <div class="col-md-4"></div>
          <div class="col-md-4 card" style="padding: 20px;">
              <h4>{{heading}}</h4>
              <p>
                  {{description}}
              </p>
              <form @submit.prevent="create_user">
                  <input type="text"
                      class="form-control"
                      placeholder="Username"
                      v-model="username"/>
                  <input type="password"
                      class="form-control"
                      placeholder="Password"
                      v-model="password"/>
                  <input type="password"
                      class="form-control"
                      placeholder="Password (repeat)"
                      v-model="password2"/>
                  <template v-if="captcha_enabled">
                      <img :src="'api/captcha/' + captcha_id + '.png'" v-if="captcha_id"/>
                      <input type="text"
                          class="form-control"
                          placeholder="Captcha"
                          v-model="captcha"/>
                  </template>
                  <button class="btn btn-primary w-100" type="submit" :disabled="!validInput">
                      <i class="fa-solid fa-user"></i>
                      {{button_label}}
                      <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                      <span class="badge bg-danger">{{error_message}}</span>
                  </button>
              </form>
          </div>
      </div>
  </default-layout>
</template>
