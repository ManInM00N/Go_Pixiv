<template>
  <div class="novel-card-wrapper">
    <el-skeleton
        class="card-skeleton"
        :loading="loading"
        animated
        :throttle="300"
    >
      <template #template>
        <div class="skeleton-card">
          <el-skeleton-item variant="image" class="skeleton-cover" />
          <div class="skeleton-content">
            <el-skeleton-item variant="h3" style="width: 80%; margin-bottom: 10px;" />
            <el-skeleton-item variant="text" style="width: 60%; margin-bottom: 8px;" />
            <el-skeleton-item variant="text" style="width: 90%; margin-bottom: 8px;" />
            <el-skeleton-item variant="text" style="width: 70%;" />
          </div>
        </div>
      </template>

      <template #default>
        <el-card :body-style="{ padding: '0px', marginBottom: '1px', width: '100%' }"
            class="novel-card"
            :class="{ 'login-required': props.isLoginOnly }"
            shadow="hover"
            @click="openNovelViewer"
        >
          <!-- 封面区域 -->
          <div class="cover-container">
            <LazyImg
                :ref="pic"
                :onload="load = false"
                :url="getImageUrl(props.cover, props.genre == '1')"
                class="image"

            />

            <!-- 遮罩层 -->
            <div class="cover-overlay">
              <div class="overlay-content">
                <el-icon class="read-icon"><View /></el-icon>
                <span class="read-text">阅读小说</span>
              </div>
            </div>

            <!-- 标签区域 -->
            <div class="cover-badges">
              <el-tag
                  v-if="props.isLoginOnly"
                  type="warning"
                  size="small"
                  class="login-badge"
              >
                登录可见
              </el-tag>
              <el-tag
                  v-if="props.genre"
                  type="danger"
                  size="small"
                  class="login-badge"
              >
                R18
              </el-tag>
              <el-tag
                  v-if="props.aiType === 0"
                  type="info"
                  size="small"
                  class="ai-badge"
              >
                AI
              </el-tag>

              <el-tag
                  v-if="props.seriesNavData && props.seriesNavData.seriesType"
                  type="success"
                  size="small"
                  class="series-badge"
              >
                连载
              </el-tag>
            </div>
          </div>

          <!-- 内容区域 -->
          <div class="novel-content">
            <!-- 标题 -->
            <div class="novel-title">
              <el-tooltip
                  placement="top"
                  :disabled="!isTitleLong"
              >
                <template #content>
                  <p style="max-width:500px;">{{props.title}}</p>
                </template>
                <h3 class="title-text">{{ props.title }}</h3>
              </el-tooltip>
            </div>

            <!-- 作者信息 -->
            <div class="author-info">
              <el-avatar
                  :size="20"
                  class="author-avatar"
                  @click.stop="openPixivUser(props.userId)"
              >
<!--                <el-icon>-->
                  <LazyImg
                      :onload="load = false"
                      :url="getProfileUrl(props.profileImageUrl)"
                  >

                  </LazyImg>
<!--                </el-icon>-->
              </el-avatar>
              <span
                  class="author-name"
                  @click.stop="openPixivUser(props.userId)"
              >
                                {{ props.userName }}
                            </span>
            </div>

            <!-- 简介 -->
            <div class="novel-description">
              <el-tooltip
                  placement="top"
                  :disabled="!isDescriptionLong"
              >
                <template #content>
                  <p style="max-width:500px;">{{props.description}}</p>
                </template>
                <p class="description-text">{{ props.description }}</p>
              </el-tooltip>
            </div>

            <!-- 统计信息 -->
            <div class="novel-stats">
              <div class="stat-item">
                <el-icon class="stat-icon"><Reading /></el-icon>
                <span class="stat-text">{{ formatCount(props.textCount) }} 字</span>
              </div>
              <div class="stat-item">
                <el-icon class="stat-icon"><Star /></el-icon>
                <span class="stat-text">{{ formatCount(props.bookmarkCount) }}</span>
              </div>
            </div>

            <!-- 类型标签 -->
            <div class="novel-genre">

              <el-tag v-for="item in props.tags" type="primary" size="small" class="genre-tag">
                {{ item }}
              </el-tag>
            </div>

            <!-- 系列信息 -->
            <div class="series-info" v-if="props.seriesNavData">
              <div class="series-title">
                <el-icon><Collection /></el-icon>
                <span>{{ props.seriesNavData.title }}</span>
              </div>
            </div>

            <!-- 操作区域 -->
            <div class="novel-actions">
              <el-button
                  type="primary"
                  size="small"
                  @click.stop="openNovelViewer"
                  class="read-btn"
              >
                <el-icon><View /></el-icon>
                阅读
              </el-button>

              <el-button
                  size="small"
                  @click.stop="downloadNovel"
                  :loading="downloading"
                  class="download-btn"
              >
                <el-icon><Download /></el-icon>
              </el-button>

              <el-button
                  size="small"
                  @click.stop="toggleFavorite"
                  :type="isFavorited ? 'danger' : ''"
                  class="favorite-btn"
              >
                <el-icon>
                  <Star v-if="isFavorited" />
                  <StarFilled v-else />
                </el-icon>
              </el-button>
            </div>
          </div>
        </el-card>
      </template>
    </el-skeleton>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import {
  Loading,
  Document,
  View,
  User,
  Reading,
  Star,
  StarFilled,
  Download,
  Collection
} from '@element-plus/icons-vue'
import {
  getImageUrl,
  getProfileUrl,
  formatCount,
  isTextLong,
  openPixivUser
} from '../assets/js/utils/index.js'

import { ElMessage } from 'element-plus'
import {form} from "../assets/js/configuration.js";
import noProfileImg from "../assets/images/NoR18.png";
import {LazyImg} from "vue-waterfall-plugin-next";
import { useNovelViewerStore } from '../assets/stores/novelViewer.js'
import {DownloadByNovelId} from "../../bindings/main/internal/pixivlib/ctl.js";

const novelViewerStore = useNovelViewerStore()
const name = "NovelCard"
const pic = ref(null)


const openNovelViewer = () => {
  // 构建小说数据对象
  const novelData = {
    id: props.id,
    title: props.title,
    userName: props.userName,
    userId: props.userId,
    genre: props.genre,
    aiType: props.aiType,
    seriesNavData: props.seriesNavData,
    description: props.description,
    cover: props.cover,
    bookmarkCount: props.bookmarkCount,
    textCount: props.textCount,
    tags : props.tags,
  }

  // 使用 store 打开查看器
  novelViewerStore.openViewer(novelData)
}

// Props 定义
const props = defineProps({
  content: {
    type: String,
    default: ""
  },
  profileImageUrl:{
     type: String,
    default: "",
  },
  cover: {
    type: String,
    default: ""
  },
  id: {
    type: [String, Number],
    required: true
  },
  userId: {
    type: [String, Number],
    required: true
  },
  userName: {
    type: String,
    required: true
  },
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    default: ""
  },
  page: {
    type: Number,
    default: 1
  },
  bookmarkCount: {
    type: Number,
    default: 0
  },
  textCount: {
    type: Number,
    default: 0
  },
  wordCount: {
    type: Number,
    default: 0
  },
  isLoginOnly: {
    type: Boolean,
    default: false
  },
  genre: {
    type: String,
    default: ""
  },
  aiType: {
    type: Number,
    default: false
  },
  seriesNavData: {
    type: Object,
    default: null
  },
  tags: {
    type: Object,
    default: null
  }
})
// 响应式数据
const loading = ref(false)
const coverLoading = ref(true)
const downloading = ref(false)
const isFavorited = ref(false)


const isTitleLong = computed(() => isTextLong(props.title, 30))
const isDescriptionLong = computed(() => isTextLong(props.description, 100))

// 事件处理
const onCoverLoad = () => {
  coverLoading.value = false
}

const onCoverError = () => {
  coverLoading.value = false
}

const downloadNovel = async () => {
  try {
    downloading.value = true
    // 这里调用小说下载函数
    await DownloadByNovelId(props.id.toString(),false)
    ElMessage.success(`开始下载小说: ${props.title}`)
  } catch (error) {
    console.error('下载失败:', error)
    ElMessage.error('下载失败，请稍后重试')
  } finally {
    setTimeout(() => {
      downloading.value = false
    }, 1000)
  }
}

const toggleFavorite = () => {
  isFavorited.value = !isFavorited.value
  ElMessage.success(isFavorited.value ? '已添加收藏' : '已取消收藏')
}

</script>


<style lang="less" scoped>
@import "../assets/style/common/loading.less";
@import "../assets/style/common/cards.less";
@import "../assets/style/common/waterfall.less";
@import "../assets/style/common/buttons.less";
@import "../assets/style/common/animations.less";

// 组件特定样式
.novel-card-wrapper {
  width: 100%;
  max-width: 300px;
  margin: 0 auto;
}

// 小说卡片特定样式
.novel-card {
  width: 100%;
  border: 2px solid transparent;
  cursor: pointer;

  &.login-required {
    border-color: #f39c12;

    &:hover {
      border-color: #e67e22;
    }
  }
}


// 简介
.novel-description {
  margin-bottom: 12px;

  .description-text {
    margin: 0;
    font-size: 13px;
    color: #909399;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
}

// 统计信息
.novel-stats {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;

  .stat-item {
    display: flex;
    align-items: center;
    gap: 4px;

    .stat-icon {
      font-size: 14px;
      color: #909399;
    }

    .stat-text {
      font-size: 12px;
      color: #909399;
    }
  }
}

// 类型标签
.novel-genre {
  margin-bottom: 10px;

  .genre-tag {
    font-size: 11px;
    padding: 2px 8px;
    border-radius: 10px;
    white-space: normal;
    height: auto;
  }
}

// 系列信息
.series-info {
  margin-bottom: 12px;
  padding: 8px;
  background: #f8f9fa;
  border-radius: 8px;

  .series-title {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 13px;
    color: #606266;
    margin-bottom: 4px;
  }

  .series-order {
    font-size: 12px;
    color: #909399;
  }
}

// 操作区域
.novel-actions {
  display: flex;
  gap: 8px;

  .read-btn {
    flex: 1;
    font-size: 12px;
  }

  .download-btn,
  .favorite-btn {
    width: 32px;
    height: 32px;
    padding: 0;

    .el-icon {
      font-size: 14px;
    }
  }

  .favorite-btn {
    &.el-button--danger {
      background: #f56c6c;
      border-color: #f56c6c;
      color: white;
    }
  }
}

// 响应式
@media (max-width: 480px) {
  .novel-card-wrapper {
    max-width: 100%;
  }

  .novel-content {
    padding: 12px;
  }

  .novel-title .title-text {
    font-size: 15px;
  }

  .novel-actions {
    gap: 6px;

    .download-btn,
    .favorite-btn {
      width: 28px;
      height: 28px;
    }
  }
}
</style>

<style lang="less">
.el-tag {
  white-space: normal;
  height: auto;
}
</style>
