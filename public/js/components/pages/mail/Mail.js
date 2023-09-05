import mail_store from '../../../store/mail.js';
import { fetch_mails } from '../../../service/mail.js';
import BoxList from './BoxList.js';
import DefaultLayout from '../../layouts/DefaultLayout.js';
import { START, MAIL } from '../../Breadcrumb.js';

export default {
    data: function() {
        return {
            breadcrumb: [START, MAIL]
        };
    },
    components: {
        "box-list": BoxList,
        "default-layout": DefaultLayout
    },
    computed: {
        busy: () => mail_store.busy,
        boxname: function() {
            return this.$route.params.boxname;
        },
        mails: function() {
            return mail_store[this.boxname];
        }
    },
    methods: {
        refresh: function() {
            mail_store.busy = true;
            fetch_mails()
            .then(() => mail_store.busy = false);
        }
    },
    template: /*html*/`
    <default-layout icon="envelope" title="Mail" :breadcrumb="breadcrumb">
        <div class="row">
            <div class="col-md-8">
                <h4>
                    Mail <small class="text-muted">Overview</small>
                    <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                </h4>
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
        <ul class="nav nav-tabs">
            <li class="nav-item">
                <router-link to="/mail/box/inbox" :class="{'nav-link': true, 'active': boxname == 'inbox'}">
                    Inbox
                </router-link>
            </li>
            <li class="nav-item">
                <router-link to="/mail/box/outbox" :class="{'nav-link': true, 'active': boxname == 'outbox'}">
                    Outbox
                </router-link>
            </li>
        </ul>
        &nbsp;
        <box-list :mails="mails"/>
    </default-layout>
    `
};