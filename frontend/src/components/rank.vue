<template>


  <el-main
    style="padding-left: 5px;padding-right: 5px"
    class="scrollbar"

  >
    <el-row
      class="ret"
    >
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
            size="large"
            @click="Rankload"
            :disabled="lock"
        >
          <el-icon size="25px"><Search/> </el-icon>
        </el-button>
      </el-col>
    </el-row>

    <el-row>
      <el-col :span="8">

          <el-text
              type="danger"
          >
            <h1>
            {{pagemsg}}
              {{tip}}
            </h1>

          </el-text>

      </el-col>

      <el-col :span="11">
        <el-text
        type="success"
        >
          <h1>
            {{remainderTime}}

          </h1>
        </el-text>
      </el-col>
      <el-col :span="5">
        <h1>
          <el-button
              @click="downloadthispage"
              size="large"
          >
            download this page
            <el-icon
                size="large"
            >
              <Download/>
            </el-icon>
          </el-button>
        </h1>

      </el-col>
    </el-row>

    <Waterfall
        :list="picitem"
        width="300"
        background-color=""
        animation-effect="fadeInUp"
        key="rankWaterfall"
    >
      <template #item="{ item, url, index }">
        <div class="card">
          <PicCard
              :author="item.Author"
              :img="item.src"
              :title="item.Title"
              :pid="item.pid"
              :authorId="item.authorId"
              :pages="item.pages"
              :r18="item.r18"
              :limit="$props.form['r-18']"
          />
        </div>
      </template>
    </Waterfall>
    <el-footer v-if="loading===true">
      <div class="loader" id ="loader" >
        <div class ="loading">
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
import {DownloadByAuthorId, DownloadByPid, DownloadByRank, PreloadRank,PopLoadPool} from "../../wailsjs/go/main/App.js";
import {EventsEmit, EventsOn} from "../../wailsjs/runtime";
import {defineComponent, onMounted, ref} from "vue";
import emitter from "../assets/js/Pub.js";
import {Download} from "@element-plus/icons-vue";
import { LazyImg, Waterfall } from 'vue-waterfall-plugin-next'
import 'vue-waterfall-plugin-next/dist/style.css'
import 'animate.css';
import {DAO} from "../../wailsjs/go/models.ts";
defineComponent({
  PicCard, DateChoose
})
name: "rank";
// const props = defineProps({
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
const picitem  =ref([])
const period=ref("daily");
const lock= ref(false)
const tip = ref("Select a page First")
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
const loading = ref(false)
const sum = ref(100)
const loadup = ref(0)
const pagemsg = ref('')
const pages = ref('1')
const downloadthispage = ()=>{
  console.log(dateSelect.value.selectedDate);
  emitter.emit("DownloadByRank",{date:dateSelect.value.selectedDate,period:period.value})
}
function Rankload(){
  pagemsg.value=dateSelect.value.selectedDate+period.value
  tip.value=""
  tip.value="Please Wait..."
  remainderTime.value=""
  PopLoadPool()

}
EventsOn("UpdateLoad",function(msg){
  console.log(msg[0])
  // picitem.value.push({pid:msg[0],Title:msg[1],Author: msg[2],src: "cache/images/"+msg[0]+".jpg",pages:msg[3],authorId:msg[4]})
  picitem.value=picitem.value.concat({pid:msg[0],Title:msg[1],Author: msg[2],src: "cache/images/"+msg[0]+".jpg",pages:msg[3],authorId:msg[4],r18:msg[5]})
  loadup.value++;
})
EventsOn("LoadOk",function(){
  loading.value=false;
  lock.value=false;
  ComputeDate()
})
EventsOn("RankmsgPopUp",function(){
  re_Date.value=new Date()
  picitem.value=[]
  loadup.value=0;
  sum.value=100
  console.log("preload ",dateSelect.value.selectedDate,period.value,pages.value)
  lock.value=true
  loading.value=true;
  PreloadRank(dateSelect.value.selectedDate,period.value,pages.value)
})
onMounted(function(){
  loading.value=false;
  // PreloadRank(dateSelect.value.selectedDate,period.value)
})

function  ComputeDate(){
  var strDate=new Date(re_Date.value);
  var endDate = new Date();
  var diffDate = endDate.getTime() - strDate.getTime()
  var leave1 = diffDate % (24 * 3600 * 1000);
  var leave2 = leave1 % (3600 * 1000);
  var minutes = Math.floor(leave2 / (60 * 1000));
  var leave3 = leave2 % (60 * 1000);
  var seconds = Math.round(leave3 / 1000);
  remainderTime.value="耗时："+minutes+'分钟' + seconds +"秒"
  console.log(remainderTime.value);
}
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