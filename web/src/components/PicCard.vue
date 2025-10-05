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
                      :url="getImageUrl()"
                      class="image"

                  />
                  <div class="image-overlay" @click="openImageViewer">
                    <div class="overlay-content">
                      <el-icon class="preview-icon"><View /></el-icon>
                      <span class="preview-text">点击预览</span>
                    </div>
                  </div>
                </div>
                <div style="padding: 14px">
                    <el-row>
                        <el-text class="w-280px mb-2" truncated
                            @click="jump('https://www.pixiv.net/artworks/' + $props.pid)">
                            {{ $props.title }}
                        </el-text>
                    </el-row>
                    <el-row>
                        <el-text class="w-280px mb-2" truncated type="primary"
                            @click="'https://www.pixiv.net/users/' + $props.authorId">{{ $props.author }}</el-text>
                    </el-row>
                    <el-row>
                        <el-col :span="20" style="text-align: left;;">
                            <el-text class="w-250px mb-2" truncated type="success">
                                Pages:{{ $props.pages }}

                            </el-text>
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
            <!-- R18标识 -->
            <div v-if="$props.r18 > 0" class="r18-badge">
              <el-tag type="danger" size="small">R18</el-tag>
            </div>
        </template>
    </el-skeleton>
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";
import noProfileImg from '../assets/images/NoR18.png';
import { DownloadByPid } from "../../bindings/main/internal/pixivlib/ctl.js";
const name = "PicCard"
import { form } from "../assets/js/configuration.js"
import { LazyImg } from "vue-waterfall-plugin-next";
import { sleep } from "../assets/js/Time.js"
import {
  Download,
  View,
  User,
  Picture,
  PictureRounded,
  Link,
  Hide,
  Loading,
  Check
} from "@element-plus/icons-vue";
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
    }
});

const pic = ref(null)
const inqueue = ref(false)
const dis = ref(false)
// 获取图片URL
const getImageUrl = () => {
  if (props.r18 && !form.value.r_18) {
    return noProfileImg
  }
  return `http://127.0.0.1:7234/api/preview?url=${props.img}`
}

const openImageViewer = () => {
  const imageData = {
    pid: props.pid,
    pages: props.pages,
    title: props.title,
    author: props.author,
    authorId: props.authorId,
    r18: props.r18,
    img: props.img
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
function jump(event) {
    console.log("jump", event)
    window.open(event, '_blank')
}
</script>

<style lang="less" scoped>
.image {
    width: 100%;
}

.loading {
    width: 28px;
    height: 28px;
    border: 2px solid #ffffff;
    border-top-color: transparent;
    border-radius: 100%;
    animation: circle infinite 0.75s linear;
}
.image-container {
  position: relative;
  cursor: pointer;
  overflow: hidden;

  .image {
    width: 100%;
    transition: transform 0.3s ease;
  }

  .image-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: opacity 0.3s ease;

    .overlay-content {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 8px;
      color: white;
      transform: translateY(10px);
      transition: transform 0.3s ease;

      .preview-icon {
        font-size: 32px;
      }

      .preview-text {
        font-size: 14px;
        font-weight: 500;
      }
    }
  }

  &:hover {
    .image {
      transform: scale(1.05);
    }

    .image-overlay {
      opacity: 1;

      .overlay-content {
        transform: translateY(0);
      }
    }
  }
}
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
    //display: flex;
    //align-items: center;
    //justify-content: center;

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
@keyframes circle {
    0% {
        transform: rotate(0);
    }

    100% {
        transform: rotate(360deg);
    }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@keyframes checkMark {
  0% {
    transform: scale(0);
    opacity: 0;
  }
  50% {
    transform: scale(1.2);
    opacity: 1;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

.r18-badge {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 10;

  .el-tag {
    font-weight: bold;
    border-radius: 20px;
  }
}
</style>
