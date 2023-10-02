import DefaultLayout from "../layouts/DefaultLayout.js";

import { START } from "../Breadcrumb.js";
import { search_messages, send_message } from "../../api/chat.js";
import { execute_chatcommand } from "../../api/chatcommand.js";
import { get_claims } from "../../service/login.js";
import format_time from "../../util/format_time.js";

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: function() {
        return {
            history: [],
            msg: "",
            handle: null,
            channel: "main",
            from_timestamp: Date.now() - (3600*1000*24*7), //7 days back or 1000 messages
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
        format_time,
        send: function() {
            if (this.is_command) {
                // send command
                execute_chatcommand(get_claims().username, this.msg.substring(1))
                .then(result => {
                    this.scroll();
                    this.history.push({
                        timestamp: +Date.now(),
                        name: "",
                        message: result.message,
                        success: result.success
                    });
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
            const later = Date.now() + (3600*1000);
            search_messages(this.channel, this.from_timestamp, later)
            .then(msgs => {
                msgs.forEach(msg => {
                    if (msg.timestamp > this.from_timestamp) {
                        this.from_timestamp = msg.timestamp;
                    }
                    this.scroll();
                    this.history.push(msg);
                });
            });
        },
        scroll: function() {
            const el = this.$refs.container;
            if (el.scrollTop == el.scrollTopMax) {
                // at the bottom, scroll further
                setTimeout(() => el.scrollTop = el.scrollTopMax, 10);
            }
        }
    },
    computed: {
        is_command: function() {
            return this.msg.length > 0 && this.msg[0] == '/';
        }
    },
    template: /*html*/`
        <default-layout title="Chat" icon="comment" :breadcrumb="breadcrumb">
            <div ref="container" style="height: 600px; overflow: scroll;">
                <div v-for="msg in history" :key="msg.id"
                    v-bind:class="{'bg-success':msg.success==true, 'bg-warning':msg.success==false}"
                    style="display: flex;">
                    <div class="text-muted" style="width: 200px; flex: 0 0 auto;">
                        {{format_time(msg.timestamp/1000)}}
                    </div>
                    <div style="width: 200px; flex: 0 0 auto;">
                        <router-link :to="'/profile/' + msg.name" v-if="msg.name != ''">
                            {{msg.name}}
                        </router-link>
                        <span v-else>
                            {{msg.name}}
                        </span>
                    </div>
                    <div style="flex: 1 1 auto;">
                        {{msg.message}}
                    </div>
                </div>
            </div>
            <form @submit.prevent="send" class="row">
                <div class="input-group">
                    <input type="text" placeholder="Message" v-model="msg" class="form-control"/>
                    <button class="btn btn-success" type="submit" v-if="!is_command" :disabled="msg == ''">
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