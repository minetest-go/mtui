<script>
import { store, dismiss } from "../store/toast.js";

const icons = {
    success: "fa-solid fa-circle-check text-success",
    danger: "fa-solid fa-triangle-exclamation text-danger",
    warning: "fa-solid fa-circle-exclamation text-warning",
    info: "fa-solid fa-circle-info text-info"
};

export default {
    data: function() {
        return {
            store,
            expanded: {}
        };
    },
    methods: {
        dismiss,
        icon: type => icons[type] || icons.info
    }
};
</script>

<template>
  <div class="toast-container app-toasts position-fixed end-0 p-3" aria-live="polite">
      <div
          class="toast show align-items-center"
          :class="'border-' + toast.type"
          role="status"
          v-for="toast in store.list"
          :key="toast.id">
          <div class="d-flex align-items-start p-2">
              <i :class="icon(toast.type)" class="mt-1 me-2"></i>
              <div class="me-auto text-break">
                  <div class="fw-semibold" v-if="toast.title">
                      {{toast.title}}
                      <span class="badge text-bg-secondary" v-if="toast.count > 1">{{toast.count}}x</span>
                  </div>
                  <div>{{toast.message}}</div>
                  <div v-if="toast.details">
                      <button type="button" class="btn btn-link btn-sm p-0" v-on:click="expanded[toast.id] = !expanded[toast.id]">
                          {{ expanded[toast.id] ? "Hide details" : "Details" }}
                      </button>
                      <pre class="small mb-0 mt-1" style="white-space: pre-wrap;" v-if="expanded[toast.id]">{{toast.details}}</pre>
                  </div>
              </div>
              <button type="button" class="btn-close ms-2" aria-label="Close" v-on:click="dismiss(toast.id)"></button>
          </div>
      </div>
  </div>
</template>

<style scoped>
/* top right, below the top bar: keeps clear of sticky bars at the bottom (e.g. the settings save bar) */
.app-toasts {
    top: 3.5rem;
}
</style>
