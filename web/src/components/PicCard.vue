<template>
    <!--    -->
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
                <LazyImg :ref="pic" :onload="load = false"
                    :url="!($props.r18 && !form.r_18) ? 'http://127.0.0.1:7234/api/preview?url=' + $props.img : noProfileImg"
                    class="image" />
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
                            <br>
                            <el-text v-if="$props.r18" class="w-250px mb-2" truncated type="danger">
                                R18

                            </el-text>
                        </el-col>

                        <el-col :span="4">
                            <div class="bottom card-header">
                                <el-button text class="button" v-if="!inqueue" @click="download" :disabled="inqueue">
                                    <el-icon size="30">
                                        <Download />
                                    </el-icon>
                                </el-button>
                                <div v-if="inqueue" style="text-align: center;height: 32px;padding-left: 15px;">
                                    <div class="loading" v-if="inqueue && dis"></div>
                                    <el-icon size="30" v-if="inqueue && !dis">
                                        <Select color="green" />
                                    </el-icon>
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
import noProfileImg from '../assets/images/NoR18.png';
import { DownloadByPid } from "../../bindings/main/internal/pixiv/ctl.js";
const name = "PicCard"
import { form } from "../assets/js/configuration.js"
import { LazyImg } from "vue-waterfall-plugin-next";
import { sleep } from "../assets/js/Time.js"
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
        type: Boolean,
        default: true,
    }
});
const pic = ref(null)
const inqueue = ref(false)
const dis = ref(false)
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
    border: 2px solid #000;
    border-top-color: transparent;
    border-radius: 100%;
    animation: circle infinite 0.75s linear;
}

@keyframes circle {
    0% {
        transform: rotate(0);
    }

    100% {
        transform: rotate(360deg);
    }
}
</style>
