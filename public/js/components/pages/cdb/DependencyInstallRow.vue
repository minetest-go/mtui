<script>

export default {
    props: ["dep", "selected_dep"],
    computed: {
        no_candidate: function(){
            return this.dep.choices.length == 0 && !this.dep.installed;
        },
        is_installed: function(){
            return this.dep.installed;
        },
        has_choices: function() {
            return this.dep.choices.length > 0 && !this.dep.installed;
        }
    }
};
</script>

<template>
  <tr v-bind:class="{'table-warning': no_candidate}">
      <td>{{dep.name}}</td>
      <td>
          <select class="form-control" v-on:change="$emit('select_dep', dep.name, $event.target.value)" v-if="has_choices">
              <option v-for="choice in dep.choices" :selected="selected_dep == choice">{{choice}}</option>
          </select>
          <i class="fa fa-spinner fa-spin" v-if="dep.busy"></i>
          <span class="badge bg-danger" v-if="no_candidate">
              <i class="fa-solid fa-triangle-exclamation"></i>
              No installation candidate found!
          </span>
          <span class="badge bg-success" v-if="is_installed">
              <i class="fa fa-check"></i>
              Already installed
          </span>
      </td>
  </tr>
</template>
