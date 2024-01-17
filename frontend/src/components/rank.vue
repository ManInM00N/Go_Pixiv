<template>
  <el-container>
    <p>rank</p>
    <el-button @click="Temp"/>
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
    const percent= ref(0);
    const tasknow= ref("No Task in queue");
    const queuenow=ref("There is no tasks waiting");
    return{
      tasknow,
      queuenow,
      percent,
    }
  }
  ,
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

EventsOn("UpdateTaskNow",function(newmsg){
  console.log(newmsg)
  console.log(newmsg[0])
  this.$refs.tasknow=newmsg[0];
});
EventsOn("UpdateQueueNow",function(newmsg){
  console.log(newmsg[0])
  this.$refs.queuenow=newmsg[0];
});
EventsOn("UpdateProcess",function(newnum){

  console.log(newnum[0])
  this.$refs.process=newnum[0];
});
</script>


<style scoped>

</style>