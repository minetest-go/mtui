
export default {
    props: ["status"],
    template: /*html*/`
    <i class="fa fa-play" v-if="status && status.running" style="color: green;"></i>
    <i class="fa fa-stop" v-if="status && !status.running" style="color: red;"></i>
    `
};