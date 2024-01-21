<template>
  <el-main>
    <el-row>
      <el-col>
        <el-text
            type="warning"
        >
          <h1
              style="text-align:center;font-size:30px"
          >
            已关注用户的作品
          </h1>
        </el-text>
      </el-col>
    </el-row>
    <Waterfall
        :list="picitem"
        width="300"
        background-color=""
        animation-effect="fadeInUp"
    >
      <template #item="{ item, url, index }" >
        <div class="card" >
          <!--            <LazyImg :url="url" />-->
          <!--            <p class="text">这是具体内容</p>-->
          <PicCard
              :author="item.Author"
              :img="item.src"
              :title="item.Title"
              :pid="item.pid"
              :authorId="item.authorId"
              :pages="item.pages"
          />
        </div>
      </template>
    </Waterfall>

  </el-main>
</template>

<script setup>
import {ref} from "vue";
import {EventsOn} from "../../wailsjs/runtime/runtime.js";
const picitem = ref([])

const name=ref("follow")

EventsOn("UpdateLoad",function(msg){
  console.log(msg[0])
  // picitem.value = [...picitem.value,{pid:msg[0],Title:msg[1],Author: msg[2],src: "cache/images/"+msg[0]+".jpg",pages:msg[3],authorId:msg[4]}]
  // picitem.value.push({pid:msg[0],Title:msg[1],Author: msg[2],src: "cache/images/"+msg[0]+".jpg",pages:msg[3],authorId:msg[4]})
  picitem.value=picitem.value.concat({pid:msg[0],Title:msg[1],Author: msg[2],src: "cache/images/"+msg[0]+".jpg",pages:msg[3],authorId:msg[4]})
})

</script>

<style scoped>

</style>