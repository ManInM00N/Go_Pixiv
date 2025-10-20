<template>
  <div class="settings-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">
        <el-icon class="title-icon"><Setting /></el-icon>
        应用设置
      </h1>
    </div>

    <!-- 设置表单 -->
    <div class="settings-content">
      <el-form
          :model="form"
          label-position="top"
          class="settings-form"
          ref="formRef"
          :rules="formRules"
      >
        <!-- 基础配置 -->
        <el-card class="settings-section" shadow="hover">
          <template #header>
            <div class="section-header">
              <el-icon><Connection /></el-icon>
              <span>连接配置</span>
            </div>
          </template>

          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item label="本地代理服务地址" prop="prefix">
                <el-input
                    v-model="form.prefix"
                    disabled
                    size="large"
                >
                  <template #prepend>
                    <el-icon><Link /></el-icon>
                  </template>
                </el-input>
                <div class="form-help">代理服务端口</div>
              </el-form-item>
            </el-col>

            <el-col :xs="24" :sm="12">
              <el-form-item label="代理端口" prop="proxy">
                <el-input
                    v-model="form.proxy"
                    :disabled="!form.useproxy"
                    size="large"
                    placeholder="如: 7890"
                >
                  <template #prepend>
                    <el-icon><DataLine /></el-icon>
                  </template>
                </el-input>
                <div class="form-help">HTTP代理端口号</div>
              </el-form-item>
            </el-col>

            <el-col :span="24">
              <el-form-item label="Cookie 认证" prop="cookie">
                <el-input
                    v-model="form.cookie"
                    type="textarea"
                    :rows="3"
                    placeholder="请输入从 Pixiv 网站获取的 Cookie"
                    show-word-limit
                    resize="none"
                >
                </el-input>
                <div class="form-help">
                  <el-icon><InfoFilled /></el-icon>
                  从浏览器开发者工具中复制 Pixiv 的完整 Cookie
                </div>
              </el-form-item>
            </el-col>
          </el-row>
        </el-card>

        <!-- 下载设置 -->
        <el-card class="settings-section" shadow="hover">
          <template #header>
            <div class="section-header">
              <el-icon><Download /></el-icon>
              <span>下载设置</span>
            </div>
          </template>

          <el-row :gutter="20">
            <el-col :xs="24" :sm="12">
              <el-form-item label="下载目录" prop="downloadposition">
                <el-input
                    v-model="form.downloadposition"
                    size="large"
                    placeholder="下载文件保存路径"
                >
                  <template #prepend>
                    <el-icon><Folder /></el-icon>
                  </template>
                </el-input>
                <div class="form-help">文件将保存到此目录</div>
              </el-form-item>
            </el-col>

            <el-col :xs="24" :sm="12">
              <el-form-item label="作品收藏数限制" prop="likelimit">
                <el-input-number
                    v-model="form.likelimit"
                    :min="0"
                    :max="999999"
                    size="large"
                    controls-position="right"
                    style="width: 100%"
                />
                <div class="form-help">只下载收藏数超过此值的作品 (0 = 无限制)</div>
              </el-form-item>
            </el-col>
          </el-row>
        </el-card>

        <!-- 性能调优 -->
        <el-card class="settings-section" shadow="hover">
          <template #header>
            <div class="section-header">
              <el-icon><Stopwatch /></el-icon>
              <span>性能调优</span>
            </div>
          </template>

          <el-row :gutter="20">
            <el-col :xs="24" :sm="8">
              <el-form-item label="下载间隔 (毫秒)" prop="downloadinterval">
                <el-input-number
                    v-model="form.downloadinterval"
                    :min="100"
                    :max="10000"
                    :step="100"
                    size="large"
                    controls-position="right"
                    style="width: 100%"
                />
                <div class="form-help">两次下载之间的等待时间</div>
              </el-form-item>
            </el-col>

            <el-col :xs="24" :sm="8">
              <el-form-item label="请求限制等待 (毫秒)" prop="retry429">
                <el-input-number
                    v-model="form.retry429"
                    :min="1000"
                    :max="60000"
                    :step="1000"
                    size="large"
                    controls-position="right"
                    style="width: 100%"
                />
                <div class="form-help">遇到 429 错误时的等待时间</div>
              </el-form-item>
            </el-col>

            <el-col :xs="24" :sm="8">
              <el-form-item label="失败重试间隔 (毫秒)" prop="retryinterval">
                <el-input-number
                    v-model="form.retryinterval"
                    :min="500"
                    :max="10000"
                    :step="500"
                    size="large"
                    controls-position="right"
                    style="width: 100%"
                />
                <div class="form-help">下载失败后的重试等待时间</div>
              </el-form-item>
            </el-col>

            <el-col :span="24">
              <el-form-item label="缓存过期时间 (天)" prop="expired_time">
                <el-input-number
                    v-model="form.expired_time"
                    :min="1"
                    :max="365"
                    size="large"
                    controls-position="right"
                    style="width: 200px"
                />
                <div class="form-help">本地缓存数据的保存天数</div>
              </el-form-item>
            </el-col>
          </el-row>
        </el-card>

        <!-- 功能开关 -->
        <el-card class="settings-section" shadow="hover">
          <template #header>
            <div class="section-header">
              <el-icon><Operation /></el-icon>
              <span>功能开关</span>
            </div>
          </template>

          <div class="switch-grid">
            <div class="switch-item">
              <div class="switch-content">
                <div class="switch-info">
                  <div class="switch-title">
                    <el-icon><Hide /></el-icon>
                    R-18 内容
                  </div>
                  <div class="switch-description">
                    允许下载和显示 R-18 分级的成人内容
                  </div>
                </div>
                <el-switch
                    v-model="form.r_18"
                    size="large"
                    active-color="#f56c6c"
                    inactive-color="#dcdfe6"
                />
              </div>
            </div>

            <div class="switch-item">
              <div class="switch-content">
                <div class="switch-info">
                  <div class="switch-title">
                    <el-icon><User /></el-icon>
                    按作者分类
                  </div>
                  <div class="switch-description">
                    下载时按作者名称创建子文件夹
                  </div>
                </div>
                <el-switch
                    v-model="form.differauthor"
                    size="large"
                />
              </div>
            </div>

            <div class="switch-item">
              <div class="switch-content">
                <div class="switch-info">
                  <div class="switch-title">
                    <el-icon><Connection /></el-icon>
                    本地代理
                  </div>
                  <div class="switch-description">
                    启用本地 HTTP 代理服务器
                  </div>
                </div>
                <el-switch
                    v-model="form.useproxy"
                    size="large"
                />
              </div>
            </div>
          </div>
        </el-card>

        <!-- 操作按钮 -->
        <div class="form-actions">
          <el-button
              type="primary"
              size="large"
              @click="handleUpdate"
              :loading="updating"
              class="update-btn"
          >
            <el-icon><CircleCheck /></el-icon>
            保存设置
          </el-button>

          <el-button
              size="large"
              @click="handleReset"
              class="reset-btn"
          >
            <el-icon><RefreshLeft /></el-icon>
            重置设置
          </el-button>

          <el-button
              type="success"
              size="large"
              @click="testConnection"
              :loading="testing"
              class="test-btn"
          >
            <el-icon><Connection /></el-icon>
            测试连接
          </el-button>
<!--          <el-button-->
<!--            size="large"-->
<!--            @click="GetWebView2Cookies('114')"-->
<!--            class="test-btn"-->
<!--          >-->
<!--            功能测试-->
<!--          </el-button>-->
        </div>
      </el-form>
    </div>

    <!-- 状态信息 -->
    <div class="status-info">
      <el-alert
          :type="connectionStatus.type"
          :title="connectionStatus.title"
          :description="connectionStatus.description"
          show-icon
          :closable="false"
          class="status-alert"
      />
    </div>
  </div>
</template>
<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import {
  Setting,
  Connection,
  Link,
  DataLine,
  InfoFilled,
  Download,
  Folder,
  Stopwatch,
  Operation,
  Hide,
  User,
  CircleCheck,
  RefreshLeft,
  ArrowUp
} from '@element-plus/icons-vue'
import { ElMessage, ElNotification, ElMessageBox } from 'element-plus'
import { Events } from "@wailsio/runtime"
import { form, updateSettings, defaultConf } from "../assets/js/configuration.js"
import {CheckLogin, GetWebView2Cookies} from "../../bindings/main/internal/pixivlib/ctl.js"

// 响应式数据
const formRef = ref(null)
const updating = ref(false)
const testing = ref(false)

// 表单验证规则
const formRules = reactive({
  cookie: [
    { required: true, message: '请输入 Cookie', trigger: 'blur' }
  ],
  downloadposition: [
    { required: true, message: '请输入下载目录', trigger: 'blur' }
  ],
  downloadinterval: [
    { type: 'number', min: 100, message: '下载间隔不能小于 100ms', trigger: 'change' }
  ],
  retry429: [
    { type: 'number', min: 1000, message: '请求限制等待时间不能小于 1000ms', trigger: 'change' }
  ],
  retryinterval: [
    { type: 'number', min: 500, message: '重试间隔不能小于 500ms', trigger: 'change' }
  ]
})

// 连接状态
const connectionStatus = computed(() => {
  if (form.value.logined) {
    return {
      type: 'success',
      title: '连接正常',
      description: '已成功登录 Pixiv，可以正常使用所有功能'
    }
  } else if (form.value.cookie) {
    return {
      type: 'warning',
      title: '未验证连接',
      description: '已配置 Cookie，但尚未验证登录状态'
    }
  } else {
    return {
      type: 'error',
      title: '未配置认证',
      description: '请配置有效的 Cookie 以使用 Pixiv 相关功能'
    }
  }
})

// 更新设置
async function handleUpdate() {
  let needReconn = false
  try {
    // 表单验证
    // await formRef.value.validate()

    updating.value = true

    ElNotification({
      type: "info",
      title: "正在保存",
      message: "正在更新设置并验证登录状态...",
      position: 'bottom-right',
      duration: 2000,
    })

    needReconn = await updateSettings()
    ElNotification({
      type:"success",
      message:"设置保存成功",
      position:"bottom-right",
      duration:2000,
    })
    console.log(needReconn)
    if (needReconn){

      await testConnection()
    }
  } catch (error) {
    if (error.errorFields) {
      ElMessage.error('请检查表单输入')
    } else {
      console.error('保存设置失败:', error)
      ElMessage.error('保存设置失败，请重试')
    }
  } finally {
    updating.value = false

  }
}

// 重置设置
async function handleReset() {
  try {
    await ElMessageBox.confirm(
        '确定要重置所有设置吗？此操作不可撤销。',
        '确认重置',
        {
          confirmButtonText: '确定重置',
          cancelButtonText: '取消',
          type: 'warning',
        }
    )
    form.value = defaultConf

    ElMessage.success('设置已重置')

  } catch (error) {
    // 用户取消操作
  }
}

// 测试连接
async function testConnection() {

  try {
    testing.value = true
    ElMessage.info('正在测试连接...')
    await CheckLogin()

  } catch (error) {
    console.error('连接测试失败:', error)
    ElMessage.error('连接测试失败')
  } finally {
    testing.value = false
  }
}

onMounted(() => {
  // 组件挂载完成
})
</script>

<style lang="less" scoped>
// 导入通用样式
@import "../assets/style/common/page-header.less";
@import "../assets/style/common/cards.less";
@import "../assets/style/common/buttons.less";
@import "../assets/style/common/responsive.less";
@import "../assets/style/variable.less";
@import "../assets/style/color.less";

// 主容器
.settings-container {
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

// 设置内容
.settings-content {
  max-width: 900px;
  margin: 0 auto;

  .settings-form {
    // 表单项标签
    :deep(.el-form-item__label) {
      font-weight: 600;
      color: #ffffff;
      font-size: 14px;
    }

    // 帮助文本
    .form-help {
      font-size: 12px;
      color: #909399;
      margin-top: 5px;
      display: flex;
      align-items: center;
      gap: 4px;
    }

    // 开关网格
    .switch-grid {
      display: grid;
      gap: 20px;

      .switch-item {
        padding: 20px;
        background: #f8f9fa;
        border-radius: 12px;
        border: 2px solid transparent;
        transition: all 0.3s ease;

        &:hover {
          border-color: #409EFF;
          background: #ecf5ff;
        }

        .switch-content {
          display: flex;
          justify-content: space-between;
          align-items: center;
          gap: 15px;

          .switch-info {
            flex: 1;

            .switch-title {
              display: flex;
              align-items: center;
              gap: 8px;
              font-weight: 600;
              color: #303133;
              margin-bottom: 5px;
            }

            .switch-description {
              font-size: 13px;
              color: #606266;
              line-height: 1.4;
            }
          }
        }
      }
    }
  }
}

// 操作按钮特定样式
.form-actions {
  justify-content: center;
  gap: 15px;
  margin: 30px 0;

  .el-button {
    padding: 12px 30px;
    border-radius: 25px;
    font-weight: 600;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
    }
  }

  .update-btn {
    background: linear-gradient(135deg, #67c23a, #85ce61);
    border: none;

    &:hover {
      background: linear-gradient(135deg, #85ce61, #95d475);
    }
  }

  .test-btn {
    background: linear-gradient(135deg, #409EFF, #66b1ff);
    border: none;

    &:hover {
      background: linear-gradient(135deg, #66b1ff, #82c4ff);
    }
  }
}

// 状态信息
.status-info {
  max-width: 900px;
  margin: 0 auto 30px auto;

  .status-alert {
    border-radius: 12px;
    border: none;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);

    :deep(.el-alert__title) {
      font-weight: 600;
    }
  }
}

// 响应式特定调整
@media (max-width: 768px) {
  .settings-content {
    .settings-form {
      .switch-grid {
        .switch-item {
          .switch-content {
            flex-direction: column;
            align-items: flex-start;
            gap: 15px;
          }
        }
      }
    }
  }
}

</style>