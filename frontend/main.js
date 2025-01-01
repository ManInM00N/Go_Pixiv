import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import router from './src/router/index.js'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/theme-chalk/base.css'
import App from './src/App.vue'
import * as config from "./src/assets/js/configuration.js"
const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
const resultElement = document.getElementById('result');
app.use(router)
app.use(ElementPlus)
app.mount('#app')

