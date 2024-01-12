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
            name:"maindownload",
            component:maindownload
        },
        {
            path: '/follow',
            name:"follow",
            component: follow,
        },
        {
            path: '/rank',
            name:"rank",
            component: rank,
        },
        {
            path: '/search',
            name:"search",
            component: search,

        },
        {
            path: '/setting',
            name:"setting",
            component: setting,
        },
        {
            path: '/user',
            name:"user",
            component: user,
        },
    ]

})
export default router