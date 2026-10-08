<script>
import format_count from "../../util/format_count.js";

export default {
    props: ["instrument", "name"],
    computed: {
        avg_micros: function() {
            return Math.floor(this.instrument.time_all / this.instrument.samples);
        },
        tableClass: function() {
            if (this.avg_micros > 50000) {
                return { "table-danger": true };
            } else if (this.avg_micros > 25000) {
                return { "table-warning": true };
            }
        }
    },
    methods: {
        format_count
    }
};
</script>

<template>
  <tr v-bind:class="tableClass">
      <td>
          <span class="badge bg-secondary">{{name}}</span>
      </td>
      <td>{{instrument.time_min}}</td>
      <td>{{avg_micros}}</td>
      <td>{{instrument.time_max}}</td>
      <td>{{ format_count(instrument.samples) }}</td>
  </tr>
</template>
