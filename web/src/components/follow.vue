<template>
  <el-main class="main-container wrap" ref="mainContainer">
    <!-- 页面标题区域 -->
    <div class="page-header">
      <h1 class="page-title">
        <el-icon class="title-icon"><UserFilled /></el-icon>
        已关注用户的作品
        <el-tag v-if="picitem.length > 0" type="info" class="count-tag">
          {{ picitem.length }} 作品
        </el-tag>
      </h1>
    </div>

    <!-- 控制面板 -->
    <el-card class="control-panel" shadow="hover">
      <el-row :gutter="20" class="control-row">
        <el-col :xs="24" :sm="12" :md="8" :lg="6">
          <div class="control-item">
            <label class="control-label">内容过滤</label>
            <el-select
                v-model="mode"
                :disabled="wait"
                size="large"
                style="width: 100%"
                @change="onModeChange"
            >
              <el-option value="all" label="全部内容" />
              <el-option value="safe" label="仅安全内容" />
              <el-option value="r18" label="仅R18内容" />
            </el-select>
          </div>
        </el-col>

        <el-col :xs="24" :sm="12" :md="10" :lg="8">
          <div class="control-item">
            <label class="control-label">批量下载页面范围</label>
            <div class="range-input-group">
              <el-input-number
                  v-model="from"
                  :min="1"
                  :max="100"
                  size="large"
                  controls-position="right"
                  @change="fixRange"
                  :disabled="wait"
              />
              <span class="range-separator">-</span>
              <el-input-number
                  v-model="to"
                  :min="1"
                  :max="100"
                  size="large"
                  controls-position="right"
                  @change="fixRange"
                  :disabled="wait"
              />
              <el-button
                  type="primary"
                  @click="downloadMultiplePages"
                  :disabled="wait"
                  :loading="downloading"
                  size="large"
              >
                <el-icon><Download /></el-icon>
                批量下载
              </el-button>
            </div>
          </div>
        </el-col>

        <!-- 当前页下载 -->
        <el-col :xs="24" :sm="24" :md="6" :lg="10">
          <div class="control-item current-page-download">
            <el-button
                type="success"
                size="large"
                @click="downloadCurrentPage"
                :disabled="wait || picitem.length === 0"
                :loading="downloading"
                class="download-current-btn"
            >
              <el-icon><Download /></el-icon>
              下载当前页 ({{ currentPage }})
            </el-button>
          </div>
        </el-col>
      </el-row>
    </el-card>
    <div class="pagination-container">
      <el-pagination
          background
          layout="prev, pager, next, jumper"
          :total="1000"
          :page-count="34"
          :current-page="currentPage"
          @current-change="handlePageChange"
          :disabled="wait"
          class="main-pagination"
      />
      <div class="floating-actions">
        <el-button
            type="primary"
            circle
            size="large"
            @click="refreshData"
            :loading="loading"
            class="refresh-fab"
            v-tooltip="{ content: '刷新页面', placement: 'left' }"
        >
          <el-icon><Refresh /></el-icon>
        </el-button>
      </div>
    </div>
    <div class="waterfall-container">
      <Waterfall
          ref="waterfall"
          :list="filteredPicItems"
          :width="300"
          :gutter="20"
          background-color="transparent"
          animation-effect="fadeInUp"
          key="followWaterfall"
      >
        <template #default="{ item, url, index }">
          <transition name="el-fade-in-linear">
            <div class="artwork-card">
              <PicCard
                  :author="item.Author"
                  :img="item.src"
                  :title="item.Title"
                  :pid="item.pid"
                  :authorId="item.authorId"
                  :pages="item.pages"
                  :r18="item.r18"
                  :key="item.pid + 'follow'"
              />
            </div>
          </transition>
        </template>
      </Waterfall>
    </div>

    <el-footer v-if="loading == true">
      <div class="loader" id="loader">
        <br>
        <br>
        <br>
        <br>
        <br>
        <br>
        <br>
        <div class="loading">
          <span></span>
          <span></span>
          <span></span>
          <span></span>
          <span></span>
        </div>
      </div>
    </el-footer>
  </el-main>

</template>

<script lang="ts" setup>
import { ref, onMounted, defineComponent, computed, nextTick } from "vue";
import { DownloadByFollowPage } from "../../bindings/main/internal/pixivlib/ctl.js";
import {DownloadFollow} from "../assets/js/download.js"
import { Waterfall } from 'vue-waterfall-plugin-next';
import axios from 'axios'
import 'vue-waterfall-plugin-next/dist/style.css'
import 'animate.css';
import PicCard from './PicCard.vue';
import "../assets/style/variable.less";
import {
  UserFilled,
  Download,
  Document,
  Filter,
  ArrowUp,
  Refresh
} from "@element-plus/icons-vue";
import {ElMessage} from "element-plus";
const waterfall = ref(null)
const from = ref(1),to=ref(10)
defineComponent({
  PicCard,
})
const picitem = ref([])
const currentPage = ref(1)
const name = ref("follow")
const mode = ref("all")
const wait = ref(false)
const loading = ref(true)
const downloading = ref(false)

// 计算属性：过滤后的图片列表
const filteredPicItems = computed(() => {
  if (mode.value === "all") {
    return picitem.value
  } else if (mode.value === "safe") {
    return picitem.value.filter(item => !item.r18)
  } else if (mode.value === "r18") {
    return picitem.value.filter(item => item.r18)
  }
  return picitem.value
})

// 计算属性：过滤后的数量
const filteredCount = computed(() => filteredPicItems.value.length)

function onModeChange() {
  console.log("内容过滤模式切换到:", mode.value)
  // 重新渲染瀑布流
  nextTick(() => {
    if (waterfall.value) {
      console.log("渲染开始")
      waterfall.value.renderer()
    }
  })
}

function refreshData() {
  fetchFollowData()
}

function fixRange() {
  from.value = Math.max(1, Math.min(30, from.value))
  to.value = Math.max(1, Math.min(30, to.value))
}
function fetchFollowData() {
  wait.value = true
  loading.value = true
  console.log("doing get")
  picitem.value = []
  axios.get("http://127.0.0.1:7234/api/followpage", {
    params: {
      p: currentPage.value.toString(),
      types: "illust",
      mode: mode.value,
    }
  }).then((res) => {
    console.log(res, res.data.data.length)
    let tmp = []
    for (var i = 0; i < res.data.data.length; i++) {
      picitem.value.push({ pid: res.data.data[i].id, Title: res.data.data[i].title, Author: res.data.data[i].userName, src: res.data.data[i].url, pages: res.data.data[i].countPage, authorId: res.data.data[i].userId, r18: res.data.data[i].r18 })
    }
    waterfall.value.renderer()
  }).catch((error) => {
    console.log(error, error)
  }).finally(() => {
    console.log("ok")
    loading.value = false
    wait.value = false
  })
}
function handlePageChange(page) {
  console.log("页面切换到:", page)
  currentPage.value = page
  fetchFollowData()
}

const debug = () => {
  console.log(picitem.value)
}
onMounted(function () {
  loading.value = true
  fetchFollowData()
  window.debug = debug
})
function downloadCurrentPage() {
  DownloadByFollowPage(currentPage.value.toString(), mode.value)
}

function downloadMultiplePages() {
  if (from.value > to.value) {
    const temp = to.value
    to.value = from.value
    from.value = temp
  }
  downloading.value = true
  try {
    DownloadFollow(from.value, to.value)
    ElMessage.success(`开始批量下载第 ${from.value} 到 ${to.value} 页的作品`)
  } catch (error) {
    ElMessage.error('批量下载失败，请稍后重试')
  } finally {
    setTimeout(() => {
      downloading.value = false
    }, 3000)
  }
}


</script>

<style lang="less" scoped>
@import "../assets/style/load.less";

.main-container {
  padding: 20px;
  min-height: 100vh;
  overflow-y: hidden;
  //background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  overflow-x: hidden;
}

// 页面标题
.page-header {
  text-align: center;
  margin-bottom: 30px;

  .page-title {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 15px;
    font-size: 28px;
    font-weight: 600;
    color: #ffffff;
    margin: 0;
    padding: 20px 0;
    text-shadow: 0 2px 4px rgba(0,0,0,0.1);

    .title-icon {
      font-size: 32px;
      color: #409EFF;
    }

    .count-tag {
      font-size: 14px;
      font-weight: normal;
    }
  }
}

// 控制面板
.control-panel {
  margin-bottom: 25px;
  border-radius: 15px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.08);

  :deep(.el-card__body) {
    padding: 25px;
  }

  .control-row {
    align-items: flex-end;
  }

  .control-item {
    .control-label {
      display: block;
      margin-bottom: 10px;
      font-weight: 600;
      color: #fdfdfd;
      font-size: 14px;
    }
  }

  .range-input-group {
    display: flex;
    align-items: center;
    gap: 10px;

    .range-separator {
      font-weight: bold;
      color: #909399;
      font-size: 16px;
    }

    .el-input-number {
      width: 100px;
    }

    @media (max-width: 768px) {
      flex-direction: column;

      .el-input-number {
        width: 100%;
      }

      .el-button {
        width: 100%;
        margin-top: 10px;
      }
    }
  }

  .current-page-download {
    display: flex;
    justify-content: flex-end;
    align-items: flex-end;

    .download-current-btn {
      width: 100%;
      height: 44px;
    }

    @media (max-width: 768px) {
      margin-top: 20px;
    }
  }
}

// 分页容器
.pagination-container {
  display: flex;
  //justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
  padding: 20px;
  background: black;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
  //width:100%;
  .stats-info {
    display: flex;
    gap: 10px;

    .el-tag {
      padding: 8px 12px;
      border-radius: 20px;
    }
  }

  @media (max-width: 768px) {
    //flex-direction: column;
    gap: 15px;

    .stats-info {
      flex-wrap: wrap;
      //justify-content: center;
    }
  }
}

// 瀑布流容器
.waterfall-container {
  min-height: 400px;

  .artwork-card {
    position: relative;
    border-radius: 15px;
    overflow: hidden;
    box-shadow: 0 4px 15px rgba(0,0,0,0.1);
    transition: all 0.3s ease;
    background: white;

    &:hover {
      box-shadow: 0 8px 30px rgba(0,0,0,0.15);
      transform: translateY(-5px);
    }
  }
}

// 加载状态
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;

  .loader {
    text-align: center;
  }

  .loading-text {
    margin-top: 20px;
    color: #606266;
    font-size: 16px;
  }
}

// 空状态
.empty-state {
  margin: 60px 0;

  .empty-description {
    p {
      margin: 5px 0;
      color: #909399;
    }
  }
}

// 浮动操作按钮
.floating-actions {
  //position: fixed;
  //right: 50px;
  //bottom: 130px;
  //z-index: 100;
  //float: left;
  .refresh-fab {
    width: 45px;
    height: 45px;
    box-shadow: 0 4px 15px rgba(64, 158, 255, 0.3);

    &:hover {
      transform: scale(1.1);
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .main-container {
    padding: 15px;
  }

  .page-header .page-title {
    font-size: 24px;
    flex-direction: column;
    gap: 10px;
  }
}

@media (max-width: 480px) {
  .page-header .page-title {
    font-size: 20px;
  }

  .waterfall-container {
    :deep(.vue-waterfall) {
      --waterfall-item-width: 100% !important;
    }
  }
}

// 动画效果
.el-fade-in-linear-enter-active {
  transition: all 0.4s ease;
}

.el-fade-in-linear-enter-from {
  opacity: 0;
  transform: translateY(30px);
}

.artwork-card {
  animation: slideInUp 0.6s ease forwards;
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(40px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>