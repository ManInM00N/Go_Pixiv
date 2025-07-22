<template>
    <el-aside name="menu" id="menu" class="asidemenu" width="60px">
        <el-menu default-active="0" background-color="#626369" :theme="theme" class="vertical-menu" mode="vertical"
            :router="true" @select="handleMenuSelect">
            <section class="top-items" height=240>
                <el-menu-item v-for="(item, idx) in items" :key="item.key" :id="item.id" :index="item.id"
                    :route="item.index" class="menu_item" :disabled="!item.logined && !form.logined"
                    :limit="form['r-18']" style="padding-left: 18px;">

                    <el-tooltip :content="item.key" :show-arrow=false placement="bottom-start" offset=2>
                        <el-container class="item_body">
                            <el-icon size="30" class="item_icon">
                                <component :is="item.iconmsg" />
                            </el-icon>
                        </el-container>
                    </el-tooltip>
                </el-menu-item>
            </section>
            <section class="placeholder-item "></section>
            <section class="bottom-items">
                <el-menu-item :key="userself.key" :id="userself.id" :index="userself.id" :route="userself.index"
                    class="menu_item" style="padding-left: 18px;">

                    <el-tooltip :content="userself.key" :show-arrow=false placement="bottom-start" offset='2'>
                        <el-container class="item_body">
                            <el-icon size="30" class="item_icon">
                                <Tools />
                            </el-icon>
                        </el-container>
                    </el-tooltip>
                </el-menu-item>
            </section>
        </el-menu>

    </el-aside>
    <el-main class="View glass" id="View" name="View"
        style="padding-right: 5px;padding-left: 5px;padding-bottom: 0px;padding-top: 0px;">
        <section
            style="width: 100%;height: calc(100% - 40px); margin-left: 15px; margin-right: 15px;margin-top: 20px;margin-bottom: 20px;">
            <section style="width:100%">

                <router-view v-slot="{ Component, route, Path }">
                    <keep-alive v-if="route.meta.keepAlive">
                        <component :is="Component" :form="form" :ws="ws" />
                    </keep-alive>
                    <component v-else :is="Component" />

                </router-view>
            </section>
        </section>
        <!--    <keep-alive>-->
        <!--      <router-view v-if="$route.meta.keepAlive"></router-view>-->
        <!--    </keep-alive>-->
        <!--    <router-view v-if="!$route.meta.keepAlive"></router-view>-->
    </el-main>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { Events } from "@wailsio/runtime";
import axios from "axios";
import { ElMessage } from "element-plus";
import { ElNotification } from 'element-plus'
import { form, ws } from "../assets/js/configuration.js"
import { useRoute, useRouter } from "vue-router";
import {DownloadGIF} from "../assets/js/download.js";
const theme = ref('dark')
const set = ref(null)
const items = ref([
    { id: '0', iconmsg: "HomeFilled", key: "maindownload", index: "/", logined: true },
    { id: '1', iconmsg: "StarFilled", key: "follow", index: "/follow", logined: false },
    { id: '2', iconmsg: "Histogram", key: "rank", index: "/rank", logined: true },
    { id: '3', iconmsg: "Search", key: "search", index: "/search", logined: true },
])
const activeIndex = ref('0')
const userself = ref({
    // id:6,key: "user",index:"/user"
    id: '4', iconmsg: "Tools", key: "settings", index: "/setting", logined: true
})
const wait = ref(false)
const route = useRoute()
const router = useRouter()
function waitchange(val) {
    wait.value = val
}
function handleMenuSelect(index) {
    activeIndex.value = index
    console.log("ee", index, activeIndex.value)

}
ElNotification({
    type: "info",
    title: "INFO",
    message: "Login ......",
    position: 'bottom-right',
    duration: 3000,
})
Events.On("downloadugoira",function(msg){
  console.log(msg, msg.data[0], msg.data)
  DownloadGIF(msg.data[0][0],msg.data[0][1],msg.data[0][2],msg.data[0][3],msg.data[0][4])
})
Events.On("login", function (msg) {
    console.log(msg, msg.data[0], msg.data)
    if (msg.data[0] === "True") {
        form.value.logined = true
        ElNotification({
            title: "Login succeed",
            type: "success",
            message: "Enjoy 🥳",
            position: 'bottom-right',
            duration: 3000,
        })
    } else {
        ElNotification({
            title: "Login failed",
            type: "warning",
            message: msg.data[1] + " 😭",
            position: 'bottom-right',
            duration: 5000,
        })
    }
    console.log("Login state: ", msg.data[0], "r18 mod: ", form.value.r_18,)
})
onMounted(function () {
    console.log(route.path)
    localStorage.setItem("cookie", form.value.cookie)
    console.log(localStorage.getItem("cookie"))
    // activeIndex.value = 1
    // startWebSocket()

})

</script>

<style lang="less" scoped>
@import "src/assets/style/color.less";
@import "src/assets/style/menu.less";

.tempheight {
    height: 60px
}

.glass {
    background-color: rgba(89, 89, 89, 0.15);
    backdrop-filter: blur(6px);
    -webkit-backdrop-filter: blur(6px);
    box-shadow: rgba(14, 14, 14, 0.19) 0px 6px 15px 0px;
    -webkit-box-shadow: rgba(14, 14, 14, 0.19) 0px 6px 15px 0px;
    border-radius: 0px;
    -webkit-border-radius: 0px;
    color: rgb(128, 128, 128);
}
</style>
