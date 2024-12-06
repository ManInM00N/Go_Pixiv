<template>
    <el-main style="padding-left: 5px;padding-right: 5px;padding-top: 0px" class="wrap">
        <el-backtop :right="100" :bottom="100" />

        <el-row>
            <el-col>
                <el-text type="warning">
                    <h1 style="text-align:center;font-size:30px">
                        已关注用户的作品
                    </h1>
                </el-text>
            </el-col>
        </el-row>
        <el-row>
            <el-col :span="19" />
            <el-col :span="5">
                <el-button @click="downloadthispage">
                    download this page
                    <el-icon size="large">
                        <Download />
                    </el-icon>
                </el-button>
            </el-col>
        </el-row>
        <el-row class="head">
            <el-col :span="20" />
            <el-col :span="4">

            </el-col>
        </el-row>
        <el-row>
            <el-col :span="5" />
            <el-col :span="14" class="center">
                <el-pagination background layout="prev, pager, next" :total="1000" :page-count="34"
                    @current-change="handlePageChange">
                </el-pagination>
            </el-col>
            <el-col :span="5">
                <el-select class="m-2" size="large" v-model="mode">
                    <el-option value="all" label="all" @click="Download" />
                    <el-option value="r18" label="r18" @click="Download" />
                </el-select>
            </el-col>
        </el-row>

        <Waterfall :list="picitem" width="300" background-color="" animation-effect="fadeInUp" key="followWaterfall">
            <template #default="{ item, url, index }">
                <transition name="el-fade-in-linear">
                    <div class="card" v-if="true">
                        <PicCard :author="item.Author" :img="item.src" :title="item.Title" :pid="item.pid"
                            :authorId="item.authorId" :pages="item.pages" :r18="item.r18" :key="item.pid + 'follow'" />
                    </div>
                </transition>
            </template>
        </Waterfall>
        <el-footer v-if="loading == true">
            <div class="loader" id="loader">
                <br>
                <br>
                <br>
                <br>
                <br>
                <br>
                <br>
                <div class="loading">
                    <span></span>
                    <span></span>
                    <span></span>
                    <span></span>
                    <span></span>
                </div>
            </div>
        </el-footer>
    </el-main>
</template>

<script lang="ts" setup>
import { ref, onMounted, defineComponent } from "vue";
import { Events } from "@wailsio/runtime";
import { PreloadFollow, PopFollowPool, DownloadByFollowPage } from "../../bindings/main/ctl.js";
import { LazyImg, Waterfall } from 'vue-waterfall-plugin-next';
import 'vue-waterfall-plugin-next/dist/style.css'
import 'animate.css';
import PicCard from './PicCard.vue';
import "../assets/style/variable.less"
defineComponent({
    PicCard,
})
// const props=defineProps({
//   limit:{
//     type:Boolean,
//     default:true
//   },
//   form:{
//     type:DAO.Settings,
//   }
// })
const props = defineProps([
    'limit',
    'form',
])
const picitem = ref([])
const currentPage = ref(1)
const name = ref("follow")
const mode = ref("all")
const loading = ref(true)
function Download() {
    console.log(currentPage.value)
    PopFollowPool()

}
Events.On("PopUp", function () {
    console.log("114514")
    loading.value = true
    picitem.value = []
    PreloadFollow(currentPage.value.toString(), mode.value)
})
Events.On("UpdateLoadFollow", function (msg) {
    console.log(msg.data[0], msg, msg.data[0][1])
    // picitem.value.push({ pid: msg.data[0][0], Title: msg.data[0][1], Author: msg.data[0][2], src: "cache/images/" + msg.data[0][0] + ".jpg", pages: msg.data[0][3], authorId: msg.data[0][4], r18: msg.data[0][5] })
    picitem.value.push({ pid: msg.data[0][0], Title: msg.data[0][1], Author: msg.data[0][2], src: msg.data[0][6], pages: msg.data[0][3], authorId: msg.data[0][4], r18: msg.data[0][5] })
})
Events.On("FollowLoadOk", function () {
    loading.value = false
})
function handlePageChange(Page) {
    console.log("Page changed")
    currentPage.value = Page
    Download()
}
const debug = () => {
    console.log(picitem.value)
}
onMounted(function () {
    loading.value = true
    // Download()
    window.debug = debug
})
function downloadthispage() {
    DownloadByFollowPage(currentPage.value.toString(), mode.value)
}
</script>

<style lang="less" scoped>
@import "../assets/style/load.less";

.center {
    display: flex;
    justify-content: center;
}

.wrap {
    height: 100%;
    overflow-x: hidden;
}
</style>
