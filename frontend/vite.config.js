import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
    plugins: [vue()],
    server: {
        port: 6060,
        // proxy: {
        //     '/api': {
        //         target: 'http://localhost:3000',
        //         changeOrigin: true,
        //         secure: false,
        //     }
        // }
    },

})



