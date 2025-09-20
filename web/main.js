import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import router from './src/router/index.js'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import 'element-plus/theme-chalk/base.css'
// import {HomeFilled,StarFilled,Histogram,Search,Tools,Download} from '@element-plus/icons-vue'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './src/App.vue'
import * as config from "./src/assets/js/configuration.js"
import { createPinia } from "pinia"

const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
// app.component("HomeFilled",HomeFilled)
// app.component("StarFilled",StarFilled)
// app.component("Histogram",Histogram)
// app.component("Search",Search)
// app.component("Tools",Tools)
// app.component("Download",Download)


const resultElement = document.getElementById('result');
app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.mount('#app')

