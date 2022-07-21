import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
    plugins: [vue()],
    server: {
        port: 7000,
        // proxy: { // it doesn't works
        //     '/api': {
        //         target: 'http://localhost:3000/',
        //         changeOrigin: true
        //     }
        // }
    },
})



