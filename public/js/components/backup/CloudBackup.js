
const store = Vue.reactive({
    endpoint: "",
    key_id: "",
    access_key: "",
    bucket: "",
    filename: "world.zip"
});

export default {
    data: () => store,
    template: /*html*/`
    <div>
        <table class="table table-striped">
            <tbody>
                <tr>
                    <td>
                        Endpoint
                    </td>
                    <td>
                        <input class="form-control" type="text" v-model="endpoint" placeholder="http(s) endpoint"/>
                    </td>
                </tr>
                <tr>
                    <td>
                        Key ID
                    </td>
                    <td>
                        <input class="form-control" type="text" v-model="key_id" placeholder="Key identifier"/>
                    </td>
                </tr>
                <tr>
                    <td>
                        Access key
                    </td>
                    <td>
                        <input class="form-control" type="text" v-model="access_key" placeholder="Secret key"/>
                    </td>
                </tr>
                <tr>
                    <td>
                        Bucket
                    </td>
                    <td>
                        <input class="form-control" type="text" v-model="bucket" placeholder="Name of the bucket"/>
                    </td>
                </tr>
                <tr>
                    <td>
                        Filename
                    </td>
                    <td>
                        <input class="form-control" type="text" v-model="filename" placeholder="Filename to use"/>
                    </td>
                </tr>
                <tr>
                    <td></td>
                    <td>
                        <button class="btn btn-primary w-100">
                            Start
                            <i class="fa fa-play"></i>
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
    `
};