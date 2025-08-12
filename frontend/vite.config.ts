import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		port: 3000,
		proxy: {
			'/api': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/login': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/register': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/refresh': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/logout': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/auth/google': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/auth/github': {
				target: 'http://localhost:8080',
				changeOrigin: true
			}
		}
	}
});
