<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
          v-if="imageViewerStore.isVisible"
          class="image-viewer-modal"
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
            <div class="image-info">
              <h3 class="image-title">{{ imageViewerStore.imageTitle }}</h3>
              <div class="image-meta">
                                <span class="author-info">
                                    <el-icon><User /></el-icon>
                                    {{ imageViewerStore.author }}
                                </span>
                <span class="pid-info">
                                    <el-icon><Link /></el-icon>
                                    {{ imageViewerStore.currentPid }}
                                </span>
                <el-tag v-if="imageViewerStore.isR18" type="danger" size="small">
                  R18
                </el-tag>
              </div>
            </div>

            <div class="modal-controls">
              <!-- 缩放控制 -->
              <div class="zoom-controls">
                <span class="zoom-info">{{ Math.round(imageViewerStore.scale * 100) }}%</span>
                <el-button-group class="zoom-buttons">
                  <el-button
                      size="small"
                      @click="imageViewerStore.zoomOut()"
                      :disabled="imageViewerStore.scale <= imageViewerStore.minScale"
                  >
                    <el-icon><Minus /></el-icon>
                  </el-button>
                  <el-button
                      size="small"
                      @click="imageViewerStore.fitToScreen()"
                  >
                    <el-icon><FullScreen /></el-icon>
                  </el-button>
                  <el-button
                      size="small"
                      @click="imageViewerStore.zoomIn()"
                      :disabled="imageViewerStore.scale >= imageViewerStore.maxScale"
                  >
                    <el-icon><Plus /></el-icon>
                  </el-button>
                </el-button-group>
              </div>

              <el-button
                  circle
                  size="large"
                  @click="downloadCurrent"
                  v-tooltip="'下载当前作品'"
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
                  @click="imageViewerStore.closeViewer"
                  v-tooltip="'关闭 (ESC)'"
                  class="control-btn close-btn"
              >
                <el-icon><Close /></el-icon>
              </el-button>
            </div>
          </div>

          <!-- 图片展示区域 -->
          <div class="image-display-area">
            <!-- 左侧导航按钮 -->
            <Transition name="nav-fade">
              <button
                  v-if="imageViewerStore.canGoPrevious"
                  class="nav-button nav-prev"
                  @click="imageViewerStore.goToPrevious"
                  v-tooltip="'上一张 (←)'"
              >
                <el-icon><ArrowLeft /></el-icon>
              </button>
            </Transition>

            <!-- 主图片容器 -->
            <div
                class="image-container"
                @click.stop
                @wheel="handleWheel"
                @mousedown="handleMouseDown"
                @mousemove="handleMouseMove"
                @mouseup="handleMouseUp"
                @mouseleave="handleMouseUp"
                @touchstart="handleTouchStart"
                @touchmove="handleTouchMove"
                @touchend="handleTouchEnd"
                :class="{ 'dragging': imageViewerStore.isDragging, 'zoomedIn': imageViewerStore.isZoomedIn }"
                ref="imageContainer"
            >
              <div class="image-wrapper">
                <el-image
                    :src="currentImageUrl"
                    fit="contain"
                    class="main-image"
                    :loading="imageLoading"
                    @load="onImageLoad"
                    @error="onImageError"
                    lazy
                    :style="{
                                        transform: imageViewerStore.imageTransform,
                                        transformOrigin: 'center center',
                                        transition: imageViewerStore.isDragging ? 'none' : 'transform 0.3s ease'
                                    }"
                    @dragstart.prevent
                >
                  <template #placeholder>
                    <div class="image-loading">
                      <el-icon class="loading-icon">
                        <Loading />
                      </el-icon>
                      <p>正在加载图片...</p>
                    </div>
                  </template>
                  <template #error>
                    <div class="image-error">
                      <el-icon class="error-icon">
                        <Picture />
                      </el-icon>
                      <p>图片加载失败</p>
                      <el-button
                          type="primary"
                          size="small"
                          @click="retryLoadImage"
                      >
                        重新加载
                      </el-button>
                    </div>
                  </template>
                </el-image>
              </div>
            </div>

            <!-- 右侧导航按钮 -->
            <Transition name="nav-fade">
              <button
                  v-if="imageViewerStore.canGoNext"
                  class="nav-button nav-next"
                  @click="imageViewerStore.goToNext"
                  v-tooltip="'下一张 (→)'"
              >
                <el-icon><ArrowRight /></el-icon>
              </button>
            </Transition>
            <!-- 缩放提示 -->
            <div class="zoom-hint" v-if="!imageViewerStore.isZoomedIn">
              <span>滚轮缩放 • 双击放大</span>
            </div>

            <!-- 拖拽提示 -->
            <div class="drag-hint" v-if="imageViewerStore.isZoomedIn && !imageViewerStore.isDragging">
              <span>拖拽移动图片</span>
            </div>
          </div>

          <!-- 底部分页控制 -->
          <div v-if="imageViewerStore.hasMultiplePages" class="pagination-controls">
            <div class="page-info">
              <span class="current-page">{{ imageViewerStore.currentPage + 1 }}</span>
              <span class="page-separator">/</span>
              <span class="total-pages">{{ imageViewerStore.totalPages }}</span>
            </div>

            <!-- 缩略图预览 -->
            <div class="thumbnail-strip">
              <div
                  v-for="page in imageViewerStore.totalPages"
                  :key="page"
                  class="thumbnail-item"
                  :class="{ 'active': page - 1 === imageViewerStore.currentPage }"
                  @click="imageViewerStore.goToPage(page - 1)"
              >
                <el-image
                    :src="getThumbnailUrl(page - 1)"
                    fit="cover"
                    class="thumbnail-image"
                />
                <div class="thumbnail-overlay">
                  <span>{{ page }}</span>
                </div>
              </div>
            </div>
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
                  <span>上一张</span>
                </div>
                <div class="shortcut-item">
                  <kbd>→</kbd>
                  <span>下一张</span>
                </div>
                <div class="shortcut-item">
                  <kbd>+ -</kbd>
                  <span>放大/缩小</span>
                </div>
                <div class="shortcut-item">
                  <kbd>0</kbd>
                  <span>适应屏幕</span>
                </div>
                <div class="shortcut-item">
                  <kbd>R</kbd>
                  <span>重置位置</span>
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
  Link,
  Download,
  Close,
  ArrowLeft,
  ArrowRight,
  Loading,
  Picture,
  QuestionFilled,
  Plus,
  Minus,
  FullScreen
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useImageViewerStore } from '../assets/stores/PicPreview.js'
import { DownloadByPid } from "../../bindings/main/internal/pixivlib/ctl.js"
import axios from "axios";

const imageViewerStore = useImageViewerStore()

const modalRef = ref(null)
const imageContainer = ref(null)
const imageLoading = ref(true)
const showShortcuts = ref(false)

// 触摸相关状态
const touchStartDistance = ref(0)
const lastTouchTime = ref(0)

// 计算属性
const currentImageUrl = computed(() => {
  if (!imageViewerStore.currentPid) return ''

  const baseUrl = 'http://127.0.0.1:7234/api/preview?url='
  const pixivUrl = `${imageViewerStore.currentImageUrl}`
  return `${baseUrl}${pixivUrl}`
})

// 获取缩略图URL
const getThumbnailUrl = (page) => {
  const baseUrl = 'http://127.0.0.1:7234/api/preview?url='
  const pixivUrl = `${imageViewerStore.thumbUrls[page]}`
  return `${baseUrl}${pixivUrl}`
}

// 鼠标事件处理
const handleMouseDown = (event) => {
  if (event.button === 0) { // 左键
    imageViewerStore.startDrag(event.clientX, event.clientY)
    event.preventDefault()
  }
}

const handleMouseMove = (event) => {
  if (imageViewerStore.isDragging) {
    imageViewerStore.onDrag(event.clientX, event.clientY)
    event.preventDefault()
  }
}

const handleMouseUp = () => {
  imageViewerStore.endDrag()
}

// 双击放大
const handleDoubleClick = (event) => {
  if (imageViewerStore.scale === 1) {
    imageViewerStore.setZoom(2)
  } else {
    imageViewerStore.fitToScreen()
  }
}

// 滚轮缩放
const handleWheel = (event) => {
  const rect = imageContainer.value.getBoundingClientRect()
  const centerX = event.clientX - rect.left
  const centerY = event.clientY - rect.top
  imageViewerStore.handleWheel(event, centerX, centerY)
}

// 触摸事件处理
const handleTouchStart = (event) => {
  if (event.touches.length === 1) {
    const now = Date.now()
    const touch = event.touches[0]

    // 检测双击
    if (now - lastTouchTime.value < 300) {
      handleDoubleClick()
    } else {
      imageViewerStore.startDrag(touch.clientX, touch.clientY)
    }

    lastTouchTime.value = now
  } else if (event.touches.length === 2) {
    // 双指缩放
    const touch1 = event.touches[0]
    const touch2 = event.touches[1]
    touchStartDistance.value = Math.sqrt(
        Math.pow(touch2.clientX - touch1.clientX, 2) +
        Math.pow(touch2.clientY - touch1.clientY, 2)
    )
  }
  event.preventDefault()
}

const handleTouchMove = (event) => {
  if (event.touches.length === 1) {
    const touch = event.touches[0]
    imageViewerStore.onDrag(touch.clientX, touch.clientY)
  } else if (event.touches.length === 2) {
    // 双指缩放
    const touch1 = event.touches[0]
    const touch2 = event.touches[1]
    const currentDistance = Math.sqrt(
        Math.pow(touch2.clientX - touch1.clientX, 2) +
        Math.pow(touch2.clientY - touch1.clientY, 2)
    )

    if (touchStartDistance.value > 0) {
      const scaleRatio = currentDistance / touchStartDistance.value
      const newScale = imageViewerStore.scale * scaleRatio
      imageViewerStore.setZoom(newScale)
      touchStartDistance.value = currentDistance
    }
  }
  event.preventDefault()
}

const handleTouchEnd = () => {
  imageViewerStore.endDrag()
  touchStartDistance.value = 0
}



// 事件处理
const handleBackdropClick = (event) => {
  if (event.target === event.currentTarget) {
    imageViewerStore.closeViewer()
  }
}

const handleKeyDown = (event) => {
  imageViewerStore.handleKeyPress(event)
}

const onImageLoad = () => {
  imageLoading.value = false
}

const onImageError = () => {
  imageLoading.value = false
  ElMessage.error('图片加载失败')
}

const retryLoadImage = () => {
  imageLoading.value = true
  // 重新触发图片加载
  const img = document.querySelector('.main-image img')
  if (img) {
    img.src = currentImageUrl.value
  }
}

const downloadCurrent = async () => {
  try {
    const success = await DownloadByPid(imageViewerStore.currentPid)
    if (success) {
      ElMessage.success('已添加到下载队列')
    } else {
      ElMessage.error('下载失败')
    }
  } catch (error) {
    console.error('下载失败:', error)
    ElMessage.error('下载失败，请稍后重试')
  }
}

const openInPixiv = () => {
  const url = `https://www.pixiv.net/artworks/${imageViewerStore.currentPid}`
  window.open(url, '_blank')
}

const toggleShortcuts = () => {
  showShortcuts.value = !showShortcuts.value
}

// 监听键盘事件
onMounted(() => {
  document.addEventListener('keydown', imageViewerStore.handleKeyPress)

  if (modalRef.value) {
    modalRef.value.focus()
  }

  // 添加双击事件
  if (imageContainer.value) {
    imageContainer.value.addEventListener('dblclick', handleDoubleClick)
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', imageViewerStore.handleKeyPress)
})

// 监听当前页变化，重置图片加载状态
watch(() => imageViewerStore.currentPage, () => {
  imageLoading.value = true
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

// 主容器
.image-viewer-modal {
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
    max-width: 1400px;
    max-height: 900px;
    display: flex;
    flex-direction: column;
    //background: rgba(255, 255, 255, 0.95);
    border-radius: 20px;
    overflow: hidden;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  }
}

// 头部信息栏
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 30px;
  //background: transparent;
  background: rgba(0, 0, 0, 0.57);
  color: white;

  .image-info {
    flex: 1;

    .image-title {
      margin: 0 0 8px 0;
      font-size: 20px;
      font-weight: 600;
      line-height: 1.3;
    }

    .image-meta {
      display: flex;
      align-items: center;
      gap: 20px;
      font-size: 14px;
      opacity: 0.9;

      .author-info,
      .pid-info {
        display: flex;
        align-items: center;
        gap: 6px;
      }
    }
  }

  .modal-controls {
    display: flex;
    gap: 10px;

    .zoom-controls {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 8px 12px;
      background: rgba(255, 255, 255, 0.1);
      border-radius: 8px;

      .zoom-info {
        font-size: 14px;
        font-weight: 600;
        min-width: 50px;
        text-align: center;
      }

      .zoom-buttons {
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

// 图片展示区域
.image-display-area {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 60%;
  //background: rgba(255, 255, 255, 0.95);
  background: rgba(0, 0, 0, 0.57);
  //background: #f8f9fa;

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

  .image-container {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px 120px;

    .image-wrapper {
      //max-width: 100%;
      height: 95%;
      .main-image {
        height: 100%;
        border-radius: 8px;
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
      }
    }
  }

  .image-loading,
  .image-error {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 15px;
    color: #ffffff;

    .loading-icon {
      font-size: 48px;
      animation: spin 1s linear infinite;
    }

    .error-icon {
      font-size: 48px;
    }

    p {
      font-size: 16px;
      margin: 0;
    }
  }
  // 提示文本
  .zoom-hint,
  .drag-hint {
    position: absolute;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 8px 16px;
    border-radius: 20px;
    font-size: 14px;
    pointer-events: none;
    z-index: 5;
    animation: fadeInOut 3s ease-in-out;
  }
}

// 分页控制
.pagination-controls {
  display: flex;
  flex-direction: column;
  gap: 15px;
  padding: 10px 30px;
  //background: white;
  background: rgba(0, 0, 0, 0.57);
  //border-top: 1px solid #e4e7ed;

  .page-info {
    text-align: center;
    font-size: 16px;
    font-weight: 600;
    color: #303133;

    .current-page {
      color: #409EFF;
    }

    .page-separator {
      margin: 0 8px;
      color: #909399;
    }
  }
  .thumbnail-strip {
    display: flex;
    gap: 8px;
    overflow-x: auto;
    padding: 10px 0;

    .thumbnail-item {
      position: relative;
      flex-shrink: 0;
      width: 60px;
      height: 60px;
      border-radius: 8px;
      overflow: hidden;
      cursor: pointer;
      border: 2px solid transparent;
      transition: all 0.3s ease;

      &:hover {
        transform: scale(1.05);
        border-color: #409EFF;
      }

      &.active {
        border-color: #409EFF;
        box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.3);
      }

      .thumbnail-image {
        width: 100%;
        height: 100%;
      }

      .thumbnail-overlay {
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        background: rgba(0, 0, 0, 0.7);
        color: white;
        font-size: 12px;
        text-align: center;
        padding: 2px;
      }
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

    .image-title {
      font-size: 18px;
    }

    .image-meta {
      gap: 15px;
      font-size: 13px;
    }

    .modal-controls {
      gap: 8px;
      .zoom-controls {
        padding: 6px 10px;

        .zoom-info {
          font-size: 12px;
          min-width: 40px;
        }
      }
    }
  }

  .image-display-area {
    .image-container {
      padding: 15px 80px;
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

    .zoom-hint,
    .drag-hint {
      font-size: 12px;
      padding: 6px 12px;
      bottom: 15px;
    }
  }

  .pagination-controls {
    padding: 15px 20px;

    .thumbnail-strip {
      .thumbnail-item {
        width: 50px;
        height: 50px;
      }
    }
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

@keyframes fadeInOut {
  0%, 100% { opacity: 0; }
  20%, 80% { opacity: 1; }
}
</style>