import mail_store from '../../store/mail.js';
import { fetch_mails } from '../../service/mail.js';
import format_time from '../../util/format_time.js';

export default {
    data: function() {
        return {
            mail_store: mail_store
        };
    },
    methods: {
        fetch_mails: fetch_mails,
        format_time: format_time
    },
    template: /*html*/`
    <div>
        <h3>Mail</h3>
        <table class="table table-condensed">
            <thead>
                <tr>
                    <th>From</th>
                    <th>Subject</th>
                    <th>Time</th>
                    <th>Unread</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(mail, index) in mail_store.mails" :key="index">
                    <td>{{mail.sender}}</td>
                    <td>
                        <router-link :to="'/mail/read/' + mail.sender + '/' + mail.time">
                            {{mail.subject}}
                        </router-link>
                    </td>
                    <td>{{format_time(mail.time)}}</td>
                    <td>{{mail.unread}}</td>
                </tr>
            </tbody>
        </table>
    </div>
    `
};