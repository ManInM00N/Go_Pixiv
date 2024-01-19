<template>
  <el-skeleton
      style="width: 240px"
      :loading="true"
      animated
      :throttle="500"
  >
    <template #template>
      <el-skeleton-item variant="image" style="width: 240px; height: 240px" />
      <div style="padding: 14px">
        <el-skeleton-item variant="h3" style="width: 50%" />
        <div
            style="
              display: flex;
              align-items: center;
              justify-items: stretch;
              margin-top: 16px;
              height: 16px;
            "
        >
          <el-skeleton-item variant="text" style="margin-right: 16px" />
          <el-skeleton-item variant="text" style="width: 30%" />
        </div>
      </div>
    </template>
    <template #default>
      <el-card :body-style="{ padding: '0px', marginBottom: '1px' }">
        <img
            :src="src"
            class="image"
        />
        <div style="padding: 14px">
          <span>Delicious hamburger</span>
          <div class="bottom card-header                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    ">
            <el-button text class="button" @click="download">
              <el-icon><Download/></el-icon>
            </el-button>
          </div>
        </div>
      </el-card>
    </template>
  </el-skeleton>
</template>

<script lang="ts" setup>
import {ref} from "vue";
import emitter from "../assets/js/Pub.js";
const name= "PicCard"
const pid = ref('')
const src = ref('')
const title = ref('')
defineProps({
  title: {
    type: String,
    default: "确定",
  },
  icons: {
    type: String,
    default: "",
  },
});
const download = ()=>{
  emitter.emit("DownloadByPid",{pid:pid})
}

</script>

<style lang="less" scoped>

</style>