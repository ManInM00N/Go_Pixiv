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
        <el-col :span="8" v-if="this.value!=='Rank'">
        </el-col>
        <el-col :span="8" v-if="this.value!=='Rank'">
          <el-input
              v-model="inputValue"
              size="large"
              placeholder="Pid/AuthorId"
              clearable
              type="number"
              @input="updateParentValue"
          ></el-input>
        </el-col>
        <el-col :span="8" v-if="this.value==='Rank'">
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
        <el-col :span="8" v-if="this.value==='Rank'">
          <date-choose></date-choose>
        </el-col>
        <el-col :span="2"/>
        <el-col :span="6">
          <el-button
              style=""
              id="bt"
              type="success"
              size="large"
              @click="Download"
          >
            Download
            <el-icon size="large">
            <Download/>
            </el-icon>
          </el-button>
        </el-col>
      </el-row >

      <el-row >
        <el-col :span="6">
          <el-text>
            {{tasknow}}
          </el-text>
        </el-col>
        <el-col :span="3" />
        <el-col :span="6">
          <el-progress
              :stroke-width="24"
              striped
              striped-flow
              :duration="10"
              :percentage="percent"
          >

          </el-progress>
        </el-col>
        <el-col :span="3"/>
        <el-col :span="6">
          <el-text>
            {{queuenow}}
          </el-text>
        </el-col>
      </el-row>
    </el-main>
  </el-container>
</template>

<script>
import DateChoose from "./DateChoose.vue";
import {DownloadByAuthorId, DownloadByPid} from "../../wailsjs/go/main/App.js";
import {EventsEmit, EventsOn, WindowReloadApp} from "../../wailsjs/runtime/runtime.js";
export default {
  name: "maindownload",
  components: {DateChoose},
  props:{
    Input: Number,
    wait: Boolean,
  },
  data(){
    return{
      now: "Pid",
      options:[
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
      ],
      modes:[
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
      ],
      inputValue: this.Input,
      period:"daily",
      percent: 0,
      tasknow:"No Task in queue",
      queuenow:"There is no tasks waiting",
    }
  },
  watch:{
    "UpdateTaskNow":function (newmsg) {
      this.tasknow=newmsg;
    },
    "UpdateQueueNow":function (newmsg) {
      this.queuenow=newmsg;
    },
    "UpdateProcess":function(newnum){
      this.percent=newnum;
    }
  },
  methods:{
    updateParentValue() {
      this.$emit('UpdateInput', this.inputValue);
    },
    changetype(data){
      console.log(data)
      this.value=data;
    },
    Download(){
      EventsOn("UpdateTaskNow",function(newmsg){
        this.tasknow=newmsg;
      });
      EventsOn("UpdateTaskNow",function(newmsg){
        this.queuenow=newmsg;
      });
      EventsOn("UpdateTaskNow",function(newnum){
        this.percent=newnum;
      });
      if(this.now=="Pid"){

        DownloadByPid(this.inputValue);
        this.inputValue="";
      }else if (this.now=="Author"){


        EventsEmit("wait",true)
        DownloadByAuthorId(this.inputValue)


        this.inputValue="";
        EventsEmit("wait",false)
      }else {

      }

    }
  },

}
</script>

<style lang="less" scoped>
@import "../assets/style/variable.less";
</style>