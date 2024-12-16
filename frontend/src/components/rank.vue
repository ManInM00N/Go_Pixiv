<template>
    <el-main style="padding-left: 5px;padding-right: 5px" class="scrollbar">
        <el-row class="ret">
            <el-col :span="1" />
            <el-col :span="6">
                <el-select :disabled="lock" v-model="period" class="m-2" size="large" style="width:150px">
                    <el-option v-for="(item) in options" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </el-col>
            <el-col :span="6">
                <date-choose :lock="lock" key="rank" ref="dateSelect" :re="true" @Rankload="Rankload"></date-choose>
            </el-col>
            <el-col :span="1" />
            <el-col :span="4">
                <el-select :disabled="lock" class="m-2" size="large" v-model="pages">
                    <el-option v-for="(item, index) in 10" :key="item" :label="item" :value="String(item)" />
                </el-select>
            </el-col>
            <el-col :span="4">
                <el-button size="large" @click="RankPage" :disabled="lock">
                    <el-icon size="25px">
                        <Search />
                    </el-icon>
                </el-button>
            </el-col>
        </el-row>

        <el-row>
            <el-col :span="8">
            </el-col>
            <el-col :span="11">
            </el-col>
            <el-col :span="5">
                <h1>
                    <el-button @click="downloadthispage" size="large">
                        download this page
                        <el-icon size="large">
                            <Download />
                        </el-icon>
                    </el-button>
                </h1>

            </el-col>
        </el-row>

        <Waterfall ref="waterfall" :list="picitem" :width=300 background-color="" :animationEffect="fadeInUp"
            key="rankWaterfall">
            <template #default="{ item, url, index }">
                <transition name="el-fade-in-linear">
                    <div class="card">
                        <PicCard :author="item.Author" :img="item.src" :title="item.Title" :pid="item.pid"
                            :authorId="item.authorId" :pages="item.pages" :r18="item.r18" :key="item.pid + 'follow'" />
                    </div>
                </transition>
            </template>
        </Waterfall>
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
    </el-main>

</template>

<script lang="ts" setup>
import DateChoose from "./DateChoose.vue";
import PicCard from "./PicCard.vue";
import { defineComponent, onMounted, ref } from "vue";
import { DownloadByRank } from "../../bindings/main/ctl.js";
import emitter from "../assets/js/Pub.js";
import { Download } from "@element-plus/icons-vue";
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
            picitem.value.push({ pid: String(res.data.data[i].illust_id), Title: res.data.data[i].title, Author: res.data.data[i].user_name, src: res.data.data[i].url, pages: Number(res.data.data[i].illust_page_count), authorId: String(res.data.data[i].user_id), r18: res.data.data[i]['illust_content_type.sexual'] })
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

.wrap {
    height: 100vh;
    overflow: scroll;
}

.scrollbar {
    height: calc(100vh - 40px);
}

/deep/.el-scrollbar__wrap {
    overflow-x: hidden;
}
</style>
