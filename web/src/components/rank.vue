<template>
  <el-main style="padding: 20px" class="scrollbar" ref="mainContainer">
    <!-- 筛选控制面板 -->
    <div class="page-header">
      <h1 class="page-title">
        <el-icon class="title-icon"><Histogram /></el-icon>
        排行榜
        <el-tag v-if="picitem.length > 0" type="info" class="count-tag">
          {{ picitem.length }} 作品
        </el-tag>
      </h1>
    </div>
      <el-card class="control-panel" shadow="hover">

      <el-row :gutter="20" class="control-row">
        <!-- 排名周期选择 -->
        <el-col :xs="24" :sm="12" :md="4" :lg="6">
          <div class="control-item">
            <label class="control-label">排行类型</label>
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
        <el-col :xs="24" :sm="12" :md="5" :lg="6">
          <div class="control-item" >
            <label class="control-label">日期</label>
            <date-choose
                :lock="lock"
                key="rank"
                ref="dateSelect"
                :re="true"
                @Rankload="Rankload"
            ></date-choose>
          </div>
        </el-col>


        <el-col :xs="24" :sm="12" :md="4" :lg="6">
          <div class="control-item">
            <label class="control-label">页号</label>
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
        <el-col :xs="24" :sm="12" :md="10" :lg="6">
          <div class="control-item">
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

    <div class="waterfall-container" v-show="!loading">
      <Waterfall
          ref="waterfall"
          :list="picitem"
          :width=waterFallConf.width
          :gutter=waterFallConf.gutter
          :breakpoints=waterFallConf.breakpoints
          background-color="transparent"
          :animationEffect="fadeInUp"
          key="rankWaterfall"
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

    <el-footer v-if="loading === true">
      <div class="loader" id="loader">
        <br v-for="_ in 6">
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
import {Download, Search, Picture, DataLine, ArrowUp, UserFilled} from "@element-plus/icons-vue";
import { LazyImg, Waterfall } from 'vue-waterfall-plugin-next'
import { Events } from "@wailsio/runtime";
import 'vue-waterfall-plugin-next/dist/style.css'
import 'animate.css';
import axios from 'axios'
import {waterFallConf} from "../assets/js/configuration.js";
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
// 导入通用样式
@import "../assets/style/common/page-header.less";
@import "../assets/style/common/cards.less";
@import "../assets/style/common/pagination.less";
@import "../assets/style/common/waterfall.less";
@import "../assets/style/common/buttons.less";
//@import "../assets/style/common/loading.less";
@import "../assets/style/common/animations.less";
@import "../assets/style/common/responsive.less";
@import "../assets/style/load.less";

// 主容器
.scrollbar {
  padding: 20px;
  min-height: 100vh;
  overflow-y: hidden;
  overflow-x: hidden;
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

@media (max-width: 480px) {
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
</style>