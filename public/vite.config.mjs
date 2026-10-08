import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

// backend api, override with API_TARGET (for example "http://ui:8080")
const target = process.env.API_TARGET || 'http://localhost:8080';

// required by the wasm client (SharedArrayBuffer)
const headers = {
	'Cross-Origin-Embedder-Policy': 'credentialless',
	'Cross-Origin-Opener-Policy': 'same-origin',
	'Cross-Origin-Resource-Policy': 'cross-origin'
};

export default defineConfig({
	base: './',
	plugins: [vue()],
	// pics and wasm files, copied as-is to the build output
	publicDir: 'static',
	build: {
		outDir: 'dist',
		emptyOutDir: true,
		sourcemap: true
	},
	server: {
		headers,
		proxy: {
			'/api': {
				target,
				changeOrigin: true,
				// websocket proxy for the wasm client
				ws: true
			}
		}
	},
	preview: { headers }
});
