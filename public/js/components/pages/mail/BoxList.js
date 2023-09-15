import { fetch_mails } from '../../../service/mail.js';
import { remove } from '../../../api/mail.js';
import format_time from '../../../util/format_time.js';

const store = Vue.reactive({
    sortfield: "time",
    sortdirection: "asc",
    sorted_mails: []
});

export default {
    props: ["mails"],
    data: () => store,
    watch: {
        "mails": function() {
            this.sort();
        }
    },
    methods: {
        format_time: format_time,
        delete_mail: function(msg){
            remove(msg)
            .then(() => fetch_mails())
            .then(() => this.sort(this.sortfield, this.sortdirection));
        },
        sort: function(field, direction) {
            this.sortfield = field || this.sortfield;
            this.sortdirection = direction || this.sortdirection;

            if (this.sortdirection == "asc"){
                this.sorted_mails = this.mails.sort((a, b) => a[this.sortfield] < b[this.sortfield]);
            } else if (this.sortdirection == "desc") {
                this.sorted_mails = this.mails.sort((a, b) => a[this.sortfield] > b[this.sortfield]);
            }
        },
        chevron_style: function(field, direction) {
            if (this.sortfield == field && this.sortdirection == direction) {
                return { color: "red" };
            } else {
                return {};
            }
        }
    },
    template: /*html*/`
    <div>
        <div class="alert alert-primary" v-if="sorted_mails.length == 0">
            No mails
        </div>
        <table class="table table-condensed" v-if="sorted_mails.length > 0">
            <thead>
                <tr>
                    <th>
                        From
                        <i class="fa-solid fa-chevron-up" v-bind:style="chevron_style('from', 'asc')" v-on:click="sort('from', 'asc')"></i>
                        <i class="fa-solid fa-chevron-down" v-bind:style="chevron_style('from', 'desc')" v-on:click="sort('from', 'desc')"></i>
                    </th>
                    <th>
                        Subject
                        <i class="fa-solid fa-chevron-up" v-bind:style="chevron_style('subject', 'asc')" v-on:click="sort('subject', 'asc')"></i>
                        <i class="fa-solid fa-chevron-down" v-bind:style="chevron_style('subject', 'desc')" v-on:click="sort('subject', 'desc')"></i>
                    </th>
                    <th>
                        Time
                        <i class="fa-solid fa-chevron-up" v-bind:style="chevron_style('time', 'asc')" v-on:click="sort('time', 'asc')"></i>
                        <i class="fa-solid fa-chevron-down" v-bind:style="chevron_style('time', 'desc')" v-on:click="sort('time', 'desc')"></i>
                    </th>
                    <th>Action</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(mail, index) in sorted_mails" :key="index" :class="{'table-info': !mail.read}">
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