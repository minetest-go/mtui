
export const browse = dir => fetch(`api/filebrowser/browse?dir=${dir}`).then(r => r.json());

export const download = filename => fetch(`api/filebrowser/file?filename=${filename}`).then(r => r.blob());

export const get_download_url = filename => `api/filebrowser/file?filename=${filename}?download=true`;

export const upload = (filename, data) => fetch(`api/filebrowser/file?filename=${filename}`, {
    method: "POST",
    body: data
});

export const remove = filename => fetch(`api/filebrowser/file?filename=${filename}`, {
    method: "DELETE"
});

export const rename = (src, dst) => fetch(`api/filebrowser/rename?src=${src}&dst=${dst}`, {
    method: "POST"
});