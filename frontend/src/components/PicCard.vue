<template

>
<!--    -->
  <el-skeleton
      style="width: 240px"
      :loading="false"
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
      <el-card
          :body-style="{ padding: '0px', marginBottom: '1px',width: '100%' }"
          v-if="!($props.r18==='r18'&&!$props.limit)"
      >
        <img
            :src="$props.img"
            class="image"
        />
        <div style="padding: 14px">
          <el-row>
            <el-text
                class="w-280px mb-2"
                truncated
                @click="jump('https://www.pixiv.net/artworks/'+$props.pid)"
            >
                {{ $props.title }}
            </el-text>
          </el-row>
          <el-row>
              <el-text
                  class="w-280px mb-2"
                  truncated
                  type="primary"
                  @click="'https://www.pixiv.net/users/'+$props.authorId"
              >{{$props.author}}</el-text>
          </el-row>
          <el-row>
            <el-col :span="20">
              <el-text class="w-250px mb-2" truncated type="success">
                Pages:{{$props.pages}}
              </el-text>
            </el-col>

            <el-col :span="4">
              <div class="bottom card-header                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    ">
                <el-button
                    text class="button"
                    @click="download"
                    :disabled="false"
                >
                  <el-icon size="30"><Download/></el-icon>
                </el-button>
              </div>
            </el-col>
          </el-row>

        </div>
      </el-card>
    </template>
  </el-skeleton>
</template>

<script lang="ts" setup>
import {ref} from "vue";
import emitter from "../assets/js/Pub.js";
const name= "PicCard"
const props = defineProps({
  limit:{
    type: Boolean,
    default:true,
  }
  ,
  pid :{
    type:String,
  },
  author:{
    type:String,
  },
  title: {
    type: String,
    default: "确定",
  },
  img: {
    type: String,
    default: "",
  },
  pages:{
    type: String,
    default: 1,
  },
  authorId:{
    type: String,
    default:1,
  },
  r18:{
    type: String,
    default:"r18",
  }
});

const download = ()=>{
  console.log("trying to download" ,props.pid,props.pid.value)
  emitter.emit("DownloadByPid",{pid:props.pid})
}
defineExpose({

})
function jump(event) {
  console.log("jump" ,event)
  window.open(event, '_blank')
}
</script>

<style lang="less" scoped>
.image{
  width:100%;
  height:100%;

}
</style>