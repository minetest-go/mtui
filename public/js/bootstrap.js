// app bootstrap script
// hostname == "127.0.0.1": debug bundles
// else: prod bundles
const prod_scripts = [
    "node_modules/vue/dist/vue.global.prod.js",
    "node_modules/vue-router/dist/vue-router.global.prod.js",
];
const dev_scripts = [
    "node_modules/vue/dist/vue.global.js",
    "node_modules/vue-router/dist/vue-router.global.js",
];
const common_scripts = [
    "node_modules/@vuepic/vue-datepicker/dist/vue-datepicker.iife.js",
    "node_modules/chart.js/dist/chart.umd.js"
];

const scripts = [];
function add_script(src) {
    const script = document.createElement("script");
    script.src = src;
    scripts.push(script);
}

if (location.hostname == "localhost" || location.hostname == "127.0.0.1") {
    // dev
    dev_scripts.forEach(add_script);
} else {
    // prod
    prod_scripts.forEach(add_script);
}
common_scripts.forEach(add_script);

// local code
let script = document.createElement("script");
script.src = "js/bundle.js";
script.onerror = function() {
    // dev/module fallback if bundle is not found
    let script = document.createElement("script");
    script.src = "js/main.js";
    script.type = "module";
    document.body.appendChild(script);
};
scripts.push(script);

function load_script() {
    if (scripts.length == 0) {
        return;
    }

    let s = scripts.shift();
    s.onload = load_script;
    document.body.appendChild(s);
}

load_script();