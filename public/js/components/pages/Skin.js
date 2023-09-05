import SkinSlot from "../SkinSlot.js";
import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

export default {
    data: function() {
        return {
            breadcrumb: [START, {
                icon: "user-astronaut",
                name: "Skin",
                link: "/skin"
            }]
        };
    },
    components: {
        "skin-slot": SkinSlot,
        "default-layout": DefaultLayout
    },
    template: /*html*/`
        <default-layout icon="user-astronaut" title="Skin" :breadcrumb="breadcrumb">
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
        </default-layout>
    `
};