import { defineConfig, loadEnv } from 'vite';
import mkcert from 'vite-plugin-mkcert';
import vue from '@vitejs/plugin-vue';
import VueRouter from 'unplugin-vue-router/vite';
import { esbuildCommonjs } from '@originjs/vite-plugin-commonjs';
import * as fs from 'fs';

const packageJson = fs.readFileSync('./package.json')
const version: string = JSON.parse(packageJson.toString()).version || '0.0.0';

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {

    const env = loadEnv(mode, process.cwd());

    return {
        plugins: [
            mkcert(),
            VueRouter({
                dataFetching: true,
                exclude: ['!*.component.vue'],
                extensions: ['.vue', '.md'],
                logs: true,
                routesFolder: [
                    {
                        src: './src/pages',
                    },
                ],
            }),
            vue(),
        ],
        base: env.VITE_BASE ?? '/',
        define: {
            __APP_VERSION__: '"' + version + '"',
        },
        optimizeDeps: {
            esbuildOptions: {
                plugins: [esbuildCommonjs()],
            },
        },
        build: {
            commonjsOptions: {
                transformMixedEsModules: true,
            },
            manifest: true,
        },
        server: {
            https: true,
            proxy: {
                '/api': 'http://localhost:8080',
            },
        },
    };
});
