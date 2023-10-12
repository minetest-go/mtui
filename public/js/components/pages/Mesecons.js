import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

import { get_mesecon_controls, set_mesecon } from "../../api/mesecons.js";

const MeseconRow = {
    props: ["mesecon"],
    methods: {
        set: function(state) {
            const new_mesecon = Object.assign({}, this.mesecon, {
                state: state
            });
            set_mesecon(new_mesecon).then(() => this.$emit("updated"));
        }
    },
    template: /*html*/`
    <tr>
        <td>{{mesecon.x}},{{mesecon.y}},{{mesecon.z}}</td>
        <td>{{mesecon.last_modified}}</td>
        <td>{{mesecon.name}}</td>
        <td>{{mesecon.state}}</td>
        <td>
            <div class="btn-group">
                <button class="btn btn-success" v-if="mesecon.nodename == 'mesecons_switch:mesecon_switch_off'" v-on:click="set('on')">
                    <i class="fa-solid fa-toggle-on"></i>
                    Set on
                </button>
                <button class="btn btn-success" v-if="mesecon.nodename == 'mesecons_switch:mesecon_switch_on'" v-on:click="set('off')">
                    <i class="fa-solid fa-toggle-off"></i>
                    Set off
                </button>
                <button class="btn btn-warning">
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
            get_mesecon_controls().then(m => this.mesecons = m);
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
                    <mesecon-row v-for="mesecon in mesecons" :mesecon="mesecon" :key="mesecon.poskey" v-on:updated="update"/>
                </tbody>
            </table> 
        </default-layout>
    `
};