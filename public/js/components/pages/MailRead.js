import format_time from '../../util/format_time.js';
import mail_store from '../../store/mail.js';

export default {
    props: ["sender", "time"],
    computed: {
        mail: function(){
            return mail_store.mails.find(m => m.sender == this.sender && m.time == this.time);
        }
    },
    methods: {
        format_time: format_time
    },
    data: function() {
        const sender = this.$route.params.sender;
        const time = +this.$route.params.time;
        return {
            sender: sender,
            time: time
        };
    },
    template: /*html*/`
    <div v-if="mail">
        <h3>
            Mail from
            <small class="text-muted">
                {{mail.sender}}
            </small>
        </h3>
        Sent: {{format_time(mail.time)}}
        <hr>
        <pre>{{mail.body}}</pre>
    </div>
    `
};