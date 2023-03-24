import { defineConfig, loadEnv } from 'vite';
import mkcert from 'vite-plugin-mkcert';
import vue from '@vitejs/plugin-vue';
import VueRouter from 'unplugin-vue-router/vite';
import { esbuildCommonjs } from '@originjs/vite-plugin-commonjs';

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
        base: env.VITE_BASE ?? '/',
        server: {
            https: true,
        },
    };
});
