<script>
import { store, close } from "../store/confirm.js";

export default {
    data: () => store,
    methods: {
        close,
        keydown: function(e) {
            if (this.show && e.key == "Escape") {
                close(false);
            }
        }
    },
    watch: {
        show: function(v) {
            if (v) {
                // focus the safe option, so a stray "enter" doesn't confirm
                this.$nextTick(() => this.$refs.cancel && this.$refs.cancel.focus());
            }
        }
    },
    mounted: function() {
        document.addEventListener("keydown", this.keydown);
    },
    unmounted: function() {
        document.removeEventListener("keydown", this.keydown);
    }
};
</script>

<template>
  <div v-if="show">
      <div class="modal show d-block" tabindex="-1" role="dialog" aria-modal="true" v-on:click.self="close(false)">
          <div class="modal-dialog modal-dialog-centered">
              <div class="modal-content">
                  <div class="modal-header">
                      <h1 class="modal-title fs-5">
                          <i class="fa-solid fa-triangle-exclamation text-danger" v-if="danger"></i>
                          {{title}}
                      </h1>
                      <button type="button" class="btn-close" aria-label="Close" v-on:click="close(false)"></button>
                  </div>
                  <div class="modal-body" v-if="message">
                      {{message}}
                  </div>
                  <div class="modal-footer">
                      <button type="button" class="btn btn-outline-secondary" ref="cancel" v-on:click="close(false)">
                          Cancel
                      </button>
                      <button type="button" class="btn" :class="danger ? 'btn-danger' : 'btn-primary'" v-on:click="close(true)">
                          <i :class="confirm_icon" v-if="confirm_icon"></i>
                          {{confirm_label}}
                      </button>
                  </div>
              </div>
          </div>
      </div>
      <div class="modal-backdrop fade show"></div>
  </div>
</template>
