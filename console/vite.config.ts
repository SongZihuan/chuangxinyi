import vue from '@vitejs/plugin-vue';
import { resolve } from 'path';
import { defineConfig, loadEnv, ConfigEnv } from 'vite';
import vueSetupExtend from 'vite-plugin-vue-setup-extend-plus';
import viteCompression from 'vite-plugin-compression';
import { buildConfig } from './src/utils/build';
import mockPlugin from './build/vite/plugin';
import createSvgIcon from './build/vite/svg-icon';
import svgLoader from 'vite-svg-loader';
import versionUpdatePlugin from './src/utils/versionUpdatePlugin';
const pathResolve = (dir: string) => {
	return resolve(__dirname, '.', dir);
};

const alias: Record<string, string> = {
	'/@': pathResolve('./src/'),
};
const CurrentTimeVersion = new Date().getTime();
const viteConfig = defineConfig((mode: ConfigEnv) => {
	const env = loadEnv(mode.mode, process.cwd());
	return {
		plugins: [
			vue(),
			vueSetupExtend(),
			viteCompression(),
			createSvgIcon(),
			JSON.parse(env.VITE_OPEN_CDN) ? buildConfig.cdn() : null,
			mockPlugin(false),
			svgLoader(),
			// @ts-ignore
			versionUpdatePlugin({
				version: CurrentTimeVersion,
			}),
		],
		root: process.cwd(),
		resolve: { alias },
		base: mode.command === 'serve' ? './' : env.VITE_PUBLIC_PATH,
		optimizeDeps: {
			exclude: ['vue-demi'],
		},
		server: {
			host: '0.0.0.0',
			port: 8080,
			open: JSON.parse(env.VITE_OPEN),
			proxy: {
				'/adminApi': {
					target: env.VITE_PROXY_TARGET || 'http://localhost:3350/api/v1/',
					changeOrigin: true,
					rewrite: (path) => path.replace(/^\/adminApi/, ''),
				},
			},
		},
		build: {
			outDir: 'dist',
			minify: 'terser',
			terserOptions: {
				compress: {
					//生产环境时移除console
					drop_console: env.VITE_NOT_CONSOLE === "true",
					drop_debugger: env.VITE_NOT_CONSOLE === "true",
				},
			},
			chunkSizeWarningLimit: 1500,
			rollupOptions: {
				output: {
					chunkFileNames: 'assets/js/[name]-[hash].js',
					entryFileNames: 'assets/js/[name]-[hash].js',
					assetFileNames: 'assets/[ext]/[name]-[hash].[ext]',
					manualChunks(id) {
						if (id.includes('node_modules')) {
							return id.toString().match(/\/node_modules\/(?!.pnpm)(?<moduleName>[^\/]*)\//)?.groups!.moduleName ?? 'vender';
						}
					},
				},
				...(JSON.parse(env.VITE_OPEN_CDN) ? { external: buildConfig.external } : {}),
			},
		},
		css: { preprocessorOptions: { css: { charset: false } } },
		define: {
			__NEXT_VERSION__: JSON.stringify(process.env.npm_package_version),
			__NEXT_NAME__: JSON.stringify(process.env.npm_package_name),
		},
	};
});

export default viteConfig;
