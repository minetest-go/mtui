<script>
import DefaultLayout from "../layouts/DefaultLayout.vue";
import { START } from "../../util/breadcrumbs.js";

import { get_mesecon_controls, set_mesecon } from "../../api/mesecons.js";
import MeseconRow from "./MeseconRow.vue";

export default {
    components: {
        "default-layout": DefaultLayout,
        "mesecon-row": MeseconRow
    },
    data: function() {
        return {
            breadcrumb: [START, {
                name: "Mesecons",
                icon: "microchip",
                link: "/mesecons"
            }],
            mesecons: [],
            update_handle: null
        };
    },
    methods: {
        update: function() {
            get_mesecon_controls().then(m => this.mesecons = m);
        },
        swap: function(i1, i2) {
            const m1 = this.mesecons[i1];
            const m2 = this.mesecons[i2];
            if (!m1 || !m2) {
                return;
            }
            const o = m1.order_id;
            m1.order_id = m2.order_id;
            m2.order_id = o;
            Promise.all([set_mesecon(m1), set_mesecon(m2)])
            .then(() => this.update());
        }
    },
    created: function() {
        this.update();
        this.update_handle = setInterval(() => this.update(), 1000);
    },
    unmounted: function() {
        clearInterval(this.update_handle);
    }
};
</script>

<template>
  <default-layout title="Mesecons" icon="microchip" :breadcrumb="breadcrumb">
      <table class="table table-condensed table-striped">
          <thead>
              <tr>
                  <th>Position</th>
                  <th>Last modified</th>
                  <th>Name</th>
                  <th>State</th>
                  <th>Actions</th>
              </tr>
          </thead>
          <tbody>
              <mesecon-row v-for="(mesecon, i) in mesecons" :mesecon="mesecon" :key="mesecon.poskey"
                  v-on:updated="update" v-on:removed="update" v-on:move-up="swap(i, i-1)" v-on:move-down="swap(i, i+1)"/>
          </tbody>
      </table> 
  </default-layout>
</template>
