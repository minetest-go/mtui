import { init, execute } from "./wasm_helper.js";

init()
.then(() => {
    document.getElementById("app").remove();
    execute([
        "--go",
        "--address", "engine",
        "--port", "30000",
        "--name", "player", //TODO
        "--password", "password" //TODO
    ]);
});