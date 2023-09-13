import format_time from '../../../util/format_time.js';
import { get_mail } from '../../../store/mail.js';
import { mark_read, remove } from '../../../api/mail.js';
import { store as mail_compose_store } from './Compose.js';
import { fetch_mails } from '../../../service/mail.js';
import { get_claims } from '../../../service/login.js';
import DefaultLayout from '../../layouts/DefaultLayout.js';
import { START, MAIL } from '../../Breadcrumb.js';

export default {
    props: ["id"],
    components: {
        "default-layout": DefaultLayout
    },
    computed: {
        mail: function() {
            const mail = get_mail(this.id);
            if (mail && !mail.read) {
                mark_read(mail)
                .then(() => mail.read = true);
            }
            
            return mail;
        },
        is_sent: function() {
            return get_claims().username == this.mail.from;
        }
    },
    methods: {
        format_time: format_time,
        reply: function() {
            mail_compose_store.recipients = [this.mail.from];
            mail_compose_store.subject = "Re: " + this.mail.subject;
            mail_compose_store.body = "\n---- Original message ----\n" + this.mail.body;
            this.$router.push({ path:"/mail/compose" });
        },
        remove: function() {
            remove(this.mail)
            .then(() => fetch_mails())
            .then(() => this.$router.push({ path:"/mail" }));
        }
    },
    data: function() {
        return {
            breadcrumb: [START, MAIL, {
                name: `Read mail`,
                icon: "envelope-open",
                link: `/mail/read/${this.id}`
            }]
        };
    },
    template: /*html*/`
    <default-layout v-if="mail" icon="envelope-open" title="Read mail" :breadcrumb="breadcrumb">
        <div class="row">
            <div class="col-md-10">
                <h4 v-if="is_sent">
                    Mail sent to
                    <small class="text-muted">
                        {{mail.to}}
                    </small>
                </h4>
                <h4 v-else>
                    Mail from
                    <small class="text-muted">
                        {{mail.from}}
                    </small>
                </h4>
            </div>
            <div class="col-md-2 btn-group">
                <a v-on:click="reply" class="btn btn-primary">
                    <i class="fa-solid fa-pen-to-square"></i>
                    Reply
                </a>
                <a v-on:click="remove" class="btn btn-danger">
                    <i class="fa-solid fa-trash-can"></i>
                    Delete
                </a>
            </div>
        </div>
        Sent: <b>{{format_time(mail.time)}}</b>
        <br>
        Subject: <b>{{mail.subject}}</b>
        <hr>
        <pre>{{mail.body}}</pre>
    </default-layout>
    `
};