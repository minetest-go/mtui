import { append, browse, remove, rename } from "../api/filebrowser.js";

export async function upload_chunked(dir, filename, data, progress_callback) {
    // temp filename to upload to
    const tmpfilename = filename + ".part";

    const list = await browse(dir);
    if (list.items.find(e => e.name == tmpfilename)){
        // remove tmp file
        await remove(dir + "/" + tmpfilename);
    }

    let offset = 0;
    do {
        const chunksize = Math.min(data.size - offset, 1000*1000); // 1 mb chunks
        await append(dir + "/" + tmpfilename, data.slice(offset, offset + chunksize));
        offset += chunksize;

        if (typeof(progress_callback) == "function") {
            progress_callback(offset / data.size); // 0...1
        }
    } while (offset < data.size);

    await rename(dir + "/" + tmpfilename, dir + "/" + filename);
}