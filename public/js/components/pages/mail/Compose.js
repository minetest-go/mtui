import { check_recipient, send } from "../../../api/mail.js";
import { fetch_mails } from "../../../service/mail.js";
import DefaultLayout from '../../layouts/DefaultLayout.js';
import { START, MAIL } from '../../Breadcrumb.js';

export const store = Vue.reactive({
    body: "",
    subject: "",
    add_recipient_name: "",
    invalid_username: false,
    recipients: [],
    busy: false,
    mail_sent: false,
    breadcrumb: [START, MAIL, {
        name: "Compose mail",
        icon: "pen-to-square",
        link: "/mail/compose"
    }]
});

export default {
    components: {
        "default-layout": DefaultLayout
    },
    data: () => store,
    methods: {
        remove_recipient: function(name){
            this.recipients = this.recipients.filter(r => r != name);
        },
        add_recipient: function() {
            if (this.recipients.find(r => r == this.add_recipient_name)) {
                // already exists in list
                this.add_recipient_name = "";
                return;
            }

            // check and add
            this.busy = true;
            check_recipient(this.add_recipient_name)
            .then(err => {
                this.busy = false;
                if (err) {
                    this.invalid_username = true;
                } else {
                    this.invalid_username = false;
                    this.recipients.push(this.add_recipient_name);
                    this.add_recipient_name = "";
                }
            });
        },
        send: function() {
            this.busy = true;
            this.mail_sent = false;

            const mail = {
                body: this.body,
                subject: this.subject,
                to: this.recipients[0]
            };

            if (this.recipients.length > 1) {
                mail.cc = this.recipients.slice(1).join(",");
            }

            send(mail)
            .then(() => {
                this.body = "";
                this.subject = "";
                this.busy = false;
                this.mail_sent = true;

                // re-read mails
                fetch_mails();
            });
            // TODO: error-handling
        }
    },
    template: /*html*/`
    <default-layout icon="paper-plane" title="Compose mail" :breadcrumb="breadcrumb">
        <h4>
            Mail <small class="text-muted">Compose</small>
            <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
        </h4>
        <form @submit.prevent="add_recipient" class="row">
            <div class="col-md-4">
                <input type="text" placeholder="Recipient" v-model="add_recipient_name" :class="{'form-control':true,'is-invalid':invalid_username}"/>
                <div class="invalid-feedback" v-if="invalid_username">
                    The player does not exist
                </div>
            </div>
            <div class="col-md-2">
                <button type="button" class="btn btn-outline-secondary" v-on:click="add_recipient" :disabled="add_recipient_name == ''">
                    <i class="fa-solid fa-plus"></i>
                    Add
                </button>
            </div>
            <div class="col-md-6">
                <h3>
                    <span class="badge bg-secondary" v-for="recipient in recipients" style="margin-left: 3px;">
                        {{recipient}}
                        <i class="fa-solid fa-trash-can" v-on:click="remove_recipient(recipient)"></i>
                    </span>
                </h3>
            </div>
        </form>
        <hr/>
        <div class="row">
            <div class="col-md-12">
                <input type="text" class="form-control" v-model="subject" placeholder="Subject"/>
            </div>
        </div>
        <div class="row">
            <div class="col-md-12">
                <textarea v-model="body" class="form-control" placeholder="Mail text" style="height: 250px;">
                </textarea>
                &nbsp;
                <div class="alert alert-warning" v-if="recipients.length == 0">
                    Add a recipient first with the <mark>Add</mark> button
                </div>
                <button class="btn btn-outline-primary w-100" :disabled="body == '' || recipients.length == 0 || subject == ''" v-on:click="send">
                    <i class="fa-solid fa-paper-plane"></i>
                    Send
                    &nbsp;
                    <span v-if="mail_sent">
                        <i class="fa-solid fa-check" style="color: green;"></i>
                        Mail sent!
                    </span>
                </button>
            </div>
        </div>
    </default-layout>
    `
};