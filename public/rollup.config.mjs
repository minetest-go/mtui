import terser from '@rollup/plugin-terser';

export default [{
	input: 'js/main.js',
	output: {
		file :'js/bundle.js',
		format: 'iife',
		sourcemap: true,
		compact: true,
		plugins: [terser()]
	}
}];
