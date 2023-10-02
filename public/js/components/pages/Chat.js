import DefaultLayout from "../layouts/DefaultLayout.js";

import { START } from "../Breadcrumb.js";
import { get_latest_chat_messages, send_message } from "../../api/chat.js";
import { execute_chatcommand } from "../../api/chatcommand.js";
import { get_claims } from "../../service/login.js";

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            history: "",
            msg: "",
            handle: null,
            channel: "main",
            breadcrumb: [START, {
                name: "Chat",
                icon: "comment",
                link: ""
            }]
        };
    },
    mounted: function() {
        this.update();
        this.handle = setInterval(() => this.update(), 1000);
    },
    unmounted: function() {
        clearInterval(this.handle);
    },
    methods: {
        send: function() {
            if (this.is_command) {
                // send command
                execute_chatcommand(get_claims().username, this.msg.substring(1))
                .then(result => {
                    if (result.success) {
                        this.history += `${result.message}\n`;
                    } else {
                        this.history += `${result.message}\n`;
                    }
                    this.msg = "";
                });

            } else {
                // send message
                send_message({
                    channel: this.channel,
                    message: this.msg
                })
                .then(() => {
                    this.update();
                    this.msg = "";
                });
            }
        },
        update: function() {
            const chat_el = this.$refs.chat_pre;
            get_latest_chat_messages(this.channel)
            .then(msgs => {
                let history = "";
                msgs.forEach(msg => history += `<${msg.name}> ${msg.message}\n`);
                const has_scrolled_down = chat_el.scrollTop == chat_el.scrollTopMax;

                this.history = history;
                if (has_scrolled_down) {
                    // scroll to the bottom
                    setTimeout(() => chat_el.scrollTop = chat_el.scrollTopMax, 10);
                }
            });
        }
    },
    computed: {
        is_command: function() {
            return this.msg.length > 0 && this.msg[0] == '/';
        }
    },
    template: /*html*/`
        <default-layout title="Chat" icon="comment" :breadcrumb="breadcrumb">
            <div class="row">
                <div class="col-12">
                    <pre ref="chat_pre" class="w-100" style="height: 400px; background-color: grey;">{{history}}</pre>
                </div>
            </div>
            <form @submit.prevent="send" class="row">
                <div class="input-group">
                    <input type="text" placeholder="Message" v-model="msg" class="form-control"/>
                    <button class="btn btn-success" type="submit" v-if="!is_command">
                        <i class="fa-solid fa-paper-plane"></i>
                        Send
                    </button>
                    <button class="btn btn-warning" type="submit" v-if="is_command">
                        <i class="fa-solid fa-terminal"></i>
                        Execute command
                    </button>
                </div>
            </form>
        </default-layout>
    `
};