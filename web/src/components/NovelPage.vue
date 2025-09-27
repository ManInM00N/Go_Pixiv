<template>
  <div class="novel-page-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">
        <el-icon class="title-icon"><Reading /></el-icon>
        小说中心
      </h1>
      <p class="page-description">发现精彩的文学作品</p>
    </div>

    <!-- 导航标签 -->
    <el-card class="nav-card" shadow="hover">
      <el-tabs
          v-model="activeTab"
          type="card"
          class="novel-tabs"
          @tab-change="handleTabChange"
      >
        <!-- 关注小说 -->
        <el-tab-pane label="关注作者" name="follow">
          <template #label>
            <div class="tab-label">
              <el-icon><StarFilled /></el-icon>
              <span>关注作者</span>
              <el-badge
                  v-if="followNovels.length > 0"
                  :value="followNovels.length"
                  :max="99"
                  class="tab-badge"
              />
            </div>
          </template>
        </el-tab-pane>

        <!-- 排行榜 -->
        <el-tab-pane label="排行榜" name="ranking">
          <template #label>
            <div class="tab-label">
              <el-icon><Trophy /></el-icon>
              <span>排行榜</span>
              <el-badge
                  v-if="rankingNovels.length > 0"
                  :value="rankingNovels.length"
                  :max="99"
                  class="tab-badge"
              />
            </div>
          </template>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 筛选控制 -->
    <el-card class="filter-card" shadow="hover">
      <div class="filter-header">
        <el-icon><Filter /></el-icon>
        <span>筛选选项</span>
      </div>

      <div class="filter-content">
        <el-row :gutter="20">
          <!-- 排行榜类型（仅排行榜标签显示） -->
          <el-col :xs="24" :sm="12" :md="6" v-if="activeTab === 'ranking'">
            <div class="filter-item">
              <label class="filter-label">排行榜类型</label>
              <el-select
                  v-model="rankingType"
                  placeholder="选择类型"
                  @change="fetchRankingNovels"
                  class="filter-select"
              >
                <el-option label="每日排行" value="daily" />
                <el-option label="每周排行" value="weekly" />
                <el-option label="每月排行" value="monthly" />
                <el-option label="新人排行" value="rookie" />
                <el-option label="完结排行" value="complete" />
              </el-select>
            </div>
          </el-col>

          <!-- 分类筛选 -->
          <el-col :xs="24" :sm="12" :md="6">
            <div class="filter-item">
              <label class="filter-label">作品分类</label>
              <el-select
                  v-model="selectedGenre"
                  placeholder="全部分类"
                  @change="applyFilters"
                  class="filter-select"
                  clearable
              >
                <el-option label="全部" value="" />
                <el-option label="恋爱" value="恋爱" />
                <el-option label="幻想" value="幻想" />
                <el-option label="现实" value="现实" />
                <el-option label="科幻" value="科幻" />
                <el-option label="历史" value="历史" />
                <el-option label="悬疑" value="悬疑" />
                <el-option label="同人" value="同人" />
              </el-select>
            </div>
          </el-col>

          <!-- 字数筛选 -->
          <el-col :xs="24" :sm="12" :md="6">
            <div class="filter-item">
              <label class="filter-label">字数范围</label>
              <el-select
                  v-model="selectedWordCount"
                  placeholder="不限"
                  @change="applyFilters"
                  class="filter-select"
                  clearable
              >
                <el-option label="不限" value="" />
                <el-option label="1万字以下" value="short" />
                <el-option label="1-5万字" value="medium" />
                <el-option label="5-20万字" value="long" />
                <el-option label="20万字以上" value="epic" />
              </el-select>
            </div>
          </el-col>

          <!-- 状态筛选 -->
          <el-col :xs="24" :sm="12" :md="6">
            <div class="filter-item">
              <label class="filter-label">作品状态</label>
              <el-select
                  v-model="selectedStatus"
                  placeholder="全部状态"
                  @change="applyFilters"
                  class="filter-select"
                  clearable
              >
                <el-option label="全部" value="" />
                <el-option label="连载中" value="ongoing" />
                <el-option label="已完结" value="completed" />
                <el-option label="AI创作" value="ai" />
                <el-option label="需要登录" value="login" />
              </el-select>
            </div>
          </el-col>
        </el-row>

        <!-- 搜索框 -->
        <div class="search-section">
          <el-input
              v-model="searchKeyword"
              placeholder="搜索小说标题、作者、简介..."
              class="search-input"
              @input="handleSearch"
              clearable
          >
            <template #prepend>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </div>
    </el-card>

    <!-- 内容区域 -->
    <div class="content-area">
      <!-- 统计信息 -->
      <div class="stats-info">
        <el-tag type="info" size="large">
          <el-icon><DataLine /></el-icon>
          共找到 {{ filteredNovels.length }} 部小说
        </el-tag>

        <div class="view-options">
          <el-radio-group v-model="viewMode" size="small">
            <el-radio-button value="waterfall">瀑布流</el-radio-button>
            <el-radio-button value="list">列表</el-radio-button>
          </el-radio-group>
        </div>
      </div>

      <!-- 瀑布流展示 -->
      <div v-if="viewMode === 'waterfall'" class="waterfall-container">
        <Waterfall
            ref="waterfall"
            :list="filteredNovels"
            :width="300"
            :gutter="20"
            background-color="transparent"
            :animationEffect="fadeInUp"
            key="novelWaterfall"
        >
          <template #default="{ item, url, index }">
            <transition name="el-fade-in-linear">
              <div class="novel-wrapper">
                <NovelCard v-bind="item" :key="item.id" />
              </div>
            </transition>
          </template>
        </Waterfall>
      </div>

      <!-- 列表展示 -->
      <div v-else class="list-container">
        <div
            v-for="novel in filteredNovels"
            :key="novel.id"
            class="list-item"
        >
          <NovelCard v-bind="novel" />
        </div>
      </div>

      <!-- 加载更多 -->
      <div class="load-more-section" v-if="hasMore">
        <el-button
            type="primary"
            size="large"
            @click="loadMore"
            :loading="loading"
            class="load-more-btn"
        >
          <el-icon><ArrowDown /></el-icon>
          加载更多
        </el-button>
      </div>

      <!-- 空状态 -->
      <el-empty
          v-if="!loading && filteredNovels.length === 0"
          description="暂无小说数据"
          class="empty-state"
      >
        <el-button type="primary" @click="refreshData">
          <el-icon><Refresh /></el-icon>
          刷新数据
        </el-button>
      </el-empty>
    </div>

    <!-- 回到顶部 -->
    <el-backtop
        target=".novel-page-container"
        :bottom="60"
        :right="60"
        :visibility-height="400"
    >
      <div class="back-to-top">
        <el-icon><ArrowUp /></el-icon>
      </div>
    </el-backtop>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import {
  Reading,
  StarFilled,
  Trophy,
  Filter,
  Search,
  DataLine,
  ArrowDown,
  ArrowUp,
  Refresh
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { Waterfall } from 'vue-waterfall-plugin-next'
import NovelCard from './NovelCard.vue'
import axios from "axios";

// 响应式数据
const activeTab = ref('follow')
const loading = ref(false)
const viewMode = ref('waterfall')
const waterfall = ref(null)

// 筛选条件
const rankingType = ref('daily')
const selectedGenre = ref('')
const selectedWordCount = ref('')
const selectedStatus = ref('')
const searchKeyword = ref('')

// 数据
const followNovels = ref([])
const rankingNovels = ref([])
const hasMore = ref(true)
const currentPage = ref(1)

// 计算属性
const currentNovels = computed(() => {
  return activeTab.value === 'follow' ? followNovels.value : rankingNovels.value
})

const filteredNovels = computed(() => {
  let novels = [...currentNovels.value]

  // 分类筛选
  if (selectedGenre.value) {
    novels = novels.filter(novel => novel.genre === selectedGenre.value)
  }

  // 字数筛选
  if (selectedWordCount.value) {
    novels = novels.filter(novel => {
      const count = novel.characterCount
      switch (selectedWordCount.value) {
        case 'short': return count < 10000
        case 'medium': return count >= 10000 && count < 50000
        case 'long': return count >= 50000 && count < 200000
        case 'epic': return count >= 200000
        default: return true
      }
    })
  }

  // 状态筛选
  if (selectedStatus.value) {
    novels = novels.filter(novel => {
      switch (selectedStatus.value) {
        case 'ongoing': return novel.seriesNavData && novel.seriesNavData.next
        case 'completed': return !novel.seriesNavData || !novel.seriesNavData.next
        case 'ai': return novel.aiType
        case 'login': return novel.isLoginOnly
        default: return true
      }
    })
  }

  // 搜索筛选
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    novels = novels.filter(novel =>
        novel.title.toLowerCase().includes(keyword) ||
        novel.userName.toLowerCase().includes(keyword) ||
        novel.description.toLowerCase().includes(keyword)
    )
  }

  return novels
})

// 方法
const handleTabChange = (tabName) => {
  if (tabName === 'follow') {
    fetchFollowNovels()
  } else if (tabName === 'ranking') {
    fetchRankingNovels()
  }
}

const fetchFollowNovels = async () => {
  try {
    loading.value = true
    // 调用获取关注作者小说的函数
    const novels = await fetchNovel('follow', currentPage.value)

    if (currentPage.value === 1) {
      followNovels.value = novels
    } else {
      followNovels.value.push(...novels)
    }

    hasMore.value = novels.length > 0

    // 重新渲染瀑布流
    await nextTick()
    if (waterfall.value) {
      waterfall.value.renderer()
    }

  } catch (error) {
    console.error('获取关注小说失败:', error)
    ElMessage.error('获取数据失败，请稍后重试')

  } finally {
    loading.value = false
  }
}

const fetchRankingNovels = async () => {
  try {
    loading.value = true
    // 调用获取排行榜小说的函数

    // 获取小说数组
    const novels = await fetchNovel('ranking', currentPage.value, { type: rankingType.value })

    if (currentPage.value === 1) {
      rankingNovels.value = novels
    } else {
      rankingNovels.value.push(...novels)
    }

    hasMore.value = novels.length > 0

    // 重新渲染瀑布流
    await nextTick()
    if (waterfall.value) {
      waterfall.value.renderer()
    }

  } catch (error) {
    console.error('获取排行榜小说失败:', error)
    ElMessage.error('获取数据失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
const mode = ref("all")
const novelItems = ref([])

// 模拟 fetchNovel 函数（实际使用时替换为真实的API调用）
const fetchNovel = async (type, page = 1, options = {}) => {
  // 模拟网络延迟
  await new Promise(resolve => setTimeout(resolve, 1000))
  axios.get("http://127.0.0.1:7234/api/followpage", {
    params: {
      p: currentPage.value.toString(),
      types: "novel",
      mode: mode.value,
    }
  }).then((res) => {
    console.log(res, res.data.data.length)
    let tmp = []
    for (var i = 0; i < res.data.data.length; i++) {
      tmp.push({ pid: res.data.data[i].id, Title: res.data.data[i].title, Author: res.data.data[i].userName, src: res.data.data[i].url, pages: res.data.data[i].countPage, authorId: res.data.data[i].userId, r18: res.data.data[i].r18,seriesId:res.data.data[i].seriesId,SeriesTitle:res.data.data[i].SeriesTitle,Description: res.data.data[i].description })
    }
    novelItems.value = tmp
    waterfall.value.renderer()
  }).catch((error) => {
    console.log(error, error)
  }).finally(() => {
    console.log("ok")
    loading.value = false
    wait.value = false
  })

  // 根据类型和页数返回不同数据
  if (type === 'follow') {
    return page === 1 ? novelItems.value : []
  } else {
    return page === 1 ? novelItems.value.reverse() : []
  }
}

const applyFilters = () => {
  currentPage.value = 1
  // 重新渲染瀑布流
  nextTick(() => {
    if (waterfall.value) {
      waterfall.value.renderer()
    }
  })
}

const handleSearch = () => {
  // 防抖处理
  clearTimeout(handleSearch.timer)
  handleSearch.timer = setTimeout(() => {
    applyFilters()
  }, 300)
}

const loadMore = () => {
  currentPage.value++
  if (activeTab.value === 'follow') {
    fetchFollowNovels()
  } else {
    fetchRankingNovels()
  }
}

const refreshData = () => {
  currentPage.value = 1
  if (activeTab.value === 'follow') {
    fetchFollowNovels()
  } else {
    fetchRankingNovels()
  }
}

// 监听标签页变化
watch(() => activeTab.value, () => {
  currentPage.value = 1
  hasMore.value = true
})

// 监听视图模式变化
watch(() => viewMode.value, () => {
  nextTick(() => {
    if (viewMode.value === 'waterfall' && waterfall.value) {
      waterfall.value.renderer()
    }
  })
})

// 生命周期
onMounted(() => {
  // 默认加载关注小说
  fetchFollowNovels()
})
</script>

<style lang="less" scoped>
.novel-page-container {
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

// 页面标题
.page-header {
  text-align: center;
  margin-bottom: 30px;
  color: white;

  .page-title {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 15px;
    font-size: 32px;
    font-weight: 700;
    margin: 0 0 10px 0;
    text-shadow: 0 2px 4px rgba(0,0,0,0.3);

    .title-icon {
      font-size: 36px;
    }
  }

  .page-description {
    font-size: 16px;
    opacity: 0.9;
    margin: 0;
  }
}

// 导航卡片
.nav-card {
  margin-bottom: 20px;
  border-radius: 15px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.1);

  .novel-tabs {
    :deep(.el-tabs__nav) {
      width: 100%;
    }

    :deep(.el-tabs__item) {
      font-weight: 600;

      .tab-label {
        display: flex;
        align-items: center;
        gap: 8px;

        .tab-badge {
          :deep(.el-badge__content) {
            font-size: 10px;
          }
        }
      }
    }
  }
}

// 筛选卡片
.filter-card {
  margin-bottom: 25px;
  border-radius: 15px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.1);

  .filter-header {
    display: flex;
    align-items: center;
    gap: 10px;
    font-weight: 600;
    color: #409EFF;
    margin-bottom: 20px;
    font-size: 16px;
  }

  .filter-content {
    .filter-item {
      margin-bottom: 15px;

      .filter-label {
        display: block;
        margin-bottom: 8px;
        font-weight: 600;
        color: #606266;
        font-size: 14px;
      }

      .filter-select {
        width: 100%;
      }
    }

    .search-section {
      margin-top: 20px;

      .search-input {
        max-width: 400px;
      }
    }
  }
}

// 内容区域
.content-area {
  .stats-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 25px;
    padding: 20px;
    background: white;
    border-radius: 12px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.05);

    .el-tag {
      padding: 8px 16px;
      border-radius: 20px;
    }

    @media (max-width: 768px) {
      flex-direction: column;
      gap: 15px;
    }
  }

  // 瀑布流容器
  .waterfall-container {
    .novel-wrapper {
      margin-bottom: 20px;
      border-radius: 15px;
      overflow: hidden;
      transition: all 0.3s ease;

      &:hover {
        transform: translateY(-3px);
      }
    }
  }

  // 列表容器
  .list-container {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 20px;

    @media (max-width: 768px) {
      grid-template-columns: 1fr;
    }
  }

  // 加载更多
  .load-more-section {
    text-align: center;
    margin: 40px 0;

    .load-more-btn {
      padding: 12px 40px;
      border-radius: 25px;
      font-weight: 600;
      box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 6px 16px rgba(64, 158, 255, 0.4);
      }
    }
  }

  // 空状态
  .empty-state {
    margin: 60px 0;
  }
}

// 回到顶部
.back-to-top {
  width: 45px;
  height: 45px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);

  &:hover {
    transform: scale(1.1);
    box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
  }
}

// 动画效果
.el-fade-in-linear-enter-active {
  transition: all 0.4s ease;
}

.el-fade-in-linear-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

// 响应式设计
@media (max-width: 768px) {
  .novel-page-container {
    padding: 15px;
  }

  .page-header {
    .page-title {
      font-size: 28px;
      flex-direction: column;
      gap: 10px;
    }
  }

  .filter-content {
    .el-row {
      margin: 0 !important;

      .el-col {
        padding: 0 !important;
        margin-bottom: 15px;
      }
    }
  }
}

@media (max-width: 480px) {
  .page-header .page-title {
    font-size: 24px;
  }

  .waterfall-container {
    :deep(.vue-waterfall) {
      --waterfall-item-width: 100% !important;
    }
  }
}
</style>