import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { esbuildCommonjs } from '@originjs/vite-plugin-commonjs';
import mkcert from 'vite-plugin-mkcert';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [mkcert(), vue()],
  optimizeDeps: {
    esbuildOptions: {
      plugins: [
        esbuildCommonjs(),
      ]
    }
  },
  server: {
    https: true,
  },
})
