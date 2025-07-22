import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
export default defineConfig( {
    plugins:[vue(),vueJsx()],
    resolve: {
        extensions: [ 'ts', 'js', 'tsx', 'jsx' ]
    },
    server: {
        proxy: {
            '/apis': {
                target: 'http://127.0.0.1:7234',
                changeOrigin: true,
                pathRewrite: {
                    '^/apis': ''
                }
            }
        }
    }
}
)