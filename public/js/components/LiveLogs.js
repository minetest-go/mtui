import events from '../events.js';

export default {
    props: ["name"],
    data: function(){
        return {
            txt: ""
        };
    },
    created: function(){
        /*
        const from = Math.floor(Date.now() / 1000) - 30;
        const to = Math.floor(Date.now() / 1000) + 30;

        get(this.name, from, to)
        .then(txt => this.txt = txt);
        */
        events.on("log", this.listener);
    },
    beforeUnmount: function(){
        events.off("log", this.listener);
    },
    methods: {
        listener: function(data){
            if (data.key != this.name)
                return;

            const newtxt = this.txt + data.message + '\n';
            const lines = newtxt.split("\n");

            while (lines.length > 20){
                lines.shift();
            }

            this.txt = lines.join('\n');
        }
    },
    template: /*html*/`
        <pre style="background-color: lightgrey; border-radius: 5px;">{{ txt }}</pre>
    `
};