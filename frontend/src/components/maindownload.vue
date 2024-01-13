<template>
  <el-container>
    <el-main>
      <el-row>
        <el-col :span="16"/>
        <el-col :span="8" >
          选择下载模式
          <el-select
            v-model="value"
            class="m-2"
            size="default"
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
      </el-row>
      <el-row>
        <el-col>

        </el-col>
      </el-row>
      <el-row>
        <el-col>

        </el-col>
      </el-row>
      <el-row>
        <el-col>

        </el-col>
      </el-row>
    </el-main>
  </el-container>
</template>

<script  >
import { ref } from 'vue'
export default {
  name: "maindownload",

  data(){
    return{
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
          label:"By Author",
        },
        {
          value:"Rank",
          label:"By Rank",
        },
      ],
      value:"daily"
    }
  },
}
</script>

<style scoped>

</style>