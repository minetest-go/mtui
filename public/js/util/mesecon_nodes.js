export const node_image_mapping = {
    "mesecons_switch:mesecon_switch_off": "mesecons_switch_off.png",
    "mesecons_switch:mesecon_switch_on": "mesecons_switch_on.png",
    "digilines:lcd": "lcd_lcd.png"
};

const colors = ["red", "green", "blue", "gray", "darkgray", "yellow", "orange", "white", "pink", "magenta", "cyan", "violet"];
colors.forEach(c => {
    node_image_mapping[`mesecons_lightstone:lightstone_${c}_off`] = `jeija_lightstone_${c}_off.png`;
    node_image_mapping[`mesecons_lightstone:lightstone_${c}_on`] = `jeija_lightstone_${c}_on.png`;
});

export const switch_nodes = {
    "mesecons_switch:mesecon_switch_off": true,
    "mesecons_switch:mesecon_switch_on": true
};

export const display_nodes = {
    "digilines:lcd": true
};
