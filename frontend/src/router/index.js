import {createRouter, createWebHashHistory} from 'vue-router'
import follow from "../components/follow.vue";
import rank from "../components/rank.vue";
import maindownload from "../components/maindownload.vue";
import search from "../components/search.vue";
import user from "../components/user.vue";
import setting from "../components/setting.vue";
const router=createRouter({
    history: createWebHashHistory(),
    routes:[
        {
            path: '/maindownload',
            component:maindownload,
        },
        {
            path: '/follow',
            component: follow,
        },
        {
            path: '/rank',
            component: rank,
        },
        {
            path: '/search',
            component: search,

        },
        {
            path: '/setting',
            component: setting,
        },
        {
            path: '/user',
            component: user,
        },
        {
            path: '/',
            redirect:'/maindownload',   //
        },
    ]

})
export default router