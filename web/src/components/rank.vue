<template>
  <el-main style="padding: 20px" class="scrollbar" ref="mainContainer">
    <!-- 筛选控制面板 -->
    <el-card class="filter-card" shadow="hover">
      <div class="filter-header">
        <h2 class="page-title">
          <el-icon><Picture /></el-icon>
          Pixiv Ranking
        </h2>
      </div>

      <el-row :gutter="20" class="filter-row">
        <!-- 排名周期选择 -->
        <el-col  :span="4">
          <div class="filter-item">
            <label class="filter-label">Ranking Period</label>
            <el-select
                :disabled="lock"
                v-model="period"
                class="filter-select"
                size="large"
                style="width: 100%"
                placeholder="Select period"
            >
              <el-option
                  v-for="item in options"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
              />
            </el-select>
          </div>
        </el-col>
        <!-- 日期选择 -->
        <el-col :span="7">
          <div class="filter-item" >
            <label class="filter-label">Date</label>
            <date-choose
                :lock="lock"
                key="rank"
                ref="dateSelect"
                :re="true"
                @Rankload="Rankload"
            ></date-choose>
          </div>
        </el-col>


        <el-col  :span="4" >
          <div class="filter-item">
            <label class="filter-label">Page</label>
            <el-select
                :disabled="lock"
                class="filter-select"
                size="large"
                v-model="pages"
                style="width: 100%"
            >
              <el-option
                  v-for="item in 10"
                  :key="item"
                  :label="`Page ${item}`"
                  :value="String(item)"
              />
            </el-select>
          </div>
        </el-col>

        <!-- 搜索和下载按钮 -->
        <el-col  :span="3" style="margin-bottom: 8px" >
          <div class="action-buttons">
            <el-button
                type="primary"
                size="large"
                @click="RankPage"
                :disabled="lock"
                :loading="loading"
                class="search-btn"
            >
              <el-icon><Search /></el-icon>
              Search
            </el-button>
          </div>
        </el-col>
        <el-col  :span="3" style="margin-left: 24px;margin-bottom: 8px">
          <div class="action-buttons">
            <el-button
                type="success"
                size="large"
                @click="downloadthispage"
                :disabled="lock || picitem.length === 0"
                class="download-btn"
            >
              <el-icon><Download /></el-icon>
              Download
            </el-button>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <div class="waterfall-container">
      <Waterfall
          ref="waterfall"
          :list="picitem"
          :width="300"
          :gutter="20"
          background-color="transparent"
          :animationEffect="fadeInUp"
          key="rankWaterfall"
      >
        <template #default="{ item, url, index }">
          <transition name="el-fade-in-linear">
            <div class="card-wrapper">
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

    <el-footer v-if="loading === true">
      <div class="loader" id="loader">
        <div class="loading">
          <span></span>
          <span></span>
          <span></span>
          <span></span>
          <span></span>
        </div>
      </div>

    </el-footer>


    <!-- 空状态 -->
    <el-empty
        v-if="!loading && picitem.length === 0"
        description="No illustrations found"
        class="empty-state"
    >
      <el-button type="primary" @click="RankPage">Try Again</el-button>
    </el-empty>
  </el-main>

</template>

<script lang="ts" setup>
import DateChoose from "./DateChoose.vue";
import PicCard from "./PicCard.vue";
import { defineComponent, onMounted, ref } from "vue";
import { DownloadByRank } from "../../bindings/main/internal/pixivlib/ctl.js";
import { Download, Search, Picture, DataLine, ArrowUp } from "@element-plus/icons-vue";
import { LazyImg, Waterfall } from 'vue-waterfall-plugin-next'
import { Events } from "@wailsio/runtime";
import 'vue-waterfall-plugin-next/dist/style.css'
import 'animate.css';
import axios from 'axios'
defineComponent({
  PicCard, DateChoose
})
name: "rank";
const picitem = ref([])
const period = ref("daily");
const lock = ref(false)
const waterfall = ref(null)
const options = ref([
  {
    value: "daily",
    label: "Daily",
  },
  {
    value: "weekly",
    label: "Weekly",
  },
  {
    value: "monthly",
    label: "Monthly",
  },
  {
    value: "daily_r18",
    label: "Daily_R18",
  },
  {
    value: "weekly_r18",
    label: "Weekly_R18",
  },
]);
const dateSelect = ref(null)
const loading = ref(false)
const sum = ref(100)
const loadup = ref(0)
const pagemsg = ref('')
const pages = ref('1')
const downloadthispage = () => {
  console.log(dateSelect.value.selectedDate);
  DownloadByRank(dateSelect.value.selectedDate, period.value)
}
const nextpage = ref(0)
function RankPage() {
  loading.value = true
  lock.value = true
  console.log("doing get")
  console.log(pages, period)
  picitem.value = []
  waterfall.value.renderer()
  axios.get("http://127.0.0.1:7234/api/rankpage", {
    params: {
      p: pages.value,
      mode: period.value,
      content: "illust",
      date: dateSelect.value.selectedDate,
    }
  }).then((res) => {
    console.log(res, res.data.data.length)
    let tmp = []
    for (var i = 0; i < res.data.data.length; i++) {
      // tmp.push({ pid: String(res.data.data[i].illust_id), Title: res.data.data[i].title, Author: res.data.data[i].user_name, src: res.data.data[i].url, pages: Number(res.data.data[i].illust_page_count), authorId: String(res.data.data[i].user_id), r18: res.data.data[i]['illust_content_type.sexual'] })
      picitem.value.push({ pid: String(res.data.data[i].illust_id), Title: res.data.data[i].title, Author: res.data.data[i].user_name, src: res.data.data[i].url, pages: Number(res.data.data[i].illust_page_count), authorId: String(res.data.data[i].user_id), r18: res.data.data[i].illust_content_type.sexual})
    }
    // picitem.value.concat(tmp)
    waterfall.value.renderer()

  }).catch((error) => {
    console.log(error, error)
  }).finally(() => {
    console.log("ok")
    lock.value = false
    loading.value = false
  })

}

onMounted(function () {
  loading.value = true;
  RankPage()
})
</script>
<style lang="less" scoped>
@import "../assets/style/load.less";

.scrollbar {
  //height: calc(100vh - 40px);
  //overflow-y: auto;
  //overflow-x: hidden;
  padding: 20px;
  min-height: 100vh;
  overflow-y: hidden;
  //background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  overflow-x: hidden;
}
.main-container {
  padding: 20px;
  min-height: 100vh;
  overflow-y: hidden;
  //background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  overflow-x: hidden;
}

// 筛选卡片样式
.filter-card {
  margin-bottom: 20px;
  border-radius: 12px;

  .filter-header {
    margin-bottom: 20px;

    .page-title {
      display: flex;
      align-items: center;
      gap: 10px;
      margin: 0;
      color: #409EFF;
      font-size: 24px;
      font-weight: 600;
    }
  }

  .filter-row {
    align-items: flex-end;
  }

  .filter-item {
    margin-bottom: 10px;

    .filter-label {
      display: block;
      margin-bottom: 8px;
      font-weight: 500;
      color: #c2c7d1;
      font-size: 14px;
    }
  }

  .action-buttons {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
    @media (max-width: 768px) {
      flex-direction: column;

      .el-button {
        width: 100%;
        margin:auto;

      }
    }
  }
}

// 结果信息
.results-info {
  margin-bottom: 20px;
  text-align: center;

  .el-tag {
    padding: 8px 16px;
    border-radius: 20px;
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

// 加载动画
.loading-footer {
  padding: 40px 0;
  text-align: center;

  .loading-text {
    margin-top: 20px;
    color: #909399;
    font-size: 14px;
  }
}


// 空状态
.empty-state {
  margin: 60px 0;
}

// 响应式适配
@media (max-width: 768px) {
  .el-main {
    padding: 15px;
  }

  .filter-card {
    .action-buttons {
      margin-top: 15px;
    }
  }

  .waterfall-container {
    :deep(.vue-waterfall) {
      --waterfall-item-width: 280px !important;
    }
  }
}

@media (max-width: 480px) {
  .waterfall-container {
    :deep(.vue-waterfall) {
      --waterfall-item-width: 100% !important;
    }
  }

  .page-title {
    font-size: 20px !important;
  }
}

// 覆盖Element Plus默认样式
:deep(.el-scrollbar__wrap) {
  overflow-x: hidden;
}

:deep(.el-card__body) {
  padding: 24px;
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