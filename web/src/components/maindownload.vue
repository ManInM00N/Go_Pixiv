<template>
  <div class="main-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">
        <el-icon class="title-icon"><Download /></el-icon>
        任务池
      </h1>
    </div>

    <!-- 下载任务填写区域 -->
    <el-card class="download-forms-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <el-icon><EditPen /></el-icon>
          <span>创建下载任务</span>
        </div>
      </template>

      <el-tabs
          v-model="activeTab"
          type="card"
          class="download-tabs"
      >
        <!-- PID下载 -->
        <el-tab-pane label="作品ID" name="pid">
          <div class="download-form">
            <div class="form-content">
              <div class="input-group">
                <label class="input-label">作品ID (PID)</label>
                <el-input
                    v-model="inputPid"
                    size="large"
                    placeholder="请输入 Pixiv 作品ID"
                    clearable
                    type="number"
                    class="download-input"
                >
                  <template #prepend>
                    <el-icon><Picture /></el-icon>
                  </template>
                </el-input>
              </div>
              <el-button
                  type="primary"
                  size="large"
                  @click="() => handleDownload('pid', inputPid)"
                  :disabled="wait || !inputPid"
                  :loading="downloading"
                  class="download-btn"
              >
                <el-icon><Download /></el-icon>
                开始下载
              </el-button>
            </div>
          </div>
        </el-tab-pane>

        <!-- 作者ID下载 -->
        <el-tab-pane label="作者ID" name="author">
          <div class="download-form">
            <div class="form-content">
              <div class="input-group">
                <label class="input-label">作者ID (UID)</label>
                <el-input
                    v-model="inputAuthorId"
                    size="large"
                    placeholder="请输入 Pixiv 作者ID"
                    clearable
                    type="number"
                    class="download-input"
                >
                  <template #prepend>
                    <el-icon><User /></el-icon>
                  </template>
                </el-input>
              </div>
              <el-button
                  type="primary"
                  size="large"
                  @click="() => handleDownload('author', inputAuthorId)"
                  :disabled="wait || !inputAuthorId"
                  :loading="downloading"
                  class="download-btn"
              >
                <el-icon><Download /></el-icon>
                下载全部作品
              </el-button>
            </div>
          </div>
        </el-tab-pane>

        <!-- 排行榜下载 -->
        <el-tab-pane label="排行榜" name="rank">
          <div class="download-form">
            <div class="form-content rank-form">
              <div class="input-group">
                <label class="input-label">排行榜类型</label>
                <el-select
                    v-model="period"
                    size="large"
                    placeholder="选择排行榜类型"
                    class="download-input"
                >
                  <el-option
                      v-for="item in options"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                      :disabled="item.disabled && !form.r_18"
                  />
                </el-select>
              </div>
              <div class="input-group">
                <label class="input-label">日期选择</label>
                <date-choose
                    key="main"
                    ref="dateSelect"
                    class="date-picker"
                />
              </div>
              <el-button
                  type="primary"
                  size="large"
                  @click="() => handleDownload('rank')"
                  :disabled="wait || !period"
                  :loading="downloading"
                  class="download-btn"
              >
                <el-icon><Download /></el-icon>
                下载排行榜
              </el-button>
            </div>
          </div>
        </el-tab-pane>

        <!-- 小说下载 -->
        <el-tab-pane label="小说" name="novel">
          <div class="download-form">
            <div class="form-content">
              <div class="input-group">
                <label class="input-label">小说ID</label>
                <el-input
                    v-model="inputNovelId"
                    size="large"
                    placeholder="请输入小说ID"
                    clearable
                    type="number"
                    class="download-input"
                >
                  <template #prepend>
                    <el-icon><Document /></el-icon>
                  </template>
                  <template #append>
                    <el-checkbox v-model="isSeries">
                      是否为系列
                    </el-checkbox>
                  </template>
                </el-input>
              </div>
              <el-button
                  type="primary"
                  size="large"
                  @click="() => handleDownload('novel', inputNovelId)"
                  :disabled="wait || !inputNovelId"
                  :loading="downloading"
                  class="download-btn"
              >
                <el-icon><Download /></el-icon>
                下载小说
              </el-button>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 主内容区域 -->
    <div class="main-content">
      <!-- 左侧主要内容区域 -->
      <div class="main-section">
        <!-- 下载进度 -->
        <el-card class="progress-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><DataLine /></el-icon>
              <span>下载进度</span>
              <div class="progress-info">
                <el-tag v-if="currentTask != null" type="primary" size="small">
                  {{ currentTask.value }}
                </el-tag>
                <el-tag type="info" size="small">
                  {{ percent }}%
                </el-tag>
                <el-tag v-if="currentTask != null" type="info" size="small">
                  {{currentTask.data.task.Current}}/{{currentTask.data.task.Total}}
                </el-tag>
              </div>
            </div>
          </template>

          <el-progress
              :percentage="percent"
              :stroke-width="20"
              :show-text="false"
              :striped="percent > 0 && percent < 100"
              :striped-flow="percent > 0 && percent < 100"
              class="main-progress"

          >
<!--            <template #default="{ percentage }">-->
<!--                            <span class="progress-text">-->
<!--                                {{ currentTask || '等待任务...' }} {{ percentage }}%-->
<!--                            </span>-->
<!--            </template>-->
          </el-progress>
        </el-card>

        <!-- 终端日志 -->
        <el-card class="terminal-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Monitor /></el-icon>
              <span>下载日志</span>
              <div class="terminal-actions">
                <el-button
                    size="small"
                    circle
                    @click="clearLogs"
                    v-tooltip="'清空日志'"
                >
                  <el-icon><Delete /></el-icon>
                </el-button>
                <el-button
                    size="small"
                    circle
                    @click="scrollToBottom"
                    v-tooltip="'滚动到底部'"
                >
                  <el-icon><Bottom /></el-icon>
                </el-button>
              </div>
            </div>
          </template>

          <div class="terminal-container">
            <el-scrollbar
                ref="scrollbar"
                class="terminal-scrollbar"
                height="450px"
            >
              <div class="terminal-content">
                <div
                    v-for="(v, index) in logs"
                    :key="index"
                    class="terminal-line"
                    :class="getLogType(v.log)"
                >
                                    <span class="terminal-time">
                                        {{ v.time }}
                                    </span>
                  <span class="terminal-text">{{ v.log }}</span>
                </div>
                <div
                    v-if="logs.length === 0"
                    class="terminal-empty"
                >
                  <el-icon><Monitor /></el-icon>
                  <span>等待下载日志...</span>
                </div>
              </div>
            </el-scrollbar>
          </div>
        </el-card>
      </div>

      <div class="sidebar-section">
        <el-card class="queue-sidebar" shadow="hover">
          <template #header>
            <div class="sidebar-header">
              <div class="header-title">
                <el-icon><List /></el-icon>
                <span>任务队列</span>
              </div>
              <el-badge
                  :value="queue.length"
                  :max="99"
                  class="queue-badge"
                  :type="queue.length > 0 ? 'primary' : 'info'"
              />
            </div>
          </template>

          <div class="sidebar-content">
            <el-scrollbar height="580px">
              <div class="queue-list">
                <div
                    v-for="(task, index) in queue"
                    :key="index"
                    class="queue-item"
                    :class="{ 'active-task': index === 0 }"
                >
                  <div class="task-header">
                    <div class="task-status">
                      <el-icon
                          v-if="index === 0"
                          class="status-icon processing"
                      >
                        <Loading />
                      </el-icon>
                      <el-icon
                          v-else
                          class="status-icon waiting"
                      >
                        <Clock />
                      </el-icon>
                    </div>
                    <div class="task-index">
                                    <span v-if="index === 0" class="current-badge">
                                        正在处理
                                    </span>
                    </div>
                    <div v-if="index > 0" class="task-actions">
                      <el-button
                          size="small"
                          type="danger"
                          circle
                          @click="removeTask(index)"
                          v-tooltip="'取消任务'"
                          class="remove-task-btn"
                      >
                        <el-icon><Close /></el-icon>
                      </el-button>
                    </div>
                  </div>

                  <div class="task-content">
                    <div class="task-name" :title="task.value">
                      {{ task.value }}
                    </div>
                  </div>

                  <div v-if="index === 0" class="task-progress">
                    <el-progress
                        :percentage="percent"
                        :show-text="false"
                        :stroke-width="4"
                        status="info"
                    />
                    <div class="progress-text">{{ percent }}%</div>
                  </div>
                </div>

                <!-- 空状态 -->
                <div
                    v-if="queue.length === 0"
                    class="queue-empty"
                >
                  <div class="empty-icon">
                    <el-icon><Box /></el-icon>
                  </div>
                  <div class="empty-text">
                    <p>暂无下载任务</p>
                    <span>创建任务后将在此显示</span>
                  </div>
                </div>
              </div>
            </el-scrollbar>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 页脚信息 -->
    <div class="page-footer">
      <el-alert
          type="info"
          :closable="false"
          show-icon

      >
        <template #title>
          <span style=" color:#ff5555">此软件为免费开源，如果是购买获得请退款举报 {{ timeElement }}</span>
        </template>
      </el-alert>
    </div>

  </div>
</template>

<script lang="js" setup>
import DateChoose from "./DateChoose.vue";
import { getLogType,debounce } from '../assets/js/utils/index.js'
import { onMounted, ref, computed, nextTick } from "vue";
import { ElNotification, ElMessage } from "element-plus";
import {
  Download,
  EditPen,
  Picture,
  User,
  Document,
  DataLine,
  Monitor,
  List,
  Delete,
  Bottom,
  Loading,
  Clock,
  Box,
  ArrowUp
} from "@element-plus/icons-vue";
import axios from "axios";
import { timeElement } from "../assets/js/Time.js"
import { Events } from "@wailsio/runtime"
import { ws, form } from "../assets/js/configuration.js";
import {
  DownloadByPid,
  DownloadByRank,
  DownloadByNovelId,
  DownloadByAuthorId,
  RemoveTask
} from "../../bindings/main/internal/pixivlib/ctl.js";

// 响应式数据
const activeTab = ref('pid')
const dateSelect = ref(null)
const scrollbar = ref(null)
const logs = ref([])
const percent = ref(0)
const queue = ref([])
const downloading = ref(false)
const wait = ref(false)

// 输入数据
const inputPid = ref('')
const inputAuthorId = ref('')
const inputNovelId = ref('')
const period = ref("daily")
const isSeries =ref(false)

onMounted(() => {
    Events.On("NotFound", (msg) => {
      ElNotification({
        title: "Error",
        type: "error",
        message: msg.data[0][0],
        position: 'bottom-right',
        duration: 3000,
      })
    })
    let cnt = 0;
    Events.On("taskPoolInfos",function (msg){
      if (cnt%6===0){
        // console.log(msg.data,msg.data[1])
      }
      cnt = (cnt+1)%6
      // msg.data[1]
      const arr = msg.data[0] || []
      let tmp = msg.data[1]
      let tt = [];
      for(let k in msg.data[1]){
        let wk = tmp[k];

        if (wk.task === "No task"){
          percent.value = 0
        }else{
          percent.value = Math.round(wk.task.Current / Math.max(1,wk.task.Total) * 100)
          tt.push({value:wk.task.Name,data:wk})
          // console.log(wk)
        }
      }
      for (let v of arr){
        // console.log(v)
        tt.push({value:v.info.Name,data:{task:v.info,status:v.status}})
      }
      queue.value = tt
    })
    Events.On("UpdateTerminal", function (newmsg) {
      // console.log(newmsg)
      logs.value.push({log:newmsg.data[0][0],time:new Date().toLocaleTimeString()})
      if (logs.value.length > 200) {
        logs.value.pop()
      }
    })
    ws.value.onmessage = (event) => {
        handleMessage(JSON.parse(event.data));
    };
})
const rows = ref(10); // 可根据需要调整展示的行数
function handleMessage(data) {
  switch (data.type) {
    case 1: // 进度更新
      percent.value = data.newnum
      break
    case 2: // 任务完成
      queue.value.shift()
      break
    case 3: // 新任务
      queue.value.push(data.newtask)
      break
  }
}
const mode = ref('')
function changetype(data) {
    console.log(data)
    now.value = data
}
function changetype2(data) {
    console.log(data)
    period.value = data
}
const now = ref("Pid")
// 选项配置
const options = ref([
  { value: "daily", label: "每日排行", disabled: false },
  { value: "weekly", label: "每周排行", disabled: false },
  { value: "monthly", label: "每月排行", disabled: false },
  { value: "daily_r18", label: "每日排行 R18", disabled: true },
  { value: "weekly_r18", label: "每周排行 R18", disabled: true },
])

const debouncedRemoveTask = debounce(
    (index)=>{
      if (index <= 0) {
        ElMessage.warning('无法删除正在执行的任务')
        return
      }
      // console.log(queue.value[index].data.info.ID)
      RemoveTask(queue.value[index].data.info.ID)
      ElNotification({
        type:"info",
        position:"bottom-right",
        message:"删除成功",
        duration: 1000,
      })
    }
    ,300)

// 删除指定任务
const removeTask = (index) => {
  debouncedRemoveTask(index)
}

const currentTask = computed(() => {
  return queue.value.length > 0 ? queue.value[0] : null
});

function clearLogs() {
  logs.value = []
  ElMessage.success('日志已清空')
}

async function handleDownload(type, value = null) {
  try {
    downloading.value = true
    wait.value = true

    let result = false
    let taskName = ''
    
    switch (type) {
      case 'pid':
        result = await DownloadByPid(value)
        taskName = `作品 ${value}`
        break
      case 'author':
        result = await DownloadByAuthorId(value)
        taskName = `作者 ${value} 的全部作品`
        break
      case 'rank':
        result = await DownloadByRank(dateSelect.value.selectedDate, period.value)
        taskName = `${period.value} 排行榜 (${dateSelect.value.selectedDate})`
        break
      case 'novel':
        result = await DownloadByNovelId(value,isSeries.value)
        taskName = `小说 ${value}`
        break
    }

    if (result) {
      // 清空对应输入框
      ElMessage.success(`已添加下载任务: ${taskName}`)
      switch (type) {
        case 'pid': inputPid.value = ''; break
        case 'author': inputAuthorId.value = ''; break
        case 'novel': inputNovelId.value = ''; break
      }
    } else {
      ElMessage.error('下载任务创建失败')
    }
  } catch (error) {
    console.error('下载失败:', error)
    ElMessage.error('下载失败，请检查网络连接')
  } finally {
    downloading.value = false
    wait.value = false
  }
}

</script>

<style lang="less" scoped>
// 导入通用样式
@import "../assets/style/common/page-header.less";
@import "../assets/style/common/cards.less";
@import "../assets/style/common/buttons.less";
@import "../assets/style/common/loading.less";
@import "../assets/style/common/animations.less";
@import "../assets/style/common/responsive.less";
@import "../assets/style/font.less";
@import "../assets/style/variable.less";
@import "../assets/style/color.less";

// 主容器
.main-container {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

// 下载表单卡片特定样式
.download-forms-card {
  margin-bottom: 25px;
  border-radius: 15px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);

  .download-tabs {
    :deep(.el-tabs__nav) {
      background: transparent;
    }

    :deep(.el-tabs__item) {
      border-radius: 8px 8px 0 0;
      margin-right: 5px;
    }
  }

  .download-form {
    padding: 20px 10px;

    .form-content {
      &.rank-form {
        flex-direction: row;
        align-items: flex-end;
      }

      .input-group {
        flex: 1;
        min-width: 200px;

        .input-label {
          display: block;
          margin-bottom: 8px;
          font-weight: 600;
          color: #606266;
          font-size: 14px;
        }

        .download-input {
          width: 100%;
        }
      }

      .download-btn {
        height: 44px;
        padding: 0 30px;
        border-radius: 22px;
        font-weight: 600;
        box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
      }
    }
  }
}

// 主要区域
.main-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

// 队列侧边栏特定样式
.sidebar-section {
  .queue-sidebar {
    .sidebar-header {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .header-title {
        display: flex;
        align-items: center;
        gap: 10px;
        font-weight: 600;
        font-size: 16px;
      }

      .queue-badge {
        :deep(.el-badge__content) {
          font-weight: 600;
        }
      }
    }

    .sidebar-content {
      .queue-list {
        .queue-item {
          border-radius: 12px;
          padding: 15px;
          margin-bottom: 12px;
          background: rgba(64, 64, 64, 0.93);
          border: 2px solid transparent;
          transition: all 0.3s ease;

          &.active-task {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: rgba(255, 255, 255, 1);
            border-color: #409EFF;
            box-shadow: 0 4px 15px rgba(64, 158, 255, 0.3);

            .current-badge {
              background: rgba(255, 255, 255, 0.29);
              color: white;
              padding: 2px 8px;
              border-radius: 10px;
              font-size: 12px;
              font-weight: 600;
            }

            .status-icon.processing {
              animation: spin 1s linear infinite;
              color: #FFD700;
            }

            .progress-text {
              color: rgba(255, 255, 255, 0.9);
              font-size: 12px;
              text-align: center;
              margin-top: 5px;
            }
          }

          .task-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 10px;

            .task-status {
              .status-icon {
                font-size: 18px;

                &.waiting {
                  color: #909399;
                }
              }
            }

            .task-index {
              .queue-number {
                background: rgb(228, 231, 237);
                color: #121212;
                padding: 2px 8px;
                border-radius: 10px;
                font-size: 12px;
                font-weight: 600;
              }
            }
          }

          .task-content {
            .task-name {
              font-weight: 600;
              font-size: 14px;
              line-height: 1.4;
              word-break: break-all;
              margin-bottom: 8px;
            }
          }

          .task-progress {
            margin-top: 10px;
          }

          &:hover:not(.active-task) {
            background: rgba(159, 160, 161, 0.18);
            border-color: #b3d8ff;
          }
        }
      }
    }
  }
}

// 页脚
.page-footer {
  .el-alert {
    border-radius: 12px;
    :deep(.el-alert__title) {
      font-size: 14px;
    }
  }
}

</style>