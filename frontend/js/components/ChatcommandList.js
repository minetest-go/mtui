import { get_chatcommand_infos } from "../api/uimod.js";
import { get_privs } from "../service/login.js";

export default {
    emits: ["selected"],
    data: function() {
        return {
            chatcommands: {},
            privs: get_privs()
        };
    },
    mounted: async function() {
        this.chatcommands = await get_chatcommand_infos();
    },
    computed: {
        available_chatcommands: function() {
            const cmds = {};
            const privs = get_privs();

            Object.keys(this.chatcommands).forEach(name => {
                const cc = this.chatcommands[name];
                const available = !cc.privs || privs.find(p => cc.privs[p]);
                if (available) {
                    cmds[name] = cc;
                }
            });

            return cmds;
        }
    },
    template: /*html*/`
    <h4>Available commands</h4>
    <table class="table table-condensed table-striped">
        <thead>
            <tr>
                <th>Name</th>
                <th>Description</th>
                <th>Params</th>
                <th>Required privileges</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="(def, name) in available_chatcommands" :key="name">
                <td>
                    <a class="btn btn-sm btn-primary" v-on:click="$emit('selected', name)">
                        {{name}}
                    </a>
                </td>
                <td>{{def.description}}</td>
                <td>{{def.params}}</td>
                <td>
                    <span class="badge bg-secondary" v-for="(_, priv) in def.privs">
                        {{priv}}
                    </span>
                </td>
            </tr>
        </tbody>
    </table>
    `
};
