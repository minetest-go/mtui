import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

import { get_luacontroller, set_luacontroller, digiline_send } from "../../api/luacontroller.js";

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
            cm: null,
            channel: "",
            message: "",
            digiline_busy: false,
            digiline_success: false
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
        },
        digiline_send: function() {
            this.digiline_busy = true;
            this.digiline_success = false;
            digiline_send({
                pos: {
                    x: +this.x,
                    y: +this.y,
                    z: +this.z
                },
                channel: this.channel,
                msg: this.message
            })
            .then(res => {
                this.digiline_success = res.success;
                this.digiline_busy = false;
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
                <div class="col-2"></div>
                <div class="col-8">
                    <div class="input-group">
                        <input type="text" class="form-control" placeholder="Channel" v-model="channel"/>
                        <input type="text" class="form-control" placeholder="Message" v-model="message"/>
                        <button class="btn btn-secondary" v-on:click="digiline_send">
                            <i class="fa fa-share"></i>
                            Send digiline message
                            <i class="fa fa-check" v-if="digiline_success"></i>
                            <i class="fa fa-spinner fa-spin" v-if="digiline_busy"></i>
                        </button>
                    </div>
                </div>
            </div>
            <hr>
            <textarea ref="textarea" class="w-100" style="height: 800px;" v-model="code"></textarea>
        </default-layout>
    `
};