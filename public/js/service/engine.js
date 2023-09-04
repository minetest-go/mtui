
export const store = Vue.reactive({
	versions: null,
	busy: false,
	status: null,
	version: ""
});

export const is_running = () => store.status ? store.status.running : false;