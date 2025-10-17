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
import axios from 'axios'
import {DownloadByNovelId} from "../../bindings/main/internal/pixivlib/ctl.js";

const novelViewerStore = useNovelViewerStore()

const modalRef = ref(null)
const novelContainer = ref(null)
const contentLoading = ref(false)
const showShortcuts = ref(false)
const showSeriesSidebar = ref(false)
const novelContent = ref('')

const copyToClipboard = async (text, label = '内容') => {
  try {
    await navigator.clipboard.writeText(text)
    ElNotification({
      position: "bottom-right",
      type: "success",
      message: `${label}已复制: ${text}`,
      duration: 2000,
    })
  } catch (err) {
    console.error('复制失败', err)
    ElNotification({
      position: "bottom-right",
      type: "warning",
      message: "复制失败",
      duration: 1000,
    })
  }
}

const toggleSeriesSidebar = () => {
  showSeriesSidebar.value = !showSeriesSidebar.value
}


async function  preventOutLink(e) {
  const el = e.target
  if (el.tagName === 'A') {
    e.preventDefault() // 阻止默认跳转
    const href = el.getAttribute('href')
    try {
      await navigator.clipboard.writeText(href)
      ElNotification({
        position:"bottom-right",
        type: "info",
        message: "链接已复制到剪切板",
        duration: 1000,
      })
    } catch (err) {
      console.error('复制失败', err)
      ElNotification({
        position:"bottom-right",
        type:"warning",
        message:"链接复制失败",
        duration: 1000,
      })
    }
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
  const url = `https://www.pixiv.net/novel/show.php?id=${novelViewerStore.currentNovelId}`
  window.open(url, '_blank')
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
// 模态框动画
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: all 0.3s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

// 导航按钮动画
.nav-fade-enter-active,
.nav-fade-leave-active {
  transition: all 0.2s ease;
}

.nav-fade-enter-from,
.nav-fade-leave-to {
  opacity: 0;
  transform: scale(0.8);
}


// 侧边栏动画
.sidebar-slide-enter-active,
.sidebar-slide-leave-active {
  transition: all 0.3s ease;
}

.sidebar-slide-enter-from,
.sidebar-slide-leave-to {
  transform: translateX(100%);
}

// 加载动画
.loading-fade-enter-active,
.loading-fade-leave-active {
  transition: all 0.3s ease;
}

.loading-fade-enter-from,
.loading-fade-leave-to {
  opacity: 0;
}

// 内容淡入动画
.content-fade-enter-active {
  transition: all 0.4s ease;
}

.content-fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

// 主容器
.novel-viewer-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 201;
  display: flex;
  align-items: center;
  justify-content: center;
  outline: none;

  .modal-backdrop {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.16);
    backdrop-filter: blur(8px);
  }

  .modal-container {
    position: relative;
    width: 95vw;
    height: 95vh;
    max-width: 1200px;
    max-height: 900px;
    display: flex;
    flex-direction: column;
    border-radius: 20px;
    overflow: hidden;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  }
}
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
    height:auto;
    overflow-x: hidden;
    .author-info,
    .novel-id {
      display: flex;
      align-items: center;
      gap: 6px;
    }
  }
}

// 加载容器
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  color: #e3e3e3;

  .loading-spinner {
    position: relative;
    width: 80px;
    height: 80px;
    margin-bottom: 30px;

    .spinner-ring {
      position: absolute;
      width: 100%;
      height: 100%;
      border: 3px solid transparent;
      border-top-color: #409EFF;
      border-radius: 50%;
      animation: spin 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;

      &:nth-child(1) {
        animation-delay: -0.45s;
      }

      &:nth-child(2) {
        animation-delay: -0.3s;
      }

      &:nth-child(3) {
        animation-delay: -0.15s;
      }
    }
  }

  .loading-text {
    font-size: 18px;
    color: #909399;
    margin: 0;
    animation: pulse 1.5s ease-in-out infinite;
  }
}

// 头部信息栏
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 30px 0 30px;
  background: rgba(0, 0, 0, 0.57);
  color: white;

  .modal-controls {
    display: flex;
    gap: 10px;
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

    .control-btn {
      background: rgba(255, 255, 255, 0.1);
      border: 1px solid rgba(255, 255, 255, 0.2);
      color: white;

      &:hover {
        background: rgba(255, 255, 255, 0.2);
        transform: scale(1.05);
      }

      &.close-btn {
        background: rgba(245, 108, 108, 0.8);
        border-color: rgba(245, 108, 108, 0.8);

        &:hover {
          background: rgba(245, 108, 108, 1);
        }
      }
    }
  }
}

// 小说显示区域
.novel-display-area {
  flex: 1;
  //flex-direction: column;
  height: calc(100% - 160px);
  position: relative;
  display: flex;
  //align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.57);

  .nav-button {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    width: 60px;
    height: 60px;
    border-radius: 50%;
    background: rgba(0, 0, 0, 0.6);
    color: white;
    border: none;
    cursor: pointer;
    z-index: 10;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    transition: all 0.3s ease;

    &:hover {
      background: rgba(0, 0, 0, 0.8);
      transform: translateY(-50%) scale(1.1);
    }

    &.nav-prev {
      left: 30px;
    }

    &.nav-next {
      right: 30px;
    }
  }

  .novel-container {
    width: 100%;
    height: calc(100% - 60px);
    display: flex;
    //flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 20px 50px;

    .novel-content-wrapper {
      width: 100%;
      //max-height: calc(100% - 60px);
      height: calc(100% - 40px);
      background: rgba(0, 0, 0, 0.79);
      margin-top: auto;
      border-radius: 12px;
      padding:30px 20px;
      box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
      .content-scrollbar {
        //height: 100%;
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
          .paragraph {
            margin: 0 0 1.5em 0;
            text-align: justify;
            text-indent: 2em;
            color: #e3e3e3;
            line-height: inherit;
          }
          padding: 0 20px
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

      .close-sidebar-btn
      .download-series-btn {
        background: rgba(255, 255, 255, 0.1);
        border-color: rgba(255, 255, 255, 0.2);
        color: white;

        &:hover {
          background: rgba(255, 255, 255, 0.2);
        }
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
}

// 分页控制
.pagination-controls {
  display: flex;
  flex-direction: column;
  gap: 15px;
  padding: 0px 30px 8px 30px;
  background: rgba(0, 0, 0, 0.57);

  .page-info {
    text-align: center;
    font-size: 16px;
    font-weight: 600;
    color: white;

    .current-page {
      color: #409EFF;
    }

    .page-separator {
      margin: 0 8px;
      color: rgba(255, 255, 255, 0.6);
    }
  }

  .chapter-list {
    display: flex;
    gap: 8px;
    overflow-x: auto;
    padding: 4px 0;

    .chapter-btn {
      flex-shrink: 0;
      min-width: 40px;
    }
  }
}

// 快捷键提示
.keyboard-shortcuts {
  position: absolute;
  bottom: 80px;
  right: 30px;
  background: rgba(0, 0, 0, 0.9);
  color: white;
  border-radius: 12px;
  padding: 20px;
  min-width: 200px;

  .shortcuts-content {
    h4 {
      margin: 0 0 15px 0;
      font-size: 16px;
      font-weight: 600;
    }

    .shortcut-list {
      display: flex;
      flex-direction: column;
      gap: 8px;

      .shortcut-item {
        display: flex;
        justify-content: space-between;
        align-items: center;

        kbd {
          background: #555;
          padding: 4px 8px;
          border-radius: 4px;
          font-family: monospace;
          font-size: 12px;
        }

        span {
          font-size: 14px;
        }
      }
    }
  }
}

// 帮助按钮
.help-button {
  position: absolute;
  bottom: 30px;
  right: 30px;

  .help-btn {
    background: rgba(0, 0, 0, 0.6);
    border: 1px solid rgba(255, 255, 255, 0.2);
    color: white;

    &:hover {
      background: rgba(0, 0, 0, 0.8);
      transform: scale(1.05);
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .modal-container {
    width: 100vw;
    height: 100vh;
    border-radius: 0;
  }

  .modal-header {
    padding: 15px 20px;

    .novel-title {
      font-size: 18px;
    }

    .novel-meta {
      gap: 15px;
      font-size: 13px;
    }

    .modal-controls {
      gap: 8px;

      .font-controls {
        padding: 6px 10px;
      }
    }
  }

  .novel-display-area {
    .novel-container {
      padding: 15px 80px;

      .novel-content-wrapper {
        padding: 20px;
      }
    }

    .nav-button {
      width: 50px;
      height: 50px;
      font-size: 20px;

      &.nav-prev {
        left: 15px;
      }

      &.nav-next {
        right: 15px;
      }
    }
  }

  .pagination-controls {
    padding: 15px 20px;
  }

  .keyboard-shortcuts {
    display: none;
  }

  .help-button {
    bottom: 20px;
    right: 20px;
  }
}

// 动画
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

</style>