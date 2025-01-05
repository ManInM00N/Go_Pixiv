<template>
    <el-main style="padding-left: 5px;padding-right: 5px;padding-top: 0px" class="wrap">


        <el-row>
            <el-col>
                <el-text type="warning">
                    <h1 style="text-align:center;font-size:30px;padding: 0;
                    line-height: 1.7em;letter-spacing: 1px;
                    font-family: Georgia, Times, serif;
                    color:#e8ffe8;">
                        已关注用户的作品
                    </h1>
                </el-text>
            </el-col>
        </el-row>
        <el-row style="margin-bottom: 5px">
            <el-col :span="19" />
            <el-col :span="5">
                <el-button @click="downloadthispage" :disabled="wait">
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
            <el-col :span="7"  style="display: flex">
              <el-col :span="8">

                <el-button @click="()=>{DownloadFollow(from,to)}" :disabled="wait">
                  下载多页
                </el-button>
              </el-col>
              <el-col :span="8" style="padding-left: 5px" >
                <el-input type="number" v-model="from" @change="fixchange"/>

              </el-col>
              <el-col :span="1">
                <span>-</span>
              </el-col>
              <el-col :span="8" >
                <el-input type="number" v-model="to" @change="fixchange"/>
              </el-col>
            </el-col>

            <el-col :span="14" class="center">
                <el-pagination background layout="prev, pager, next" :total="1000" :page-count="34"
                    @current-change="handlePageChange" :disabled="wait">
                </el-pagination>
            </el-col>
            <el-col :span="3">
<!--                <el-select class="m-2" size="large" v-model="mode">-->
<!--                    <el-option value="all" label="all" @click="Download" />-->
<!--                    <el-option value="r18" label="r18" @click="Download" />-->
<!--                </el-select>-->
            </el-col>
        </el-row>
        <Waterfall ref="waterfall" :list="picitem" width=300 background-color="" animation-effect="fadeInUp"
            key="followWaterfall">
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
      <el-backtop :right="100" :bottom="100" />

    </el-main>
</template>

<script lang="ts" setup>
import { ref, onMounted, defineComponent } from "vue";
import { DownloadByFollowPage } from "../../bindings/main/ctl.js";
import {DownloadFollow} from "../assets/js/download.js"
import { Waterfall } from 'vue-waterfall-plugin-next';
import axios from 'axios'
import 'vue-waterfall-plugin-next/dist/style.css'
import 'animate.css';
import PicCard from './PicCard.vue';
import "../assets/style/variable.less"
const waterfall = ref(null)
const from = ref(1),to=ref(10)
defineComponent({
    PicCard,
})
const picitem = ref([])
const currentPage = ref(1)
const name = ref("follow")
const mode = ref("all")
const wait = ref("false")
const loading = ref(true)
function fixchange(){
  from.value = Math.max(1,from.value)
  to.value = Math.max(1,to.value)
  from.value = Math.min(20,from.value)
  to.value = Math.min(20,to.value)
    if (from.value > to.value ){
      let tt = to.value
      to.value = from.value
      from .value = tt
    }
}
function FollowMsg() {
    wait.value = true
    loading.value = true
    console.log("doing get")
    picitem.value = []
    axios.get("http://127.0.0.1:7234/api/followpage", {
        params: {
            p: currentPage.value.toString(),
            mode: mode.value,
        }
    }).then((res) => {
        console.log(res, res.data.data.length)
        let tmp = []
        for (var i = 0; i < res.data.data.length; i++) {
            picitem.value.push({ pid: res.data.data[i].id, Title: res.data.data[i].title, Author: res.data.data[i].userName, src: res.data.data[i].url, pages: res.data.data[i].countPage, authorId: res.data.data[i].userId, r18: res.data.data[i].r18 })
        }
        waterfall.value.renderer()
    }).catch((error) => {
        console.log(error, error)
    }).finally(() => {
        console.log("ok")
        loading.value = false
        wait.value = false
    })
}
function handlePageChange(Page) {
    console.log("Page changed")
    currentPage.value = Page
    FollowMsg()
}
const debug = () => {
    console.log(picitem.value)
}
onMounted(function () {
    loading.value = true
    FollowMsg()
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

.image-slot {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
    background: var(--el-fill-color-light);
    color: var(--el-text-color-secondary);
    font-size: 30px;
}

.image-slot .el-icon {
    font-size: 30px;
}
</style>
