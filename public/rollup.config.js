
module.exports = [{
	input: 'js/main.js',
	output: {
		file :'js/bundle.js',
		format: 'iife',
		sourcemap: true,
		compact: true
	}
},{
	input: 'js/wasm_main.js',
	output: {
		file :'js/wasm_bundle.js',
		format: 'iife',
		sourcemap: true,
		compact: true
	}
}];
