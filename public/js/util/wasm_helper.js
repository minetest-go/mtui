/* globals _malloc, cwrap, HEAPU32, HEAPU8, stringToNewUTF8 */

function makeArgv(args) {
    // Assuming 4-byte pointers
    const argv = _malloc((args.length + 1) * 4);
    let i;
    for (i = 0; i < args.length; i++) {
        HEAPU32[(argv >>> 2) + i] = stringToNewUTF8(args[i]);
    }
    HEAPU32[(argv >>> 2) + i] = 0; // argv[argc] == NULL
    return [i, argv];
}

let emloop_pause, emloop_unpause, emloop_init_sound, emloop_invoke_main, emloop_install_pack;
let irrlicht_resize, irrlicht_force_pointerlock;
let emsocket_init, emsocket_set_proxy;

window.emloop_request_animation_frame = function() {
    emloop_pause();
    window.requestAnimationFrame(() => { emloop_unpause(); });
};

export function is_supported() {
    return !!window.SharedArrayBuffer;
}

function addPack(name) {
    return fetch(`wasm/packs/${name}.pack`)
    .then(r => r.blob())
    .then(blob => blob.arrayBuffer())
    .then(ab => new Uint8Array(ab))
    .then(u8arr => {
        const len = u8arr.byteLength;
        const offset = _malloc(len);
        HEAPU8.set(u8arr, offset);
        emloop_install_pack(stringToNewUTF8(name), offset, len);
    });
}

function setProgress(msg) {
    document.getElementById("loading-msg").innerText = msg;
}

export const ready = new Promise((resolve, reject) => {
    window.emloop_ready = function() {
        emloop_pause = cwrap("emloop_pause", null, []);
        emloop_unpause = cwrap("emloop_unpause", null, []);
        emloop_init_sound = cwrap("emloop_init_sound", null, []);
        emloop_invoke_main = cwrap("emloop_invoke_main", null, ["number", "number"]);
        emloop_install_pack = cwrap("emloop_install_pack", null, ["number", "number", "number"]);
        irrlicht_resize = cwrap("irrlicht_resize", null, ["number", "number"]);
        irrlicht_force_pointerlock = cwrap("irrlicht_force_pointerlock", null);
        emsocket_init = cwrap("emsocket_init", null, []);
        emsocket_set_proxy = cwrap("emsocket_set_proxy", null, ["number"]);

        setProgress("base pack");
        addPack("base")
        .then(() => resolve())
        .catch(() => reject("could not fetch base pack"));
    };
});

export function init(){
    document.body.innerHTML = /*html*/`
        <div id="error_div" style="display: none;">
            <div class="row">
                <div class="col-4"></div>
                <div class="col-4 alert alert-warning">
                    <i class="fa fa-triangle-exclamation"></i>
                    &nbsp;
                    <span id="error_msg"></span>
                </div>
                <div class="col-4"></div>
            </div>
        </div>
        <div id="app">
            <div class="row">
                <div class="col-4"></div>
                <div class="col-4 text-center">
                    <div style="height: 250px;"></div>
                    <i class="fa fa-spinner fa-spin"></i>
                    Loading client: <b id="loading-msg">initializing...</b>
                </div>
                <div class="col-4"></div>
            </div>
        </div>
        <canvas id="canvas" style="width: 100%; height: 100%; display: none;"></canvas>
    `;

    // provide global callbacks
    window.Module = {
        canvas: document.getElementById("canvas"),
        print: s => {
            console.log("print", s);
            if (s.startsWith("main() exited with return value")) {
                console.warn("game exited");
                show_error("game exited");
                location.reload();
            }
            if (s.startsWith("Unhandled exception:")) {
                console.warn("unhandled exception");
                show_error("unhandled exception");
            }
        },
        printErr: s => {
            console.warn("printErr", s);
            if (s.startsWith("emsocket_getaddrinfo: emsocket_read failed")) {
                // comes after "Unhandled exception:"
                console.warn("socket error");
            }
            if (s.indexOf("Connection timed out") >= 0) {
                // comes before "main() exited with return value 1"
                console.warn("connection timed out");
            }
            if (s.indexOf("Access denied. Reason: ") >= 0) {
                // "Timed out"
                // "Invalid password"
                // comes before "main() exited with return value 1"
                console.warn("access denied");
            }
        }
    };

    const script_el = document.createElement("script");
    script_el.src = "wasm/minetest.js";
    document.body.appendChild(script_el);

    return ready;
}

export function resize() {
    const width  = window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth;
    const height = window.innerHeight|| document.documentElement.clientHeight|| document.body.clientHeight;

    irrlicht_resize(width, height);
    irrlicht_force_pointerlock();
}

export function execute(args) {
    document.getElementById("app").remove();

    const canvas_el = document.getElementById("canvas");
    canvas_el.style.display = "block";

    const [argc, argv] = makeArgv(["./minetest", ...args]);
    emloop_invoke_main(argc, argv);

    emloop_init_sound();
    emsocket_init();
    emsocket_set_proxy(stringToNewUTF8(location.protocol.replace("http", "ws") + "//" + location.host + location.pathname + "api/wasm/proxy"));

    resize();
    window.addEventListener('resize', resize);
    window.emloop_request_animation_frame();
}

function show_error(err) {
    document.getElementById("canvas").remove();
    document.getElementById("error_div").style.display = "block";
    document.getElementById("error_msg").innerText = err;
}
