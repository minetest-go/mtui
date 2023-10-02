import DefaultLayout from "../layouts/DefaultLayout.js";

import { START } from "../Breadcrumb.js";
import { get_latest_chat_messages, send_message } from "../../api/chat.js";

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
    created: function() {
        this.update();
        this.handle = setInterval(() => this.update(), 1000);
    },
    unmounted: function() {
        clearInterval(this.handle);
    },
    methods: {
        send: function() {
            send_message({
                channel: this.channel,
                message: this.msg
            })
            .then(() => {
                this.update();
                this.msg = "";
            });
        },
        update: function() {
            get_latest_chat_messages(this.channel)
            .then(msgs => {
                let history = "";
                msgs.forEach(msg => history += `<${msg.name}> ${msg.message}\n`);
                this.history = history;
                this.$refs.chat_pre.scrollTop = this.$refs.chat_pre.scrollTopMax;
            });
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
                    <button class="btn btn-success" type="submit">
                        <i class="fa-solid fa-paper-plane"></i>
                        Send
                    </button>
                </div>
            </form>
        </default-layout>
    `
};