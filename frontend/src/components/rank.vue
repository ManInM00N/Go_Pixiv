<template>
  <el-main>
    <el-row>
      <el-col :span="1"/>
      <el-col :span="6">
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
          />
        </el-select>
      </el-col>
      <el-col :span="6">
        <date-choose
            key="rank"
            ref="dateSelect"
            :re="true"
        ></date-choose>
      </el-col>
      <el-col :span="1"/>
      <el-col :span="8">
        <p>rank</p>
        <el-button @click="downloadthispage"/>
      </el-col>
    </el-row>
    <el-row>
      <el-space
          wrap
          size="30px"
      >
        <PicCard v-for="items in picitem"


        />
      </el-space>
    </el-row>


  </el-main>
</template>

<script lang="ts" setup>
import DateChoose from "./DateChoose.vue";
import PicCard from "./PicCard.vue";
import {DownloadByAuthorId, DownloadByPid,DownloadByRank} from "../../wailsjs/go/main/App.js";
import {EventsEmit, EventsOn} from "../../wailsjs/runtime";
import {defineComponent, onMounted, ref} from "vue";
import emitter from "../assets/js/Pub.js";
defineComponent({
  PicCard, DateChoose
})
name: "maindownload";

const picitem  =ref([
  {
    pid:"1234",
    type:"",
    src:"",
  },
  {
    pid:"1234",
    type:"",
    src:"",
  },
  {
    pid:"1234",
    type:"",
    src:"",
  },
  {
    pid:"1234",
    type:"",
    src:"",
  },
  {
    pid:"1234",
    type:"",
    src:"",
  },
])
const period=ref("daily");
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
const dateSelect= ref(null)
const downloadthispage = ()=>{
  console.log(dateSelect.value.selected.value)
  emitter.emit("DownloadByRank",{date:dateSelect.value.selected.value,period:period})
}

</script>


<style lang="less" scoped>


</style>