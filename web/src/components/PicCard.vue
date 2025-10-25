<template>
    <el-skeleton style="width: 240px" :loading="load" animated :throttle="500">
        <template #template>
            <el-skeleton-item variant="image" style="width: 240px; height: 240px" />
            <div style="padding: 14px">
                <el-skeleton-item variant="h3" style="width: 50%" />
                <div style="
              display: flex;
              align-items: center;
              justify-items: stretch;
              margin-top: 16px;
              height: 16px;
            ">
                    <el-skeleton-item variant="text" style="margin-right: 16px" />
                    <el-skeleton-item variant="text" style="width: 30%" />
                </div>
            </div>
        </template>
        <template #default>
            <el-card :body-style="{ padding: '0px', marginBottom: '1px', width: '100%' }">
                <div class="image-container" >
                  <LazyImg
                      :ref="pic"
                      :onload="load = false"
                      :url="getImageUrl(props.img, props.r18)"
                      class="image"

                  />
                  <div class="image-overlay" @click="openImageViewer">
                    <div class="overlay-content">
                      <el-icon class="preview-icon"><View /></el-icon>
                      <span class="preview-text">点击预览</span>
                    </div>
                  </div>
                  <!-- R18标识 -->
                  <div class="cover-badges">
                    <el-tag v-if="$props.r18 > 1"  type="danger" size="small">R18</el-tag>
                    <el-tag v-if="$props.pages > 1" type="normal" size="small">{{$props.pages}}页</el-tag>
                    <el-tag v-if="$props.illustType == 2" type="primary" size="small">GIF</el-tag>
                  </div>
                </div>
                <div style="padding: 14px">
                    <el-row>
                        <el-text class="w-280px mb-2" truncated
                             @click="openPixivArtwork(props.pid)">
                            {{ $props.title }}
                        </el-text>
                    </el-row>
                    <el-row>
                        <el-text class="w-280px mb-2" truncated type="primary"
                            @click="openPixivUser( $props.authorId)">{{ $props.author }}</el-text>
                    </el-row>
                    <el-row>
                        <el-col :span="20" style="text-align: left;;">
                        </el-col>

                        <el-col :span="4">
                            <div class="bottom card-header download-area">
                                <el-button text class="button download-btn" v-if="!inqueue" @click="download" :disabled="inqueue">
                                    <el-icon size="30">
                                        <Download />
                                    </el-icon>
                                </el-button>
                                <div v-if="inqueue" class="download-status" style="text-align: center;height: 32px;">
                                    <div class="loading" v-if="inqueue && dis"></div>
                                    <div
                                        v-else
                                        class="download-success"
                                    >
                                      <el-icon class="success-icon"><Check /></el-icon>
                                    </div>
                                </div>

                            </div>
                        </el-col>
                    </el-row>

                </div>
            </el-card>

        </template>
    </el-skeleton>
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";
import {DownloadByPid} from "../../bindings/main/internal/pixivlib/ctl.js";
const name = "PicCard"
import { form } from "../assets/js/configuration.js"
import { LazyImg } from "vue-waterfall-plugin-next";
import {
  Download,
  View,
  Check
} from "@element-plus/icons-vue";
import {getImageUrl, openPixivArtwork, openPixivUser} from '../assets/js/utils/index.js'
import { sleep } from '../assets/js/utils/debounce.js'
import { useImageViewerStore } from "../assets/stores/PicPreview.js"
const store = useImageViewerStore()


const load = ref(true)
const props = defineProps({
    pid: {
        type: String,
    },
    author: {
        type: String,
    },
    title: {
        type: String,
        default: "确定",
    },
    img: {
        type: String,
        default: "",
    },
    pages: {
        type: Number,
        default: 1,
    },
    authorId: {
        type: String,
        default: 1,
    },
    r18: {
        type: Number,
        default: true,
    },
    illustType:{
      type: Number,
      default : 0,
    }
});

const pic = ref(null)
const inqueue = ref(false)
const dis = ref(false)

const openImageViewer = () => {
  const imageData = {
    pid: props.pid,
    pages: props.pages,
    title: props.title,
    author: props.author,
    authorId: props.authorId,
    r18: props.r18,
    img: props.img,
    illustType:props.illustType,
  }
  store.openViewer(imageData)
}

async function download() {
  if (DownloadByPid(props.pid)) {
    console.log(props.pid)
    dis.value = true
    inqueue.value = true
    await sleep(1000)
    dis.value = false
  }
}

onMounted(() => {
})
</script>


<style lang="less" scoped>
// 导入通用样式
@import "../assets/style/common/loading.less";
@import "../assets/style/common/waterfall.less";
@import "../assets/style/common/animations.less";

// 卡片特定的下载区域样式
.download-area {
  .download-btn {
    width: 32px;
    height: 32px;
    transition: all 0.3s ease;

    &:hover {
      transform: scale(1.1);
      box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
    }
  }

  .download-status {
    width: 32px;
    height: 32px;

    .loading-spinner {
      .spinning {
        animation: spin 1s linear infinite;
        color: #409EFF;
      }
    }

    .download-success {
      .success-icon {
        color: #67c23a;
        font-size: 32px;
        animation: checkMark 0.5s ease;
      }
    }
  }
}
</style>