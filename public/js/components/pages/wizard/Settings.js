
const SettingSuggestion = {
    props: ["type", "name", "title", "description"],
    data: function() {
        return {
            busy: false
        };
    },
    template: /*html*/`
    <div class="card w-100" style="padding: 10px;">
        <div class="row">
            <div class="col-3">
                <h4>{{title}}</h4>
            </div>
            <div class="col-3">
                {{description}}
            </div>
            <div class="col-4">
                
            </div>
            <div class="col-2">
                <button class="btn btn-success w-100">
                    <i class="fa fa-plus"></i>
                    Set stuff
                    <i class="fa fa-spinner fa-spin" v-if="busy"></i>
                </button>
            </div>
        </div>
    </div>
    `
};

export default {
    components: {
        "setting-suggestion": SettingSuggestion
    },
    template: /*html*/`
    <div>
        <h4>Common settings</h4>
        <setting-suggestion name="max_users" type="int" title="Max players" description="Maximum players allowed on the server"/>
        <setting-suggestion name="server_announce" type="bool" title="Announce" description="Announce the server in the public serverlist"/>
    </div>
    `
};