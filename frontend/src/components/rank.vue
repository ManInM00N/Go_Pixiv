<template>
  <el-main
    style="padding-left: 5px;padding-right: 5px"
  >
    <el-row>
      <el-col :span="1"/>
      <el-col :span="6">
        <el-select
            :disabled="lock"
            v-model="period"
            class="m-2"
            size="large"
            style="width:150px"
        >
          <el-option
              v-for="(item) in options"
              :key="item.value"
              :label="item.label"
              :value="item.value"
          />
        </el-select>
      </el-col>
      <el-col :span="6">
        <date-choose
            :lock="lock"
            key="rank"
            ref="dateSelect"
            :re="true"
            @Rankload="Rankload"
        ></date-choose>
      </el-col>
      <el-col :span="1"/>
      <el-col :span="4">
        <el-select
            :disabled="lock"
            class="m-2"
            size="large"
            v-model="pages"
        >
          <el-option
              value="1"
              label="Page 1"
          />
          <el-option
              value="2"
              label="Page 2"
          />
        </el-select>
      </el-col>
      <el-col :span="4">
        <el-button
            @click="Rankload"
        >
          <el-icon size="25px"><Search/> </el-icon>
        </el-button>

      </el-col>
    </el-row>
    <br>
    <el-skeleton
        :loading="loading"
        animated
        :throttle="500"
    >
      <template #template>
        <el-space
          wrap
        >
          <el-skeleton-item variant="image" style="width: 290px; height: 240px" />
          <el-skeleton-item variant="image" style="width: 290px; height: 240px" />
          <el-skeleton-item variant="image" style="width: 290px; height: 240px" />
          <el-skeleton-item variant="image" style="width: 290px; height: 240px" />
          <el-skeleton-item variant="image" style="width: 290px; height: 240px" />
          <el-skeleton-item variant="image" style="width: 290px; height: 240px" />

        </el-space>
        <br>
        <h1 size="60px">{{tip}}  {{loadup}}/{{sum}}</h1>
      </template>
      <template #default>
        <el-row>
          <el-col :span="8">
            <el-text
              type="primary"
              align="center"
            >
              {{pagemsg}}
            </el-text>
            <el-text>
              {{remainderTime}}
            </el-text>
          </el-col>
          <el-col :span="4">
            <el-text>
              download this page
            </el-text>
          </el-col>
          <el-col>
            <el-icon>
              <Download/>
              <el-button @click="downloadthispage"/>

            </el-icon>
          </el-col>
        </el-row>
<!--        <el-container-->
<!--            v-masonry-->
<!--        >-->
<!--            <PicCard-->
<!--                v-for="(item, index) in picitem"-->
<!--                :author="item.Author"-->
<!--                :img="item.src"-->
<!--                :title="item.Title"-->
<!--                :pid="item.pid"-->
<!--                :pages="item.pages"-->
<!--                :key="index"-->
<!--                v-masonry-tile-->
<!--            />-->
<!--        </el-container>-->
        <MasonryWall
            :gap="9"
            :column-width="260"
            :items="picitem"
        >
          <template
              #default="{item,index}"
          >
            <PicCard
                :author="item.Author"
                :img="item.src"
                :title="item.Title"
                :pid="item.pid"
                :authorId="item.authorId"
            />
          </template>
        </MasonryWall>
      </template>

    </el-skeleton>

  </el-main>
</template>

<script lang="ts" setup>
import DateChoose from "./DateChoose.vue";
import PicCard from "./PicCard.vue";
import {DownloadByAuthorId, DownloadByPid, DownloadByRank, PreloadRank,PopLoadPool} from "../../wailsjs/go/main/App.js";
import {EventsEmit, EventsOn} from "../../wailsjs/runtime";
import {defineComponent, onMounted, ref,toRef} from "vue";
import emitter from "../assets/js/Pub.js";
import MasonryWall from '@yeger/vue-masonry-wall'
import {Download} from "@element-plus/icons-vue";
defineComponent({
  PicCard, DateChoose,MasonryWall
})
name: "maindownload";
const picitem  =ref([])
const period=ref("daily");
const lock= ref(false)
const tip = ref("Search a page Please")
const options = ref([
  {
    value:"daily",
    label:"Daily",
  },
  {
    value:"weekly",
    label:"Weekly",
  },
  {
    value:"monthly",
    label:"Monthly",
  },
  {
    value:"daily_r18",
    label:"Daily_R18",
  },
  {
    value:"weekly_r18",
    label:"Weekly_R18",
  },
]);
const re_Date=ref(new Date());
const remainderTime=ref('')
const dateSelect= ref(null)
const loading = ref(true)
const sum = ref(100)
const loadup = ref(0)
const pagemsg = ref('')
const pages = ref(1)
const downloadthispage = ()=>{
  console.log(dateSelect.value.selectedDate);
  emitter.emit("DownloadByRank",{date:dateSelect.value.selectedDate,period:period.value})
}
function Rankload(){
  pagemsg.value=dateSelect.value.selectedDate+period.value
  tip.value="Please Wait..."
  PopLoadPool()
  re_Date.value=new Date()
  picitem.value=[]
  loadup.value=0;
  sum.value=100
  console.log("preload ",dateSelect.value.selectedDate,period.value)
  lock.value=true
  loading.value=true;
  PreloadRank(dateSelect.value.selectedDate,period.value)
}
EventsOn("UpdateLoad",function(msg){
  console.log(msg[0])
  picitem.value = [...picitem.value,{pid:msg[0],Title:msg[1],Author: msg[2],src: "cache/images/"+msg[0]+".jpg",pages:msg[3],authorId:msg[4]}]
  loadup.value++;
})
EventsOn("LoadOk",function(){
  loading.value=false;
  lock.value=false;
  ComputeDate()
})
onMounted(function(){
  loading.value=true;
  // PreloadRank(dateSelect.value.selectedDate,period.value)
})

function  ComputeDate(){
  var strDate=new Date(re_Date.value);
  var endDate = new Date(); // 结束时间
  var diffDate = endDate.getTime() - strDate.getTime()
  var days = Math.floor(diffDate / (24 * 3600 * 1000));
  var leave1 = diffDate % (24 * 3600 * 1000);
  var hours = Math.floor(leave1 / (3600 * 1000));
  var leave2 = leave1 % (3600 * 1000);
  var minutes = Math.floor(leave2 / (60 * 1000));
  var leave3 = leave2 % (60 * 1000);
  var seconds = Math.round(leave3 / 1000);
  remainderTime.value=minutes+'分钟' + seconds +"秒"
  console.log(remainderTime.value);
}
</script>
<style lang="less" scoped>

</style>