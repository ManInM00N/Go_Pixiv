<template>
  <el-container>
    <el-main>
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
        <el-col :span="8" v-if="this.now!=='Rank'">
        </el-col>
        <el-col :span="8" v-if="this.now!=='Rank'">
          <el-input
              v-model="inputValue"
              size="large"
              placeholder="Pid/AuthorId"
              clearable
              type="number"
              @input="updateParentValue"
          ></el-input>
        </el-col>
        <el-col :span="8" v-if="this.now==='Rank'">
          <el-select
              v-model="period"
              ref="period"
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
        <el-col :span="8" v-if="this.now==='Rank'">
          <date-choose
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
        <el-col :span="6" >
          <el-text class="Tre">
            {{tasknow}}
          </el-text>
        </el-col>
        <el-col :span="2" />
        <el-col :span="8" class="Tre">
          <el-progress
              stroke-width="24"
              striped
              striped-flow
              :duration="10"
              :percentage="percent"
          >

          </el-progress>
        </el-col>
        <el-col :span="2"/>
        <el-col :span="6">
          <el-text class="Tre">
            {{queuenow}}
          </el-text>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="1" />
        <el-col :span="14">

        </el-col>
        <el-col :span="2" />
        <el-col :span="7" >
          <el-table
              :data="queue"
              :cell-style="cellStyle"
              class="Half_light Tre"

          >
              <el-table-column label="TaskQueue" prop="value" />
          </el-table>
        </el-col>/

      </el-row>
    </el-main>
  </el-container>
</template>

<script lang="ts">
import DateChoose from "./DateChoose.vue";
import {DownloadByAuthorId, DownloadByPid,DownloadByRank} from "../../wailsjs/go/main/App.js";
import {EventsEmit, EventsOn} from "../../wailsjs/runtime";
import { ref} from "vue";
export default {
  name: "maindownload",
  components: {DateChoose},
  props:{
    Input: Number,
    wait: Boolean,
  },
  setup(){
    const cellStyle = ({  rowIndex }) => {
      if (rowIndex === 0) {
        return {
          border: "3px solid #CD7F32",
        }
      }
    }
    const percent= ref(0);
    const tasknow= ref("No Task in queue");
    const queuenow=ref("There is no tasks waiting");
    const  now = ref("Pid")
    const queue=ref([
      {value:'tasknow'},
      {value: 'taskb'},
      {value:'taskc'},
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
    EventsOn("UpdateTaskNow",function(newmsg){
      console.log(newmsg)
      console.log(newmsg[0])
      tasknow.value=newmsg[0];
    });
    EventsOn("UpdateQueueNow",function(newmsg){
      console.log(newmsg[0])
      queuenow.value=newmsg[0];
    });
    EventsOn("UpdateProcess",function(newnum){

      console.log(newnum[0])
      percent.value=newnum[0];
    });
    EventsOn("Push",function(newmsg){
      queue.value.push(newmsg[0]);
    });
    EventsOn("Push",function(newmsg){
      queue.value.shift()
    });
    return{
      queue,
      tasknow,
      queuenow,
      percent,
      inputValue,
      period,
      now,
      options,
      modes,
      cellStyle,
    }
  }
  ,
  methods:{
    updateParentValue() {
      this.$emit('UpdateInput', this.inputValue);
    },
    changetype(data){
      console.log(data)
      this.value=data;
    },
    Download(){
      if(this.now=="Pid"){

        DownloadByPid(this.inputValue);
        this.inputValue="";
      }else if (this.now=="Author"){
        this.$emit("wait",true)
        DownloadByAuthorId(this.inputValue)
        this.inputValue="";
        this.$emit("wait",false)
      }else {
        this.$emit("wait",true)
        DownloadByRank(this.$refs.dateSelect.selectedDate,this.period)
        this.$emit("wait",false)
      }

    }
  },
}

</script>

<style lang="less" scoped>
@import "../assets/style/font.less";
@import "../assets/style/variable.less";
</style>