
<template >
  <el-aside
      name="menu" id="menu"
      class="asidemenu"
      style="width: 60px"
  >
    <el-menu
        :default-active="$route.path"
        :theme="theme"
        class="vertical-menu"
        mode="vertical"
        :router="true"
    >
      <el-container class="top-items" >
        <el-menu-item
            v-for="(item,idx) in items"
            :key="item.key"
            :id="item.id"
            :index="idx+''"
            :route="item.index"
            class="menu_item"
            @select="handleMenuSelect"
        >
          <el-container class="item_body">
            <el-icon size="30" class="item_icon">
              <component :is="item.iconmsg"  style="color: #32CD99"/>
            </el-icon>
          </el-container>
        </el-menu-item>
      </el-container >
      <el-container class="placeholder-item" ></el-container>
      <el-menu-item
          :key = "userself.key"
          :id = "userself.id"
          :index="userself.id.toString()"
          :route="userself.index"
          class="menu_item"
          style="padding-left: 10px;padding-right: 0px"
      >
        <el-container>
          <el-avatar shape="circle " src="/src/assets/images/no_profile.png" size="default"/>
        </el-container>
      </el-menu-item>
    </el-menu>

  </el-aside>
  <el-main
      class="View"
      id = "View"
      name="View"
      style="padding-right: 5px"
  >

    <router-view v-slot="{ Component,route }">
      <keep-alive v-if="route.meta.keepAlive ">
        <component :is="Component" />
      </keep-alive>
      <component v-else :is="Component" />
    </router-view>
<!--    <keep-alive>-->
<!--      <router-view v-if="$route.meta.keepAlive"></router-view>-->
<!--    </keep-alive>-->
<!--    <router-view v-if="!$route.meta.keepAlive"></router-view>-->
  </el-main>
</template>

<script setup>
import DateChoose from "./DateChoose.vue";
import {defineComponent, ref} from "vue";

const activeIndex=ref("/maindownload")
const theme= ref('dark')
const items=ref([
  {id: 1,iconmsg :"HomeFilled",key:"maindownload" ,index:"/maindownload"},
  {id: 2,iconmsg :"StarFilled",key:"follow",index:"/follow"},
  {id: 3,iconmsg :"Histogram",key:"rank",index:"/rank"},
  {id: 4,iconmsg :"Search",key:"search",index:"/search"},
  {id: 5,iconmsg :"Setting",key:"setting",index:"/setting"},
])
const userself=ref({
  id:6,key: "user",index:"/user"
})
const Input=ref('')
const wait=ref(false)
function waitchange(val){
  wait.value=val
}
function  handleMenuSelect(index) {
  console.log("ee",this.activeIndex)
}
</script>

<style lang="less" scoped>
@import "src/assets/style/Color.less";
@import "src/assets/style/menu.less";
</style>