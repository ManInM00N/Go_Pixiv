<template>
  <el-container>
    <el-main>
      <el-row >
        <el-col :span="16"/>
        <el-col :span="8">
          下载模式
          <el-select
            v-model="value"
            ref="mode"
            class="m-2"
            size="default"
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
          <el-input  ref="re" size="large" placeholder="Pid/AuthorId"></el-input>
        </el-col>
        <el-col :span="8" v-if="this.value==='Rank'">
          <el-select
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
            />
          </el-select>
        </el-col>
        <el-col :span="8" v-if="this.value==='Rank'">
          <date-choose></date-choose>
        </el-col>
        <el-col :span="2"/>
        <el-col :span="6">
          <el-button
              :disabled="true"
              style=""
              :icon="Download"
              type="success"
              size="large"
          >
            Download</el-button>
        </el-col>
      </el-row >
      <el-row >
      </el-row>
    </el-main>
  </el-container>
</template>

<script  >
import { ref } from 'vue'
import DateChoose from "./DateChoose.vue";
export default {
  name: "maindownload",
  components: {DateChoose},
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
      value:"Pid"
    }
  },
}
</script>

<style lang="less" scoped>
@import "../assets/style/variable.less";
</style>