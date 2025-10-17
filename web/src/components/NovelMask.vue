<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
          v-if="novelViewerStore.isVisible"
          class="novel-viewer-modal"
          @click="handleBackdropClick"
          @keydown="handleKeyDown"
          tabindex="0"
          ref="modalRef"
      >
        <!-- 蒙版背景 -->
        <div class="modal-backdrop"></div>

        <!-- 主容器 -->
        <div class="modal-container">
          <!-- 头部信息栏 -->
          <div class="modal-header">
            <div >
              <!-- 章节导航 -->
              <div class="chapter-nav" v-if="novelViewerStore.seriesData != null">
                <el-tooltip placement="top" >

                  <template #content>
                    <div class="tips-content" style="width: 300px;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">
                      {{novelViewerStore.seriesData.prev.title}}
                    </div>
                  </template>
                  <el-button
                      type="primary"
                      @click="novelViewerStore.goToPrevious"
                      :disabled="!novelViewerStore.canGoPrevious  || novelViewerStore.isLoading"
                  >
                    <el-icon><ArrowLeft /></el-icon>
                    上一章
                  </el-button>
                </el-tooltip>

                <el-tooltip placement="top" >
                  <template #content>
                    <div class="tips-content" style="width: 300px;text-overflow: ellipsis;white-space: nowrap;overflow: hidden;">
                      {{novelViewerStore.seriesData.next.title}}
                    </div>
                  </template>
                  <el-button
                      type="primary"
                      @click="novelViewerStore.goToNext"
                      :disabled="!novelViewerStore.canGoNext  || novelViewerStore.isLoading"
                  >
                    下一章
                    <el-icon><ArrowRight /></el-icon>
                  </el-button>
                </el-tooltip>
              </div>
            </div>
            <div class="modal-controls">
              <!-- 字体控制 -->
              <div class="font-controls">
                <el-button-group class="font-buttons">
                  <el-button
                      size="small"
                      @click="decreaseFontSize"
                      :disabled="novelViewerStore.fontSize <= novelViewerStore.minFontSize || novelViewerStore.isLoading"
                  >
                    <el-icon><Minus /></el-icon>
                  </el-button>
                  <el-button size="small" disabled>
                    {{ novelViewerStore.fontSize }}px
                  </el-button>
                  <el-button
                      size="small"
                      @click="increaseFontSize"
                      :disabled="novelViewerStore.fontSize >= novelViewerStore.maxFontSize  || novelViewerStore.isLoading"
                  >
                    <el-icon><Plus /></el-icon>
                  </el-button>
                </el-button-group>
              </div>

              <!-- 章节列表按钮 -->
              <el-button
                  v-if="novelViewerStore.hasSeries"
                  circle
                  size="large"
                  @click="toggleSeriesSidebar"
                  v-tooltip="'章节列表'"
                  class="control-btn"
              >
                <el-icon><Menu /></el-icon>
              </el-button>

              <el-button
                  circle
                  size="large"
                  @click="downloadCurrent"
                  v-tooltip="'下载当前小说'"
                  class="control-btn"
              >
                <el-icon><Download /></el-icon>
              </el-button>
              <el-button
                  circle
                  size="large"
                  @click="openInPixiv"
                  v-tooltip="'在 Pixiv 中查看'"
                  class="control-btn"
              >
                <el-icon><Link /></el-icon>
              </el-button>
              <el-button
                  circle
                  size="large"
                  @click="novelViewerStore.closeViewer"
                  v-tooltip="'关闭 (ESC)'"
                  class="control-btn close-btn"
              >
                <el-icon><Close /></el-icon>
              </el-button>
            </div>
          </div>

          <!-- 小说内容显示区域 -->
          <div class="novel-display-area">
            <!-- 左侧导航按钮 -->
            <Transition name="nav-fade">
              <button
                  v-if="novelViewerStore.canGoPrevious && !novelViewerStore.isLoading"
                  class="nav-button nav-prev"
                  @click="novelViewerStore.goToPrevious"
                  v-tooltip="'上一章 (←)'"
              >
                <el-icon><ArrowLeft /></el-icon>
              </button>
            </Transition>

            <!-- 主内容容器 -->
            <div class="novel-container" ref="novelContainer">
              <div class="novel-content-wrapper" :style="contentStyle">
                <Transition name="loading-fade" v-if="novelViewerStore.isLoading">
                  <div  class="loading-container">
                    <div class="loading-spinner">
                      <div class="spinner-ring"></div>
                      <div class="spinner-ring"></div>
                      <div class="spinner-ring"></div>
                    </div>
                    <p class="loading-text">加载中...</p>
                  </div>
                </Transition>
                <el-scrollbar class="content-scrollbar"  v-else>
                  <div class="novel-info">
                    <el-text type="info" style="font-size: 16px" v-if="novelViewerStore.seriesData!==null">{{novelViewerStore.seriesData.title}} #{{novelViewerStore.seriesData.order}}</el-text>
                    <h3 class="novel-title">{{ novelViewerStore.novelTitle }}</h3>
                    <div class="novel-meta">
                      <span class="author-info">
                        <el-icon><User /></el-icon>
                        {{ novelViewerStore.author }}
                      </span>
                      <!-- Novel ID (可点击复制) -->
                      <span class="novel-id clickable" @click="copyToClipboard(novelViewerStore.currentNovelId, '小说ID')">
                        <el-icon><Document /></el-icon>
                        {{ novelViewerStore.currentNovelId }}
                      </span>

                      <!-- Series ID (可点击复制) -->
                      <span
                          v-if="novelViewerStore.seriesData!==null"
                          class="series-id clickable"
                          @click="copyToClipboard(novelViewerStore.seriesId, '系列ID')"
                      >
                        <el-icon><Collection /></el-icon>
                        系列: {{ novelViewerStore.seriesId }}
                      </span>

                    </div>
                    <div class="novel-meta" >
                      <div  style="float:left;display: flex;gap:0px 12px">
                        <div v-if="novelViewerStore.cover!==''">
                          <el-image :src="`http://127.0.0.1:7234/api/preview?url=${novelViewerStore.cover}`"
                                   class="image" style="width: 200px" fit="fill"></el-image>
                        </div>
                        <div v-html="novelViewerStore.description" @click="preventOutLink">
                        </div>
                      </div>
                    </div>
                    <div class="novel-meta">
                      <el-tag v-for="item in novelViewerStore.tags" :type="item==='R-18'? 'danger':'primary'" style="">
                        {{item}}
                      </el-tag>
                    </div>
                  </div>
                  <el-divider></el-divider>
                  <!-- 小说正文 -->
                  <div class="novel-text" v-if="!contentLoading">
                    <p v-for="(paragraph, index) in novelParagraphs" :key="index" class="paragraph">
                      {{ paragraph }}
                    </p>
                  </div>

                </el-scrollbar>
              </div>
            </div>

            <!-- 右侧导航按钮 -->
            <Transition name="nav-fade">
              <button
                  v-if="novelViewerStore.canGoNext && !novelViewerStore.isLoading"
                  class="nav-button nav-next"
                  @click="novelViewerStore.goToNext"
                  v-tooltip="'下一章 (→)'"
              >
                <el-icon><ArrowRight /></el-icon>
              </button>
            </Transition>
            <!-- 章节列表侧边栏 -->
            <Transition name="sidebar-slide">
              <div v-if="showSeriesSidebar && novelViewerStore.hasSeries" class="series-sidebar">
                <div class="sidebar-header">
                  <h4>章节列表</h4>
                  <div class="sidebar-header-buttons">
                    <el-button
                        circle
                        size="small"
                        @click="DownloadByNovelId(novelViewerStore.seriesId.toString(), true)"
                        v-tooltip="'下载整个系列'"
                        class="download-series-btn"
                    >
                      <el-icon><Download /></el-icon>
                    </el-button>
                    <el-button
                        circle
                        size="small"
                        @click="toggleSeriesSidebar"
                        v-tooltip="'关闭'"
                        class="close-sidebar-btn"
                    >
                      <el-icon><Close /></el-icon>
                    </el-button>
                  </div>
                </div>
                <el-scrollbar class="sidebar-scrollbar">
                  <div class="chapter-list">
                    <div
                        v-for="chapter in novelViewerStore.seriesList"
                        :key="chapter.id"
                        class="chapter-item"
                        :class="{ active: chapter.id === novelViewerStore.currentNovelId }"
                        @click="novelViewerStore.goToChapter(chapter.id)"
                    >
                      <span class="chapter-order">#{{ chapter.order }}</span>
                      <span class="chapter-title">{{ chapter.title }}</span>
                    </div>
                  </div>
                </el-scrollbar>
              </div>
            </Transition>
          </div>


          <!-- 快捷键提示 -->
          <div class="keyboard-shortcuts" v-if="showShortcuts">
            <div class="shortcuts-content">
              <h4>快捷键</h4>
              <div class="shortcut-list">
                <div class="shortcut-item">
                  <kbd>ESC</kbd>
                  <span>关闭</span>
                </div>
                <div class="shortcut-item">
                  <kbd>←</kbd>
                  <span>上一章</span>
                </div>
                <div class="shortcut-item">
                  <kbd>→</kbd>
                  <span>下一章</span>
                </div>
                <div class="shortcut-item">
                  <kbd>+ -</kbd>
                  <span>调整字号</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 快捷键提示按钮 -->
        <div class="help-button">
          <el-button
              circle
              size="large"
              @click="toggleShortcuts"
              v-tooltip="'快捷键帮助'"
              class="help-btn"
          >
            <el-icon><QuestionFilled /></el-icon>
          </el-button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import {
  User,
  Document,
  Link,
  Download,
  Close,
  ArrowLeft,
  ArrowRight,
  Loading,
  QuestionFilled,
  Plus,
  Minus,
  Menu,
  Collection
} from '@element-plus/icons-vue'
import {ElMessage, ElNotification} from 'element-plus'
import { useNovelViewerStore } from '../assets/stores/novelViewer.js'
import { copyToClipboard, copyLink, openPixivNovel } from '../assets/js/utils/index.js'
import axios from 'axios'
import {DownloadByNovelId} from "../../bindings/main/internal/pixivlib/ctl.js";

const novelViewerStore = useNovelViewerStore()

const modalRef = ref(null)
const novelContainer = ref(null)
const contentLoading = ref(false)
const showShortcuts = ref(false)
const showSeriesSidebar = ref(false)
const novelContent = ref('')


const toggleSeriesSidebar = () => {
  showSeriesSidebar.value = !showSeriesSidebar.value
}


async function preventOutLink(e) {
  const el = e.target
  if (el.tagName === 'A') {
    e.preventDefault()
    const href = el.getAttribute('href')
    await copyLink(href)
  }
}

// 计算属性
const contentStyle = computed(() => ({
  fontSize: `${novelViewerStore.fontSize}px`,
  lineHeight: novelViewerStore.lineHeight
}))

const novelParagraphs = computed(() => {
  if (!novelViewerStore.novelContent) return []
  return novelViewerStore.novelContent.split('\n').filter(p => p.trim())
})


// 字体控制
const increaseFontSize = () => {
  novelViewerStore.increaseFontSize()
}

const decreaseFontSize = () => {
  novelViewerStore.decreaseFontSize()
}

// 事件处理
const handleBackdropClick = (event) => {
  if (event.target === event.currentTarget) {
    novelViewerStore.closeViewer()
  }
}

const handleKeyDown = (event) => {
  novelViewerStore.handleKeyPress(event)
}

const downloadCurrent = async () => {
  try {

    DownloadByNovelId(novelViewerStore.currentNovelId)
    ElMessage.success(`开始下载: ${novelViewerStore.novelTitle}`)
    // 调用下载API
  } catch (error) {
    console.error('下载失败:', error)
    ElMessage.error('下载失败,请稍后重试')
  }
}

const openInPixiv = () => {
  openPixivNovel(novelViewerStore.currentNovelId)
}

const toggleShortcuts = () => {
  showShortcuts.value = !showShortcuts.value
}

// 监听键盘事件
onMounted(() => {
  document.addEventListener('keydown', novelViewerStore.handleKeyPress)
  if (modalRef.value) {
    modalRef.value.focus()
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', novelViewerStore.handleKeyPress)
})

// 监听当前页变化,加载新内容
watch(() => novelViewerStore.currentPage, () => {
  if (novelViewerStore.isVisible) {
    novelViewerStore.loadNovelContent()
  }
})

// 监听可见性变化
watch(() => novelViewerStore.isVisible, (newVal) => {
  if (newVal) {
    novelViewerStore.loadNovelContent()
  }
})
</script>


<style lang="less" scoped>
// 导入通用样式
@import "../assets/style/common/modal.less";
@import "../assets/style/common/buttons.less";
@import "../assets/style/common/loading.less";
@import "../assets/style/common/animations.less";

// 小说查看器特定样式
.novel-viewer-modal {

  // 小说内容包装器
  .novel-content-wrapper {
    width: 100%;
    height: calc(100% - 40px);
    background: rgba(0, 0, 0, 0.79);
    margin-top: auto;
    border-radius: 12px;
    padding: 30px 20px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);

    .content-scrollbar {
      height: 100%;

      .chapter-header {
        margin-bottom: 30px;
        padding-bottom: 15px;
        border-bottom: 2px solid #409EFF;

        .chapter-title {
          margin: 0;
          font-size: 24px;
          font-weight: 600;
          color: #ffffff;
        }
      }

      .novel-text {
        padding: 0 20px;

        .paragraph {
          margin: 0 0 1.5em 0;
          text-align: justify;
          text-indent: 2em;
          color: #e3e3e3;
          line-height: inherit;
        }
      }

      .content-loading {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 15px;
        color: #909399;
        padding: 60px 0;

        .loading-icon {
          font-size: 48px;
          animation: spin 1s linear infinite;
        }
      }

      .chapter-nav {
        display: flex;
        justify-content: space-between;
        margin-top: 40px;
        padding-top: 20px;
        border-top: 1px solid #e4e7ed;
      }
    }
  }
}

// 小说信息特定样式
.novel-info {
  flex: 1;
  color: #fbf1f1;

  .novel-title {
    margin: 0 0 8px 0;
    font-size: 20px;
    font-weight: 600;
    line-height: 1.3;
  }

  .novel-meta {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 10px 20px;
    font-size: 14px;
    opacity: 0.9;
    white-space: normal;
    height: auto;
    overflow-x: hidden;

    .author-info,
    .novel-id,
    .series-id {
      display: flex;
      align-items: center;
      gap: 6px;
      cursor: pointer;

      &.clickable:hover {
        color: #409EFF;
      }
    }
  }
}

// 字体控制
.font-controls {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;

  .font-buttons {
    .el-button {
      background: transparent;
      border-color: rgba(255, 255, 255, 0.3);
      color: white;

      &:hover:not(:disabled) {
        background: rgba(255, 255, 255, 0.1);
        border-color: rgba(255, 255, 255, 0.5);
      }

      &:disabled {
        opacity: 0.5;
      }
    }
  }
}

// 系列侧边栏
.series-sidebar {
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  width: 320px;
  background: rgba(0, 0, 0, 0.9);
  border-left: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  flex-direction: column;
  z-index: 20;

  .sidebar-header {
    padding: 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);

    h4 {
      margin: 0;
      color: white;
      font-size: 18px;
      font-weight: 600;
    }

    .sidebar-header-buttons {
      display: flex;
      gap: 8px;
    }
  }

  .sidebar-scrollbar {
    flex: 1;
    padding: 10px;

    .chapter-list {
      display: flex;
      flex-direction: column;
      gap: 8px;

      .chapter-item {
        padding: 12px 15px;
        background: rgba(255, 255, 255, 0.05);
        border-radius: 8px;
        cursor: pointer;
        transition: all 0.2s ease;
        display: flex;
        gap: 10px;
        align-items: flex-start;

        .chapter-order {
          color: #409EFF;
          font-weight: 600;
          flex-shrink: 0;
          min-width: 35px;
        }

        .chapter-title {
          color: rgba(255, 255, 255, 0.9);
          line-height: 1.4;
          flex: 1;
          word-break: break-word;
        }

        &:hover {
          background: rgba(255, 255, 255, 0.1);
          transform: translateX(-5px);
        }

        &.active {
          background: rgba(64, 158, 255, 0.2);
          border-left: 3px solid #409EFF;

          .chapter-title {
            color: white;
            font-weight: 500;
          }
        }
      }
    }
  }
}

// 章节导航
.chapter-nav {
  display: flex;
  gap: 10px;
}

// 提示文本样式
.tips-content {
  width: 300px;
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
}
</style>

<style lang="less" scoped>

// 小说显示区域
.novel-display-area {
  flex: 1;
  height: calc(100% - 160px);
  position: relative;
  display: flex;
  justify-content: center;
  background: rgba(0, 0, 0, 0.57);

  .novel-container {
    width: 100%;
    height: calc(100% - 60px);
    display: flex;
    //flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 20px 50px;

  }
}

</style>