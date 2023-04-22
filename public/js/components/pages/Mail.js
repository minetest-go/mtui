import mail_store from '../../store/mail.js';
import { fetch_mails } from '../../service/mail.js';
import { remove } from '../../api/mail.js';
import format_time from '../../util/format_time.js';

export default {
    data: function() {
        return {
            mail_store: mail_store,
            busy: false
        };
    },
    methods: {
        refresh: function(){
            this.busy = true;
            fetch_mails()
            .then(() => this.busy = false);
        },
        format_time: format_time,
        delete_mail: function(msg){
            this.busy = true;
            remove(msg)
            .then(() => fetch_mails())
            .then(() => this.busy = false);
        }
    },
    template: /*html*/`
    <div>
        <div class="row">
            <div class="col-md-8">
                <h3>
                    Mail <small class="text-muted">Overview</small>
                    <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                </h3>
            </div>
            <div class="col-md-4 btn-group">
                <router-link to="/mail/compose" class="btn btn-primary">
                    <i class="fa-solid fa-pen-to-square"></i>
                    Compose
                </router-link>
                <a class="btn btn-success" v-on:click="refresh">
                    <i class="fa-solid fa-rotate"></i>
                    Refresh
                </a>
            </div>
        </div>
        &nbsp;
        <div class="alert alert-primary" v-if="mail_store.mails.length == 0">
            No mails
        </div>
        <table class="table table-condensed" v-if="mail_store.mails.length > 0">
            <thead>
                <tr>
                    <th>From</th>
                    <th>Subject</th>
                    <th>Time</th>
                    <th>Action</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(mail, index) in mail_store.mails" :key="index" :class="{'table-info': !mail.read}">
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