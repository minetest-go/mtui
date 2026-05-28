import DefaultLayout from "../layouts/DefaultLayout.js";
import CodeEditor from "../CodeEditor.js";

import { START } from "../Breadcrumb.js";

import { get_luacontroller, set_luacontroller, digiline_send } from "../../api/luacontroller.js";

export default {
    props: ["x", "y", "z"],
    components: {
        "default-layout": DefaultLayout,
        "code-editor": CodeEditor
    },
    mounted: async function() {
        const res = await get_luacontroller({
            pos: {
                x: +this.x,
                y: +this.y,
                z: +this.z
            }
        });
        if (res.success) {
            this.code = res.code;
        }
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
            channel: "",
            message: "",
            digiline_busy: false,
            digiline_success: false
        };
    },
    methods: {
        program: async function() {
            this.busy = true;
            const res = await set_luacontroller({
                pos: {
                    x: +this.x,
                    y: +this.y,
                    z: +this.z
                },
                code: this.code
            });
            if (res.success) {
                this.success = true;
            }
            this.busy = false;
        },
        digiline_send: async function() {
            this.digiline_busy = true;
            this.digiline_success = false;
            const res = await digiline_send({
                pos: {
                    x: +this.x,
                    y: +this.y,
                    z: +this.z
                },
                channel: this.channel,
                msg: this.message
            });
            this.digiline_success = res.success;
            this.digiline_busy = false;
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
                        <i class="fa fa-spinner fa-spin" v-if="busy"></i>
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
            <code-editor v-model="code" class="w-100" style="height: 800px;" mode="lua"/>
        </default-layout>
    `
};