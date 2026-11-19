import { START, FILEBROWSER } from "../../../util/breadcrumbs.js";

const modes = [{
    name: "lua", match: /.*(lua)$/i
},{
    name: "javascript", match: /.*(js|json)$/i
},{
    name: "htmlmixed", match: /.*(html)$/i
},{
    name: "xml", match: /.*(xml)$/i
},{
    name: "text/css", match: /.*(css)$/i
},{
    name: "toml", match: /.*(toml)$/i
}];

export const can_edit = filename => {
    return filename.match(/.*(js|lua|txt|conf|cfg|json|md|mt|html|css|toml)$/i);
};

export const get_mode_name = filename => {
    const mode = modes.find(m => filename.match(m.match));
    return mode ? mode.name : null;
};
// breadcrumb for the given path inside the filebrowser, the last item is shown as file if is_file is set
export const filebrowser_breadcrumb = (pathMatch, is_file) => {
    const bc = [START, FILEBROWSER];
    let path = "";
    pathMatch.split("/")
    .filter(p => p != "")
    .forEach(p => {
        path = path == "" ? p : path + "/" + p;
        bc.push({
            name: p,
            icon: "folder-open",
            link: "/filebrowser/" + path
        });
    });

    if (is_file) {
        const lastitem = bc[bc.length-1];
        lastitem.icon = "file";
        lastitem.link = null;
    }

    return bc;
};
