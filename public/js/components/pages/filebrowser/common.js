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