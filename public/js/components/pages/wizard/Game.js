import { get_mods_by_type } from "../../../service/mods.js";

export default {
    computed: {
        game: function() {
            const list = get_mods_by_type("game");
            return list.length == 1 ? list[0] : null;
        },
        complete: function() {
            return this.game;
        }
    },
    template: /*html*/`
    <div>
        <div class="row" v-if="complete">
            <div class="col-12">
                <div class="alert alert-success">
                    <i class="fa fa-check"></i> Game installed: <b>{{game.name}}</b>
                </div>
            </div>
        </div>
        <div v-if="!complete">
            Select the game to install
        </div>
    </div>
    `
};