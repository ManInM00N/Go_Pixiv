import { createRouter, createWebHashHistory } from 'vue-router'
import follow from "../components/follow.vue";
import rank from "../components/rank.vue";
import maindownload from "../components/maindownload.vue";
import search from "../components/search.vue";
import setting from "../components/settings.vue";
const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            component: maindownload,
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
        },
        {
            path: '/rank',
            component: rank,
            meta: {
                keepAlive: true,
            },

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
    ]
})
export default router
