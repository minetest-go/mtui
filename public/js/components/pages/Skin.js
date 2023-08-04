
import SkinSlot from "../SkinSlot.js";

export default {
    components: {
        "skin-slot": SkinSlot
    },
    template: /*html*/`
        <div>
            <h3>
                Skin
                <small class="text-muted">manager</small>
            </h3>
            <div class="row">
                <div class="col-md-2">
                    <skin-slot skin_id="1"/>
                </div>
                <div class="col-md-2">
                    <skin-slot skin_id="2"/>
                </div>
                <div class="col-md-2">
                    <skin-slot skin_id="3"/>
                </div>
                <div class="col-md-2">
                    <skin-slot skin_id="4"/>
                </div>
                <div class="col-md-2">
                    <skin-slot skin_id="5"/>
                </div>
                <div class="col-md-2">
                    <skin-slot skin_id="6"/>
                </div>
            </div>
            <hr>
            <div class="row">
                <div class="col-12">
                    <div class="alert alert-info">
                        <i class="fa fa-info-circle"></i>
                        <b>Note:</b> New skins in empty slots need a server-restart to become available ingame
                    </div>
                </div>
            </div>
        </div>
    `
};