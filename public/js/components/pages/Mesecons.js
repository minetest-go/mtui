import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

import format_time from "../../util/format_time.js";

import { get_mesecon_controls, set_mesecon, delete_mesecon } from "../../api/mesecons.js";

const node_image_mapping = {
    "mesecons_switch:mesecon_switch_off": "mesecons_switch_off.png",
    "mesecons_switch:mesecon_switch_on": "mesecons_switch_on.png",
};

const colors = ["red", "green", "blue", "gray", "darkgray", "yellow", "orange", "white", "pink", "magenta", "cyan", "violet"];
colors.forEach(c => {
    node_image_mapping[`mesecons_lightstone:lightstone_${c}_off`] = `jeija_lightstone_${c}_off.png`;
    node_image_mapping[`mesecons_lightstone:lightstone_${c}_on`] = `jeija_lightstone_${c}_on.png`;
});

const MeseconRow = {
    props: ["mesecon"],
    data: function() {
        return {
            busy: false,
            error: false
        };
    },
    methods: {
        format_time,
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
            return node_image_mapping[this.mesecon.nodename];
        }
    },
    template: /*html*/`
    <tr>
        <td>{{mesecon.x}},{{mesecon.y}},{{mesecon.z}}</td>
        <td>{{format_time(mesecon.last_modified/1000)}}</td>
        <td>{{mesecon.name}}</td>
        <td>
            <img :src="'pics/' + img_src" v-if="img_src" style="height: 32px; width: 32px; image-rendering: crisp-edges"/>
            <span v-else>{{mesecon.state}}</span>
        </td>
        <td>
            <div class="btn-group">
                <button class="btn btn-success" v-if="mesecon.nodename == 'mesecons_switch:mesecon_switch_off'" v-on:click="set('on')">
                    <i class="fa-solid fa-toggle-on"></i>
                    Set on
                    <i class="fa fa-spinner fa-spin" v-if="busy"></i>
                    <i class="fa-solid fa-xmark" v-if="error"></i>
                </button>
                <button class="btn btn-success" v-if="mesecon.nodename == 'mesecons_switch:mesecon_switch_on'" v-on:click="set('off')">
                    <i class="fa-solid fa-toggle-off"></i>
                    Set off
                    <i class="fa fa-spinner fa-spin" v-if="busy"></i>
                    <i class="fa-solid fa-xmark" v-if="error"></i>
                </button>
                <button class="btn btn-warning" v-on:click="remove">
                    <i class="fa-solid fa-trash"></i>
                    Remove
                </button>
            </div>
        </td>
    </tr>
    `
};

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
            get_mesecon_controls().then(m => {
                m.sort((a,b) => a.poskey > b.poskey);
                this.mesecons = m;
            });
        }
    },
    created: function() {
        this.update();
        this.update_handle = setInterval(() => this.update(), 1000);
    },
    unmounted: function() {
        clearInterval(this.update_handle);
    },
    template: /*html*/`
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
                    <mesecon-row v-for="mesecon in mesecons" :mesecon="mesecon" :key="mesecon.poskey"
                        v-on:updated="update" v-on:removed="update"/>
                </tbody>
            </table> 
        </default-layout>
    `
};