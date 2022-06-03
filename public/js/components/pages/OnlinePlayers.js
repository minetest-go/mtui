import onlineplayers_store from "../../store/onlineplayers.js";
import format_seconds from "../../util/format_seconds.js";

export default {
    data: () => onlineplayers_store,
    methods: {
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
                    <th>Address</th>
                    <th>Protocol-Version</th>
                    <th>Connected since</th>
                    <th>RTT (min/avg/max)</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="player in players" :key="player.name">
                    <td>
                        <router-link :to="'/player/' + player.name">
                            {{player.name}}
                        </router-link>
                    </td>
                    <td>
                        {{player.address}}
                        <span class="badge bg-info">IPv{{player.ip_version}}</span>
                    </td>
                    <td>{{player.protocol_version}}</td>
                    <td>{{ format_seconds(player.connection_uptime) }}</td>
                    <td>
                        <i class="fa-solid fa-signal" v-bind:style="{'color': signal_color(player.avg_rtt) }"></i>
                        {{ Math.floor(player.min_rtt*1000)/1000 }}/
                        {{ Math.floor(player.avg_rtt*1000)/1000 }}/
                        {{ Math.floor(player.max_rtt*1000)/1000 }} s
                    </td>
                </tr>
            </tbody>
        </table>
    `
};