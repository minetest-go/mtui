import onlineplayers_store from "../../store/onlineplayers.js";
import format_seconds from "../../util/format_seconds.js";
import { has_priv } from "../../service/login.js";

export default {
    data: () => onlineplayers_store,
    methods: {
        has_priv: has_priv,
        format_seconds: format_seconds,
        signal_color: function(rtt) {
            if (rtt < 0.1) return "green";
            if (rtt < 0.5) return "orange";
            return "red";
        }
    },
    template: /*html*/`
        <table class="table table-condensed table-striped">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Health</th>
                    <th v-if="has_priv('ban')">Position</th>
                    <th v-if="has_priv('ban')">Address</th>
                    <th v-if="has_priv('ban')">Protocol-Version</th>
                    <th v-if="has_priv('ban')">Connected since</th>
                    <th v-if="has_priv('ban')">RTT (min/avg/max)</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="player in players" :key="player.name">
                    <td>
                        <router-link :to="'/profile/' + player.name">
                            {{player.name}}
                        </router-link>
                    </td>
                    <td>
                        <i class="fa-solid fa-heart" style="color: red;"></i>
                        {{ player.hp }}
                    </td>
                    <td v-if="has_priv('ban')">
                        {{Math.floor(player.pos.x)}}/{{Math.floor(player.pos.y)}}/{{Math.floor(player.pos.z)}}
                    </td>
                    <td v-if="has_priv('ban')">
                        {{player.info.address}}
                    </td>
                    <td v-if="has_priv('ban')">
                        {{player.info.protocol_version}}
                    </td>
                    <td v-if="has_priv('ban')">
                        {{ format_seconds(player.info.connection_uptime) }}
                    </td>
                    <td v-if="has_priv('ban')">
                        <i class="fa-solid fa-signal" v-bind:style="{'color': signal_color(player.info.avg_rtt) }"></i>
                        {{ Math.floor(player.info.min_rtt*1000)/1000 }}/
                        {{ Math.floor(player.info.avg_rtt*1000)/1000 }}/
                        {{ Math.floor(player.info.max_rtt*1000)/1000 }} s
                    </td>
                </tr>
            </tbody>
        </table>
    `
};