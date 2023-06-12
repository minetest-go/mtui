import mail_store from '../../../store/mail.js';
import { fetch_mails } from '../../../service/mail.js';
import { remove } from '../../../api/mail.js';
import format_time from '../../../util/format_time.js';


export default {
    props: ["mails"],
    methods: {
        format_time: format_time,
        delete_mail: function(msg){
            mail_store.busy = true;
            remove(msg)
            .then(() => fetch_mails())
            .then(() => mail_store.busy = false);
        }
    },
    template: /*html*/`
    <div>
        <div class="alert alert-primary" v-if="mails.length == 0">
            No mails
        </div>
        <table class="table table-condensed" v-if="mails.length > 0">
            <thead>
                <tr>
                    <th>From</th>
                    <th>Subject</th>
                    <th>Time</th>
                    <th>Action</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(mail, index) in mails" :key="index" :class="{'table-info': !mail.read}">
                    <td>
                        <router-link :to="'/profile/' + mail.from">
                            {{mail.from}}
                        </router-link>                    
                    </td>
                    <td>
                        <router-link :to="'/mail/read/' + mail.id">
                            {{mail.subject}}
                        </router-link>
                    </td>
                    <td>{{format_time(mail.time)}}</td>
                    <td>
                        <div class="btn-group">
                            <router-link class="btn btn-primary" :to="'/mail/read/' + mail.id">
                                <i class="fa-solid fa-envelope-open"></i>
                                Open
                            </router-link>
                            <a class="btn btn-danger" v-on:click="delete_mail(mail)">
                                <i class="fa-solid fa-trash-can"></i>
                                Delete
                            </a>
                        </div>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
    `
};