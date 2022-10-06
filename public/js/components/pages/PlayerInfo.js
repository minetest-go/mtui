import Userprofile from "../Userprofile.js";

export default {
	components: {
        "user-profile": Userprofile
    },
	template: /*html*/`
		<user-profile :username="$route.params.name"/>
	`
};
