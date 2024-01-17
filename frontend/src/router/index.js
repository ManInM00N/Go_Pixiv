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
            meta: {
                keepAlive: true,
                refresh: false,
            },
        },
        {
            path: '/follow',
            component: follow,
            meta: {
                keepAlive: true,
            },
            beforeRouteLeave(to,from,next){
                to.meta.keepAlive = true
                next(0)
            }
        },
        {
            path: '/rank',
            component: rank,
            meta: {
                keepAlive: true,
            },
            beforeRouteLeave(to,from,next){
                to.meta.keepAlive = true
                next(0)
            }
        },
        {
            path: '/search',
            component: search,

            meta: {
                keepAlive: true,
                refresh: false,
            },

        },
        {
            path: '/setting',
            component: setting,
            meta: {
                keepAlive: true,
                refresh: false,
            },
        },
        {
            path: '/user',
            component: user,
            meta: {
                keepAlive: true,
                refresh: false,
            },
        },
        {
            path: '/',
            redirect:'/maindownload',
            meta: {
                keepAlive: true,
                refresh: false,
            },
        },
    ]
})
export default router