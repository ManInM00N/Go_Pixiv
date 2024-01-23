import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import router from './router/index.js'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/theme-chalk/base.css'
import App from './App.vue'
import V3waterfall from 'v3-waterfall'
import 'v3-waterfall/dist/style.css'
const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.use(router)
app.use(ElementPlus)
app.use(V3waterfall)
app.mount('#app')
