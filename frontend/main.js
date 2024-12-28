import * as GreetService from "./bindings/main/greetservice.js";
import { Events } from "@wailsio/runtime";

import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import router from './src/router/index.js'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import 'element-plus/theme-chalk/dark/css-vars.css'
import 'element-plus/theme-chalk/base.css'
import App from './src/App.vue'
import V3waterfall from 'v3-waterfall'
import 'v3-waterfall/dist/style.css'
import * as config from "./src/assets/js/configuration.js"
const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
const resultElement = document.getElementById('result');
app.use(router)
app.use(ElementPlus)
app.use(V3waterfall)
app.mount('#app')

window.doGreet = () => {
    let name = document.getElementById('name').value;
    if (!name) {
        name = 'anonymous';
    }
    GreetService.Greet(name).then((result) => {
        resultElement.innerText = result;
    }).catch((err) => {
        console.log(err);
    });
}

