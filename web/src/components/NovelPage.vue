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
            </div>
          </template>
        </el-tab-pane>

        <!-- 排行榜 -->
        <el-tab-pane v-if="false" label="排行榜" name="ranking">
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
                <el-option label="每日排行R18" value="daily_r18" />
                <el-option label="每周排行" value="weekly" />
                <el-option label="每周排行R18" value="weekly_r18" />
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
                <el-option v-for="tag in all_tags" :label="tag" :value="tag">

                </el-option>
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
      <div class="pagination-container">
        <el-pagination
            background
            layout="prev, pager, next, jumper"
            :total="1000"
            :page-count="activeTab === 'ranking'? 2 : 7"
            :current-page="currentPage"
            @current-change="handlePageChange"
            :disabled="loading"
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
            <el-icon v-if="!loading"><Refresh /></el-icon>
          </el-button>
          <el-tag type="info" size="large">
            <el-icon><DataLine /></el-icon>
            共找到 {{ filteredNovels.length }} 部小说
          </el-tag>
        </div>
      </div>
      <!-- 瀑布流展示 -->
      <div  class="waterfall-container">
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
import { createDebouncedSearch } from '../assets/js/utils/debounce.js'
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
const all_tags = ref(null)
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
    novels = novels.filter(novel => novel.tags.includes(selectedGenre.value))
  }

  // 字数筛选
  if (selectedWordCount.value) {
    novels = novels.filter(novel => {
      const count = novel.textCount
      switch (selectedWordCount.value) {
        case 'short': return count < 10000
        case 'medium': return count >= 10000 && count < 50000
        case 'long': return count >= 50000 && count < 200000
        case 'epic': return count >= 200000
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

const handlePageChange = (page) =>{
  currentPage.value = page
  if(activeTab.value === "follow"){
    fetchFollowNovels()
  }else{
    fetchRankingNovels()
  }
}

// 方法
const handleTabChange = (tabName) => {
  if (tabName === 'follow') {
    fetchFollowNovels()
  } else if (tabName === 'ranking') {
    currentPage.value = Math.min(currentPage.value,2)
    fetchRankingNovels()
  }
}

const fetchFollowNovels = async () => {
  try {
    loading.value = true
    // 调用获取关注作者小说的函数
    let novels = await fetchNovel('follow', currentPage.value)
    followNovels.value = novels
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

    rankingNovels.value = novels

    // hasMore.value = novels.length > 0

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

const fetchNovel = async (type, page = 1, options = {}) => {

  try {
    let url = ''
    let params = {
      p: page.toString(),
      type: 'novel'
    }

    if (type === 'follow') {
      url = 'http://127.0.0.1:7234/api/followpage'
      params.mode = 'all' // 或根据需要设置
    } else if (type === 'ranking') {
      url = 'http://127.0.0.1:7234/api/rankpage'
      params.mode = options.type || 'daily' // 使用排行榜类型
    }

    const response = await axios.get(url, { params })
    let tmp_tags = []
    console.log(response)
    // 映射 API 返回的数据到组件需要的格式
    const novels = response.data.data.map(item => ({
      // 基础信息
      id: item.id,
      title: item.title,
      content: item.content || '',
      cover: item.url || item.cover || '', // 封面图片

      // 作者信息
      userId: item.userId,
      userName: item.userName,
      profileImageUrl: item.profileImageUrl,

      // 详细信息
      description: item.description || '',
      page: item.countPage || 1,
      bookmarkCount: item.bookmarkCount || 0,
      textCount: item.textCount || 0,
      wordCount: item.wordCount || 0,

      // 状态和类型
      isLoginOnly: item.isLoginOnly || false,
      genre: item.genre || '',
      aiType: item.aiType || false,

      // 系列信息
      seriesNavData: item.seriesId ? {
        seriesType: 'novel',
        seriesId: item.seriesId,
        title: item.SeriesTitle || item.seriesTitle || '',
        order: item.seriesOrder || 1,
        prev: item.prevId ? {
          title: item.prevTitle || '',
          order: item.prevOrder || 0,
          id: item.prevId
        } : null,
        next: item.nextId ? {
          title: item.nextTitle || '',
          order: item.nextOrder || 0,
          id: item.nextId
        } : null
      } : null,

      // 标签信息
      tags:  item.tags || []
    }))
    for(let it in novels){
      tmp_tags = tmp_tags.concat(novels[it].tags)
    }
    all_tags.value = Array.from(new Set(tmp_tags))
    return novels

  } catch (error) {
    console.error('获取小说数据失败:', error)
    throw error
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

const handleSearch = createDebouncedSearch(() => {
  applyFilters()
}, 300)

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
// 导入通用样式
@import "../assets/style/common/page-header.less";
@import "../assets/style/common/cards.less";
@import "../assets/style/common/pagination.less";
@import "../assets/style/common/waterfall.less";
@import "../assets/style/common/buttons.less";
@import "../assets/style/common/loading.less";
@import "../assets/style/common/animations.less";
@import "../assets/style/common/responsive.less";

// 主容器
.novel-page-container {
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

// 导航卡片
.nav-card {
  margin-bottom: 20px;
  border-radius: 15px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);

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

// 筛选卡片特定样式
.filter-card {
  .filter-content {
    .filter-item {
      margin-bottom: 15px;

      .filter-label {
        display: block;
        margin-bottom: 8px;
        font-weight: 600;
        color: #ffffff;
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

// 内容区域特定配置
.content-area {
  .stats-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 25px;
    padding: 20px;
    background: white;
    border-radius: 12px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);

    .el-tag {
      padding: 8px 16px;
      border-radius: 20px;
    }

    @media (max-width: 768px) {
      flex-direction: column;
      gap: 15px;
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
}

// 响应式
@media (max-width: 768px) {
  .novel-page-container {
    padding: 15px;
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
  .waterfall-container {
    :deep(.vue-waterfall) {
      --waterfall-item-width: 100% !important;
    }
  }
}
</style>