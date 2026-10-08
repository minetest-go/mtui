<script>
import format_time from "../../util/format_time.js";
import { set_mesecon, delete_mesecon } from "../../api/mesecons.js";
import { node_image_mapping, switch_nodes, display_nodes } from "../../util/mesecon_nodes.js";

export default {
    props: ["mesecon"],
    data: function() {
        return {
            busy: false,
            error: false,
            name: this.mesecon.name,
            name_edit: false
        };
    },
    methods: {
        format_time,
        save_name: function() {
            set_mesecon(Object.assign({}, this.mesecon, { name: this.name }))
            .then(() => this.name_edit = false);
        },
        set: function(state) {
            this.busy = true;
            this.error = false;
            const new_mesecon = Object.assign({}, this.mesecon, {
                state: state
            });
            set_mesecon(new_mesecon)
            .then(res => {
                if (res.success) {
                    this.$emit("updated");
                } else {
                    this.error = true;
                }
                this.busy = false;
            });
        },
        remove: function() {
            delete_mesecon(this.mesecon.poskey)
            .then(() => this.$emit("removed"));
        }
    },
    computed: {
        img_src: function() {
            if (this.is_mooncontroller) {
                return "mooncontroller_top.png";
            }
            if (this.is_luacontroller) {
                return "jeija_luacontroller_top.png";
            }
            return node_image_mapping[this.mesecon.nodename];
        },
        is_switch: function() {
            return switch_nodes[this.mesecon.nodename];
        },
        is_display: function() {
            return display_nodes[this.mesecon.nodename];
        },
        is_luacontroller: function() {
            return this.mesecon.nodename.startsWith("mesecons_luacontroller");
        },
        is_mooncontroller: function() {
            return this.mesecon.nodename.startsWith("mooncontroller");
        }
    }
};
</script>

<template>
  <tr>
      <td>{{mesecon.x}},{{mesecon.y}},{{mesecon.z}}</td>
      <td>{{format_time(mesecon.last_modified/1000)}}</td>
      <td>
          <div class="input-group" v-if="name_edit">
              <input type="text" class="form-control" v-model="name"/>
              <button class="btn btn-success" v-on:click="save_name">
                  <i class="fa fa-floppy-disk"></i>
              </button>
          </div>
          <div v-else>
              <span class="badge bg-primary">
                  {{name}}
              </span>
              <i class="fa-regular fa-pen-to-square" v-on:click="name_edit = true"></i>
          </div>
      </td>
      <td>
          <img :src="'pics/' + img_src" v-if="img_src" style="height: 32px; width: 32px; image-rendering: crisp-edges"/>
          <span v-if="!img_src || is_display" class="badge bg-secondary">{{mesecon.state}}</span>
          <i class="fa fa-spinner fa-spin" v-if="busy"></i>
          <i class="fa-solid fa-xmark" v-else-if="error"></i>
      </td>
      <td>
          <div class="btn-group">
              <button class="btn btn-secondary" v-on:click="$emit('move-down')">
                  <i class="fa-solid fa-chevron-down"></i>
              </button>
              <button class="btn btn-secondary" v-on:click="$emit('move-up')">
                  <i class="fa-solid fa-chevron-up"></i>
              </button>
              <button class="btn btn-warning" v-on:click="remove">
                  <i class="fa-solid fa-trash"></i>
                  Remove
              </button>
              <router-link :to="'/mesecons/luacontroller/' + mesecon.x + '/' + mesecon.y + '/' + mesecon.z"
                  class="btn btn-primary" v-if="is_luacontroller || is_mooncontroller">
                  <i class="fa fa-microchip"></i>
                  Program
              </router-link>
              <button class="btn btn-success" v-if="is_switch" :disabled="busy" v-on:click="set('on')">
                  <i class="fa-solid fa-toggle-on"></i>
                  Set ON
              </button>
              <button class="btn btn-success" v-if="is_switch" :disabled="busy" v-on:click="set('off')">
                  <i class="fa-solid fa-toggle-off"></i>
                  Set OFF
              </button>
          </div>
      </td>
  </tr>
</template>
