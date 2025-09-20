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
                <el-tag v-if="currentTask" type="primary" size="small">
                  {{ currentTask }}
                </el-tag>
                <el-tag type="info" size="small">
                  {{ percent }}%
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

    <!-- 回到顶部 -->
    <el-backtop
        :bottom="60"
        :right="60"
        :visibility-height="300"
    >
      <div class="back-to-top">
        <el-icon><ArrowUp /></el-icon>
      </div>
    </el-backtop>
  </div>
</template>

<script lang="js" setup>
import DateChoose from "./DateChoose.vue";
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
import { DownloadByPid, DownloadByRank, DownloadByNovelId, DownloadByAuthorId } from "../../bindings/main/internal/pixivlib/ctl.js";

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
      let tmp = msg.data[1]
      let tt = [];
      for(let k in msg.data[1]){
        let wk = tmp[k];

        if (wk.task === "No task"){
          percent.value = 0
        }else{
          percent.value = Math.round(wk.task.Current / Math.max(1,wk.task.Total) * 100)
          tt.push({value:wk.task.Name})
          // console.log(wk.task.Name)
        }
      }
      for (let v of msg.data[0]){
        // console.log(v.info)
        tt.push({value:v.info.Name})
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

// 删除指定任务
const removeTask = (index) => {
  if (index <= 0) {
    ElMessage.warning('无法删除正在执行的任务')
    return
  }

  queue.value.splice(index, 1)


}

const currentTask = computed(() => {
  return queue.value.length > 0 ? queue.value[0].value : null
});

function getLogType(log) {
  if (log.includes('错误') || log.includes('Error')) {
    return 'log-error'
  } else if (log.includes('警告') || log.includes('Warning')) {
    return 'log-warning'
  } else if (log.includes('完成') || log.includes('Success')) {
    return 'log-success'
  }
  return 'log-info'
}

function getCurrentTime() {
  return new Date().toLocaleTimeString()
}

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
        result = await DownloadByNovelId(value)
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

// function Download() {
//     // console.log("Downloading ", now.value)
//
//     return
// }
</script>

<style lang="less" scoped>
@import "../assets/style/font.less";
@import "../assets/style/variable.less";
@import "../assets/style/color.less";

.main-container {
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

// 页面标题
.page-header {
  text-align: center;
  margin-bottom: 15px;

  .page-title {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 15px;
    font-size: 32px;
    font-weight: 700;
    color: white;
    margin: 0;
    text-shadow: 0 2px 4px rgba(0,0,0,0.3);

    .title-icon {
      font-size: 36px;
    }
  }
}

// 下载表单卡片
.download-forms-card {
  margin-bottom: 25px;
  border-radius: 15px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.1);

  .card-header {
    display: flex;
    align-items: center;
    gap: 10px;
    font-weight: 600;
    font-size: 16px;
  }

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
      display: flex;
      align-items: flex-end;
      gap: 20px;
      flex-wrap: wrap;

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

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 6px 16px rgba(64, 158, 255, 0.4);
        }
      }
    }
  }
}

// 主内容区域 - 左主右侧边栏布局
.main-content {
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: 25px;
  margin-bottom: 30px;

  @media (max-width: 1200px) {
    grid-template-columns: 1fr 300px;
    gap: 20px;
  }

  @media (max-width: 900px) {
    grid-template-columns: 1fr;
    gap: 20px;
  }
}

// 左侧主要内容区域
.main-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.progress-card {
  border-radius: 15px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.1);

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    font-weight: 600;

    .progress-info {
      display: flex;
      gap: 8px;
    }
  }

  .main-progress {
    :deep(.el-progress__text) {
      font-size: 16px;
      font-weight: 600;
    }
  }

  .progress-text {
    font-size: 14px;
    font-weight: 600;
    color: #409EFF;
  }
}

.terminal-card {
  border-radius: 15px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.1);

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    font-weight: 600;

    .terminal-actions {
      display: flex;
      gap: 5px;
    }
  }

  .terminal-container {
    background: #1e1e1e;
    border-radius: 8px;
    overflow: hidden;

    .terminal-scrollbar {
      :deep(.el-scrollbar__wrap) {
        background: transparent;
      }
    }

    .terminal-content {
      padding: 15px 0;
      font-family: 'Consolas', 'Monaco', monospace;
      font-size: 14px;
      line-height: 1.5;

      .terminal-line {
        display: flex;
        margin-bottom: 8px;
        align-items: flex-start;

        .terminal-time {
          color: #666;
          margin-right: 10px;
          font-size: 12px;
          white-space: nowrap;
        }

        .terminal-text {
          flex: 1;
          color: #e0e0e0;
          word-break: break-all;
        }

        &.log-error .terminal-text {
          color: #ff4757;
        }

        &.log-warning .terminal-text {
          color: #ffa726;
        }

        &.log-success .terminal-text {
          color: #26de81;
        }

        &.log-info .terminal-text {
          color: #74b9ff;
        }
      }

      .terminal-empty {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 200px;
        color: #666;
        gap: 10px;
      }
    }
  }
}

// 右侧任务队列侧边栏
.sidebar-section {
  .queue-sidebar {
    border-radius: 15px;
    box-shadow: 0 8px 32px rgba(0,0,0,0.1);
    height: fit-content;
    position: sticky;
    top: 20px;

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
            //transform: translateX(3px);
          }
        }

        .queue-empty {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          height: 400px;
          color: #ffffff;

          .empty-icon {
            margin-bottom: 15px;

            .el-icon {
              font-size: 48px;
            }
          }

          .empty-text {
            text-align: center;

            p {
              margin: 0 0 5px 0;
              font-size: 16px;
              font-weight: 600;
            }

            span {
              font-size: 14px;
              opacity: 0.8;
            }
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

// 动画
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

// 响应式设计
@media (max-width: 768px) {
  .main-container {
    padding: 15px;
  }

  .page-header .page-title {
    font-size: 28px;
    flex-direction: column;
    gap: 10px;
  }

  .download-form .form-content {
    flex-direction: column;
    align-items: stretch;

    .download-btn {
      width: 100%;
      margin-top: 15px;
    }
  }

  .main-content {
    .sidebar-section {
      order: -1; // 在移动端将侧边栏移到上方
    }
  }
}
</style>
