<script>
export default {
    props: ["items"],
    methods: {
        get_icon_class: function(item) {
            const cl = {};
            if (item.icon){
                cl.fa = true;
                cl["fa-" + item.icon] = true;
            }
            return cl;
        }
    }
};
</script>

<template>
  <nav>
      <ol class="breadcrumb text-bg-secondary" style="padding: 5px; border-radius: 5px;">
          <li class="breadcrumb-item" v-for="item in items">
              <router-link class="link-light" :to="item.link" v-if="item.link">
                  <i v-bind:class="get_icon_class(item)" v-if="item.icon"></i>
                  {{item.name}}
              </router-link>
              <span v-else>
                  <i v-bind:class="get_icon_class(item)" v-if="item.icon"></i>
                  {{item.name}}
              </span>
          </li>
      </ol>
  </nav>
</template>
