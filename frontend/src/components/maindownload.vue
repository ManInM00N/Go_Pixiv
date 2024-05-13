<template>
  <el-container>
    <el-main
      style="
      display: flex;
        flex-direction: column;
        height: 100%;
      "
    >
      <el-row >
        <el-col :span="16"/>
        <el-col :span="8">
          下载模式
          <el-select
            v-model="now"
            ref="mode"
            class="m-2"
            size="default"
            @change="changetype"
            style="width:150px"
          >
            <el-option
                v-for="item in modes"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-col>
      </el-row >
      <el-divider class="Divide_Line"/>
      <el-row>
        <el-col :span="8" v-if="now!=='Rank'">
        </el-col>
        <el-col :span="8" v-if="now!=='Rank'">
          <el-input
              v-model="inputValue"
              size="large"
              placeholder="Pid/AuthorId"
              clearable
              type="number"
          ></el-input>
        </el-col>
        <el-col :span="8" v-if="now==='Rank'">
          <el-select
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
                :disabled="wait"
            />
          </el-select>
        </el-col>
        <el-col :span="8" v-if="now==='Rank'">
          <date-choose
            key="main"
            ref="dateSelect"
          ></date-choose>
        </el-col>
        <el-col :span="2"/>
        <el-col :span="6">
          <el-button
              style=""
              id="bt"
              type="success"
              size="large"
              @click="Download"
              :disabled="wait"
          >
            Download
            <el-icon size="large">
            <Download/>
            </el-icon>
          </el-button>
        </el-col>
      </el-row >
      <el-row style="height: 20px"/>
      <el-row >
        <el-col :span="1" />
        <el-col :span="15" class="Tre">
          <el-progress
              stroke-width="28"
              striped
              striped-flow
              :duration="10"
              :percentage="percent"
              text-inside="true"
              style=""
          >
            <span style="font-size:16px" v-if="queue.length>0">{{queue[0].value}} {{percent}}%</span>
          </el-progress>
        </el-col>
        <el-col :span="2"/>
        <el-col :span="6">
          <el-text class="Tre">
            {{queuenow}}
          </el-text>
        </el-col>
      </el-row>
      <br>
      <el-row
        class="Get_Remain"
      >
        <el-col :span="1" />
        <el-col
            :span="15"
            class="terminal-text"
        >
          <el-scrollbar
              class="text Micro"
              style="height:500px"
          >
            <p v-for="item in logs" >
              {{item}}
            </p>
          </el-scrollbar>
        </el-col>
        <el-col :span="1" />
        <el-col :span="7" >
          <el-table
              :data="queue"
              :cell-class-name="cellStyle"
              class="Half_light Tre queueTable"

              style="height:490px"
          >
              <el-table-column
                  label="TaskQueue"
                  prop="value"
                  style="height:60px"
              />
          </el-table>
        </el-col>
      </el-row>
    </el-main>
    <el-footer>
      <el-text
      type="danger">
        此软件为免费开源，如果是购买获得请退款举报
      </el-text>
    </el-footer>
  </el-container>
</template>

<script lang="ts" setup >
import DateChoose from "./DateChoose.vue";
import {DownloadByAuthorId, DownloadByPid,DownloadByRank} from "../../wailsjs/go/main/App.js";
import {EventsEmit, EventsOn} from "../../wailsjs/runtime";
import {defineComponent, onMounted, ref} from "vue";
import emitter from "../assets/js/Pub.js";
import {ElMessage} from "element-plus";
import axios from "axios";

const props = defineProps({
  Input: Number,
  wait: Boolean,
  ws:WebSocket,
})
onMounted(()=>{
    props.ws.value.onmessage = (event) => {
      // res.value = event.data;
      handleMessage(JSON.parse(event.data));
    };
})
const logs=ref([]);
const rows= ref(10); // 可根据需要调整展示的行数
const cellStyle = ({  rowIndex }) => {
  if (rowIndex === 0) {
    return 'Xbord'
  }
}
function handleMessage(data){
  if (data.type==1){
    percent.value=data.newnum
  }else if (data.type==2){

  }else if ()
}
const percent= ref(0);
const  now = ref("Pid")
const queue=ref([
])
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
const modes=ref([
  {
    value:"Pid",
    label:"By Pid",
  },
  {
    value:"Author",
    label:"By AuthorId",
  },
  {
    value:"Rank",
    label:"By Rank",
  },
]);
const inputValue= ref('');
const period=ref("daily");
EventsOn("UpdateProcess",function(newnum){
  console.log(newnum[0])
  percent.value=newnum[0];
});
EventsOn("Push",function(newmsg){
  console.log(newmsg[0])
  queue.value.push({value:newmsg[0]})
});
EventsOn("Pop",function(){
  queue.value.shift()
});
EventsOn("UpdateTerminal",function (newmsg) {
  logs.value.push(newmsg[0])
  if (logs.value.length>50){
    logs.value.pop()
  }
})

</script>
<style lang="less" scoped>
@import "../assets/style/font.less";
@import "../assets/style/variable.less";
@import "../assets/style/color.less";
.terminal-text {
  background-color: rgba(0,0,0,0.3);
  font-family: monospace;
  border: none;
  white-space: pre-wrap;
}
.text{
  background: rgba(@quartz,0.1);
}
.queueTable {
  -webkit-background-clip: text;
  //opacity: 0.5;
}
.No_Background{
  background: rgba(ff, ff, ff, 0.3);
  border: 2px solid #CD7F32;
}
/deep/.Xbord{
  width: 80%;
  position: relative;
  text-align: center;
  font-size: 24px;
  &::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    border: 2px solid rgb(17, 36, 100);
    transition: all .5s;
    animation: clippath 3s infinite linear;
  }
  @keyframes clippath {
    0%, 100% { clip-path: inset(0 0 95% 0); }
    25% { clip-path: inset(0 95% 0 0); }
    50% { clip-path: inset(95% 0 0 0); }
    75% { clip-path: inset(0 0 0 95%); }
  }
}
/deep/ .el-table {
  thead {
    color: #fff;
    font-weight: 500;
    background: linear-gradient(to right, rgba(#6fa3fe,0.5), rgba(#4cdafe,0.5)) !important;
    & th {
      background-color: transparent;
    }
    & tr {
      background-color: transparent;
    }
  }
}

</style>