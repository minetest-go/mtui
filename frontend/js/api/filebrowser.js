
export const browse = dir => fetch(`api/filebrowser/browse?dir=${dir}`).then(r => r.json());

export const download = filename => fetch(`api/filebrowser/file?filename=${filename}`).then(r => r.blob());

export const download_text = filename => fetch(`api/filebrowser/file?filename=${filename}`).then(r => r.text());

export const get_download_url = filename => `api/filebrowser/file?filename=${filename}&download=true`;

export const get_zip_url = dir => `api/filebrowser/zip?dir=${dir}`;
export const get_targz_url = dir => `api/filebrowser/targz?dir=${dir}`;

function postProgress(url, body, callback) {
    callback = callback ? callback : () => {};

    return new Promise(resolve => {
        var request = new XMLHttpRequest();
        request.upload.addEventListener('progress', e => {
            const progress = e.total / e.loaded;
            callback(progress);
        });
        request.upload.addEventListener('load', () => resolve());
        request.open('POST', url);
        request.send(body);
    });
}

export const append = (filename, data, offset) => fetch(`api/filebrowser/file?filename=${filename}&offset=${offset}`, {
    method: "PUT",
    body: data
});

export const upload = (filename, data, callback) => postProgress(`api/filebrowser/file?filename=${filename}`, data, callback);

export const upload_zip = (dir, data, callback) => postProgress(`api/filebrowser/zip?dir=${dir}`, data, callback);

export const upload_targz = (dir, data, callback) => postProgress(`api/filebrowser/targz?dir=${dir}`, data, callback);

export const unzip = filename => fetch(`api/filebrowser/unzip?filename=${filename}`, {
    method: "POST"
});

export const mkdir = dir => fetch(`api/filebrowser/mkdir?dir=${dir}`, {
    method: "POST"
});

export const remove = filename => fetch(`api/filebrowser/file?filename=${filename}`, {
    method: "DELETE"
});

export const rename = (src, dst) => fetch(`api/filebrowser/rename?src=${src}&dst=${dst}`, {
    method: "POST"
});