
export default {
    template: /*html*/`
    <div>
        <div class="card">
            <div class="card-header">
                Download backup <i class="fa fa-download"></i>
            </div>
            <div class="card-body">
                <a class="btn btn-primary" href="api/filebrowser/zip?dir=/">
                    <i class="fa fa-file-zipper"></i>
                    Download world-backup as zip-file
                </a>
            </div>
        </div>
    </div>
    `
};