<template>
  <div class="search-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">
        <el-icon class="title-icon"><Search /></el-icon>
        以图搜图
      </h1>
      <p class="page-description">上传图片以查找相似的图片来源</p>
      <p class="page-tip">
        <el-icon><InfoFilled /></el-icon>
        支持拖拽、点击上传，或使用 <kbd>Ctrl+V</kbd> (Mac: <kbd>⌘+V</kbd>) 粘贴剪切板图片
      </p>
    </div>

    <el-card class="upload-card">
      <template #header>
        <div class="card-header">
          <el-icon><Upload /></el-icon>
          <span>图片上传</span>
        </div>
      </template>

      <div class="upload-area">
        <el-upload
            class="image-uploader"
            :auto-upload="false"
            :show-file-list="false"
            :on-change="handleImageSelect"
            accept="image/jpeg,image/png"
            drag
        >
          <div v-if="!previewUrl" class="upload-content">
            <el-icon class="upload-icon"><Plus /></el-icon>
            <div class="upload-text">
              <p>点击或拖拽图片到此处</p>
              <span>支持 JPG、PNG 格式，文件小于 10MB</span>
              <span class="paste-hint">或按 <kbd>Ctrl+V</kbd> 粘贴剪切板图片</span>
            </div>
          </div>
          <div v-else class="preview-container">
            <img :src="previewUrl" class="preview-image" alt="预览图片" height="200px" />
            <div class="preview-overlay">
              <el-icon class="change-icon"><RefreshRight /></el-icon>
              <span>点击更换图片</span>
            </div>
          </div>
        </el-upload>

        <div class="upload-actions">
          <el-button
              type="primary"
              class="search-btn"
              :loading="searching"
              :disabled="!selectedFile"
              @click="handleSearch"
          >
            <el-icon v-if="!searching"><Search /></el-icon>
            {{ searching ? '搜索中...' : '开始搜索' }}
          </el-button>
          <el-button v-if="selectedFile" @click="handleClear">
            <el-icon><Delete /></el-icon>
            清除
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 搜索结果 -->
    <div v-if="searchResults.length > 0" class="results-section">
      <div class="results-header">
        <h2>
          <el-icon><Collection /></el-icon>
          SauceNAO搜索结果
          <el-tag type="info" class="count-tag">{{ searchResults.length }} 条结果</el-tag>
        </h2>
        <div v-if="searchHeader" class="search-info">
          <el-tag>剩余短期搜索: {{ searchHeader.short_remaining }}/{{ searchHeader.short_limit }}</el-tag>
          <el-tag type="success">剩余长期搜索: {{ searchHeader.long_remaining }}/{{ searchHeader.long_limit }}</el-tag>
        </div>
      </div>

      <div class="results-grid">
        <el-card
            v-for="(result, index) in searchResults"
            :key="index"
            class="result-card"
            :class="{ 'pixiv-result': isPixivResult(result) }"
        >
          <!-- 相似度标签 -->
          <div class="similarity-badge">
            <el-tag
                :type="getSimilarityConfig(result.header.similarity).type"
                size="large"
            >
              {{ result.header.similarity }}%
            </el-tag>
          </div>

          <!-- 缩略图 -->
          <div class="result-thumbnail" @click="openOriginal(result)">
            <img :src="result.header.thumbnail" :alt="getResultTitle(result)" />
            <div class="thumbnail-overlay" v-if="result.data.ext_urls && result.data.ext_urls.length > 0">
              <el-icon class="preview-icon"><View /></el-icon>
              <span>查看原图</span>
            </div>
          </div>

          <!-- 结果信息 -->
          <div class="result-info">
            <!-- 来源标识 -->
            <div class="result-source">
              <el-icon><Platform /></el-icon>
              <span>{{ getSourceName(result.header.index_name) }}</span>
            </div>

            <!-- Pixiv 特殊标识 -->
            <div v-if="isPixivResult(result)" class="pixiv-info">
              <el-tag style="cursor: pointer" type="primary" effect="dark" @click="openPixivArtwork(extractPixivId(result))">
                <el-icon><Picture /></el-icon>
                Pixiv: {{ extractPixivId(result) }}
              </el-tag>
            </div>

            <!-- 标题 -->
            <h3 v-if="result.data.title" class="result-title">
              {{ result.data.title }}
            </h3>
            <h3 v-else-if="result.data.source" class="result-title">
              {{ result.data.source }}
            </h3>

            <!-- 作者信息 -->
            <div v-if="result.data.creator || result.data.author_name || result.data.member_name" class="author-info">
              <el-icon><User /></el-icon>
              <span>{{ result.data.creator || result.data.author_name || result.data.member_name }}</span>
            </div>

            <!-- 角色/系列信息 -->
            <div v-if="result.data.material" class="metadata-item">
              <el-icon><StarFilled /></el-icon>
              <span>{{ result.data.material }}</span>
            </div>
            <div v-if="result.data.characters" class="metadata-item characters">
              <el-icon><Avatar /></el-icon>
              <span>{{ result.data.characters }}</span>
            </div>

            <!-- 外部链接 -->
            <div v-if="result.data.ext_urls && result.data.ext_urls.length > 0" class="external-links">
              <el-divider />
              <div class="links-title">
                <el-icon><Link /></el-icon>
                <span>相关链接</span>
              </div>
              <div class="links-list">
                <el-button
                    v-for="(url, urlIndex) in result.data.ext_urls"
                    :key="urlIndex"
                    type="primary"
                    link
                    @click="OpenInBrowser(url)"
                >
                  <el-icon><TopRight /></el-icon>
                  {{ formatUrl(url) }}
                </el-button>
              </div>
            </div>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="!searching && hasSearched" class="empty-state">
      <el-empty description="未找到相似的图片">
        <el-button type="primary" @click="handleClear">重新搜索</el-button>
      </el-empty>
    </div>
  </div>
</template>

<script setup>
import {ref, onMounted, onUnmounted} from 'vue'
import {ElMessage, ElNotification} from 'element-plus'
import {
  Search,
  Upload,
  Plus,
  RefreshRight,
  Delete,
  Collection,
  View,
  Platform,
  Picture,
  User,
  StarFilled,
  Avatar,
  Link,
  TopRight,
  InfoFilled
} from '@element-plus/icons-vue'
import {
  fetchSearch,
  extractPixivId,
  isPixivResult,
  getResultTitle,
  getSourceName,
  getSimilarityConfig,
  formatUrl
} from '../assets/js/utils/saucenao.js'
import {OpenInBrowser} from "../../bindings/main/internal/pixivlib/ctl.js";
import {openPixivArtwork} from "../assets/js/utils/index.js";

const selectedFile = ref(null)
const previewUrl = ref('')
const searching = ref(false)
const hasSearched = ref(false)
const searchResults = ref([])
const searchHeader = ref(null)

const handleImageSelect = (file) => {
  const isImage = file.raw.type === 'image/jpeg' || file.raw.type === 'image/png'
  if (!isImage) {
    ElMessage.error('只能上传 JPG/PNG 格式的图片!')
    return
  }

  const isLt10M = file.raw.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('图片大小不能超过 10MB!')
    return
  }

  selectedFile.value = file.raw
  previewUrl.value = URL.createObjectURL(file.raw)
}

// 处理剪切板粘贴
const handlePaste = (e) => {
  const items = e.clipboardData?.items
  if (!items) return

  for (let i = 0; i < items.length; i++) {
    const item = items[i]

    if (item.type.indexOf('image') !== -1) {
      const file = item.getAsFile()
      if (file) {
        if (file.type !== 'image/jpeg' && file.type !== 'image/png') {
          ElMessage.warning('剪切板中的图片格式不支持，请使用 JPG 或 PNG 格式')
          return
        }

        if (file.size > 10 * 1024 * 1024) {
          ElMessage.warning('剪切板中的图片大小超过 10MB')
          return
        }

        selectedFile.value = file
        previewUrl.value = URL.createObjectURL(file)

        ElMessage.success('已从剪切板粘贴图片')
        e.preventDefault()
        return
      }
    }
  }
}

onMounted(() => {
  document.addEventListener('paste', handlePaste)
})

onUnmounted(() => {
  document.removeEventListener('paste', handlePaste)
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
  }
})

const handleSearch = async () => {
  if (!selectedFile.value) {
    ElMessage.warning('请先选择图片')
    return
  }

  searching.value = true
  hasSearched.value = true

  try {
    const result = await fetchSearch(selectedFile.value)

    if (result && result.header && result.results) {
      searchHeader.value = result.header
      searchResults.value = result.results

      ElNotification({
        type: 'success',
        title: '搜索完成',
        message: `找到 ${result.results.length} 条相似结果`,
        position: 'bottom-right',
        duration: 3000,
      })
    } else {
      searchResults.value = []
      ElMessage.warning('未找到相似的图片')
    }
  } catch (error) {
    console.error('搜索失败:', error)
    ElNotification({
      type: 'error',
      title: '搜索失败',
      message: error.message || '发生未知错误',
      position: 'bottom-right',
      duration: 5000,
    })
  } finally {
    searching.value = false
  }
}

const handleClear = () => {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value)
  }

  selectedFile.value = null
  previewUrl.value = ''
  searchResults.value = []
  searchHeader.value = null
  hasSearched.value = false
}

const openOriginal = (result) => {
  if (result.data.ext_urls && result.data.ext_urls.length > 0) {
    OpenInBrowser(result.data.ext_urls[0])
  }
}
</script>

<style scoped lang="less">
@import "../assets/style/common/search.less";
</style>