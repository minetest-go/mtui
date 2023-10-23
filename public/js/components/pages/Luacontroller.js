import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

import { get_luacontroller, set_luacontroller } from "../../api/luacontroller.js";

export default {
    props: ["x", "y", "z"],
    components: {
        "default-layout": DefaultLayout
    },
    mounted: function() {
        get_luacontroller({
            pos: {
                x: +this.x,
                y: +this.y,
                z: +this.z
            }
        })
        .then(res => {
            if (res.success) {
                this.code = res.code;
            }
        })
        .then(() => {
            this.cm = CodeMirror.fromTextArea(this.$refs.textarea, {
                lineNumbers: true,
                viewportMargin: Infinity,
                mode: {
                    name: "lua"
                }
            });
        });
    },
    data: function() {
        return {
            breadcrumb: [START, {
                name: "Mesecons",
                icon: "microchip",
                link: "/mesecons"
            }, {
                name: "Luacontroller",
                icon: "microchip",
                link: `/mesecons/luacontroller/${this.x}/${this.y}/${this.z}`
            }],
            code: "",
            errmsg: "",
            success: false,
            busy: false,
            cm: null
        };
    },
    methods: {
        program: function() {
            set_luacontroller({
                pos: {
                    x: +this.x,
                    y: +this.y,
                    z: +this.z
                },
                code: this.cm.getValue()
            })
            .then(res => {
                if (res.success) {
                    this.success = true;
                }
            });
        }
    },
    template: /*html*/`
        <default-layout title="Luacontroller" icon="microchip" :breadcrumb="breadcrumb">
            <div class="row">
                <div class="col-2">
                    <button class="btn btn-success w-100" v-on:click="program">
                        <i class="fa fa-microchip"></i>
                        Program
                        <i class="fa fa-check" v-if="success"></i>
                    </button>
                </div>
            </div>
            <hr>
            <textarea ref="textarea" class="w-100" style="height: 800px;" v-model="code"></textarea>
        </default-layout>
    `
};