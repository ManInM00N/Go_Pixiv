import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import router from './router/index.js'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import MasonryWall from '@yeger/vue-masonry-wall'
import { VueMasonryPlugin } from 'vue-masonry';
const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.use(VueMasonryPlugin);
app.use(router)
app.use(ElementPlus)
app.use(MasonryWall)
app.mount('#app')
