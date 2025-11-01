<template>
  <div class="app-container">
    <!-- 左侧导航栏 -->
    <aside class="sidebar-container">
      <div class="sidebar-content">
        <!-- Logo区域 -->
        <div class="logo-section">
          <div class="logo-icon">
            <el-icon><img src="../../public/icon.ico" width="40" height="40"/></el-icon>
          </div>
          <div class="logo-text">Pixiv</div>
        </div>

        <!-- 主要导航 -->
        <nav class="main-navigation">
          <el-menu
              :default-active="activeIndex"
              class="sidebar-menu"
              mode="vertical"
              :router="true"
              @select="handleMenuSelect"
              :collapse="false"
              text-color="#e2e8f0"
              active-text-color="#ffffff"
              background-color="transparent"
          >
            <el-menu-item
                v-for="item in items"
                :key="item.id"
                :index="item.id"
                :route="item.index"
                :disabled="!item.logined && !form.pixivConf.logined"
                class="nav-item"
            >
              <div class="nav-item-content">
                <div class="nav-icon">
                  <el-icon :size="22">
                    <component :is="item.iconmsg" />
                  </el-icon>
                </div>
                <span class="nav-text">{{ item.key }}</span>
                <div v-if="!item.logined && !form.pixivConf.logined" class="lock-indicator">
                  <el-icon :size="14"><Lock /></el-icon>
                </div>
              </div>
            </el-menu-item>
          </el-menu>
        </nav>

        <!-- 用户状态指示 -->
        <div class="user-status">
          <div class="status-indicator" :class="{ 'online': form.pixivConf.logined }">
            <div class="status-dot"></div>
            <span class="status-text">
                {{ form.pixivConf.logined ? '已登录' : '未登录' }}
            </span>
          </div>
        </div>

        <!-- 底部设置 -->
        <div class="bottom-section">
          <el-menu
              :default-active="activeIndex"
              class="sidebar-menu"
              mode="vertical"
              :router="true"
              text-color="#94a3b8"
              active-text-color="#ffffff"
              background-color="transparent"
              @select="handleMenuSelect"
          >
            <el-menu-item
                :index="userself.id"
                :route="userself.index"
                class="nav-item"
            >
              <div class="nav-item-content">
                <div class="nav-icon">
                  <el-icon :size="22">
                    <Setting />
                  </el-icon>
                </div>
                <span class="nav-text">设置</span>
              </div>
            </el-menu-item>
          </el-menu>

          <!-- 版本信息 -->
          <div class="version-info">
            <a href="javascript:void(0);" @click="OpenInBrowser('https://github.com/ManInM00N/Go_Pixiv/releases')">v1.3.2</a>
          </div>
        </div>
      </div>
    </aside>

    <!-- 主内容区域 -->
    <main class="main-content">
      <div class="content-wrapper">
        <router-view v-slot="{ Component, route }">
          <transition name="fade-slide" mode="out-in">
            <keep-alive v-if="route.meta.keepAlive">
              <component
                  :is="Component"
                  :form="form"
                  :ws="ws"
                  :key="route.path"
              />
            </keep-alive>
            <component
                v-else
                :is="Component"
                :key="route.path"
            />
          </transition>
        </router-view>
      </div>
    </main>
    <el-backtop
        target=".main-content"
        :bottom="45"
        :right="45"
        :visibility-height="300"
    >
      <div class="back-to-top-btn">
        <el-icon :size="18"><ArrowUp /></el-icon>
      </div>
    </el-backtop>

    <div class="notification-container" id="notification-container"></div>
    <PicMask/>
    <NovelMask/>
  </div>
</template>
<script setup>
import { onMounted, ref, computed } from "vue";
import { Events } from "@wailsio/runtime";
import { ElMessage, ElNotification } from "element-plus";
import {
  HomeFilled,
  StarFilled,
  Histogram,
  Search,
  Setting,
  PictureFilled,
  Lock,
  ArrowUp
} from "@element-plus/icons-vue";
import { OpenInBrowser} from "../../bindings/main/internal/pixivlib/ctl.js";
import { form } from "../assets/js/configuration.js"
import { useRoute, useRouter } from "vue-router";
import { DownloadGIF } from "../assets/js/download.js";
import axios from "axios";
import PicMask from "./PicMask.vue";
import NovelMask from "./NovelMask.vue";

// 响应式数据
const activeIndex = ref('0')
const route = useRoute()
const router = useRouter()
const theme = ref('dark')
const set = ref(null)
const items = ref([
    { id: '0', iconmsg: "HomeFilled", key: "下载中心", index: "/", logined: true, description: "管理下载任务"},
    { id: '1', iconmsg: "StarFilled", key: "关注作品", index: "/follow", logined: false,description: "已关注用户的作品" },
    { id: '2', iconmsg: "Histogram", key: "排行榜", index: "/rank", logined: true,description: "热门作品排行" },
    { id : '3', iconmsg:"Reading" ,key:"小说",index: '/novel',logined: false , description: "小说"},
    { id: '4', iconmsg: "Search", key: "搜索", index: "/search", logined: true, description: "搜索作品和作者" },
])
const userself = ref({
    // id:6,key: "user",index:"/user"
    id: '10086', iconmsg: "Setting", key: "设置", index: "/setting", logined: true
})
const wait = ref(false)

// 菜单选择处理
function handleMenuSelect(index) {
  activeIndex.value = index
  console.log("菜单切换:", index)
}

Events.On("downloadugoira",function(msg){
  console.log("GIF下载事件:", msg.data,msg.data[0])
  DownloadGIF(msg.data[0][0],msg.data[0][1],msg.data[0][2],msg.data[0][3],msg.data[0][4],msg.data[0][5])
})
Events.On("login", function (msg) {
  console.log("登录事件:", msg.data)
  if (msg.data[0] === "True") {
    form.value.pixivConf.logined = true
    ElNotification({
      title: "登录成功",
      type: "success",
      message: "欢迎回来！现在可以访问所有功能",
      position: 'bottom-right',
      duration: 3000,
    })
  } else {
    form.value.pixivConf.logined = false
    ElNotification({
      title: "登录失败",
      type: "warning",
      message: msg.data[1] || "请检查Cookie配置",
      position: 'bottom-right',
      duration: 5000,
    })
  }
})
onMounted(function () {
    console.log("当前路径:", route.path)
    const currentItem = items.value.find(item => item.index === route.path) || userself.value

    if (currentItem) {
      activeIndex.value = currentItem.id
    } else if (route.path === userself.value.index) {
      activeIndex.value = userself.value.id
    }
    console.log(currentItem)
    ElNotification({
      type: "info",
      title: "系统启动",
      message: "正在检查登录状态...",
      position: 'bottom-right',
      duration: 2000,
    })
})

</script>


<style lang="less" scoped>
@import "../assets/style/color.less";
@import "../assets/style/menu.less";

</style>