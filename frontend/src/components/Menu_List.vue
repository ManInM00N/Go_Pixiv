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
            :disabled="!item.logined"
            :limit="form['r-18']"
        >
          <el-container
              class="item_body"

          >
            <el-icon size="30" class="item_icon">
              <component
                  :is="item.iconmsg"
                  style="color: #32CD99"
              />
            </el-icon>
          </el-container>
        </el-menu-item>
      </el-container >
      <el-container class="placeholder-item " ></el-container>
      <el-menu-item
          :key = "userself.key"
          :id = "userself.id"
          :index="userself.id.toString()"
          :route="userself.index"
          class="menu_item"
          @select="handleMenuSelect"
      >
        <el-container>
          <el-icon
              size="30" class="item_icon"
          >
            <Tools/>
          </el-icon>
<!--          <el-avatar shape="circle " src="/src/assets/images/no_profile.png" size="default"/>-->
        </el-container>
      </el-menu-item>
    </el-menu>

  </el-aside>
  <el-main
      class="View"
      id = "View"
      name="View"
      style="padding-right: 5px;padding-left: 5px"
  >
    <router-view v-slot="{ Component,route,Path }">
      <keep-alive v-if="route.meta.keepAlive ">
        <component
            :is="Component"
            :form="form"
        />
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
import {defineComponent, onMounted, ref} from "vue";
import settings from "./settings.vue";
import {CheckLogin, GetSetting, UpdateSetting} from "../../wailsjs/go/main/App.js";
import emitter from "../assets/js/Pub.js"
import {EventsOn} from "../../wailsjs/runtime/runtime.js";
import {DAO} from "../../wailsjs/go/models.ts";
const form = ref(DAO.Settings)
const activeIndex=ref("/maindownload")
const theme= ref('dark')
const set = ref(null)
const items=ref([
  {id: 1,iconmsg :"HomeFilled",key:"maindownload" ,index:"/maindownload",logined: true},
  {id: 2,iconmsg :"StarFilled",key:"follow",index:"/follow",logined: false},
  {id: 3,iconmsg :"Histogram",key:"rank",index:"/rank",logined: true},
  {id: 4,iconmsg :"Search",key:"search",index:"/search",logined: true},
])
const userself=ref({
  // id:6,key: "user",index:"/user"
  id: 5,iconmsg :"Tools",key:"settings",index:"/setting",logined: true
})
const Input=ref('')
const wait=ref(false)
function waitchange(val){
  wait.value=val
}
function  handleMenuSelect(index) {
  console.log("ee",this.activeIndex)
}
emitter.on("Relogin",function(){
  // logined()
  CheckLogin()
})
EventsOn("login",function(msg){
  if(msg==="True"){
    items.value[1].logined=true
  }else {
    items.value[1].logined=false
    form.value.cookie=""
    const temp = GetSetting()
    console.log(temp)
    temp.then(res=>{
      form.value=res
      console.log(form.value)
    })
  }
  console.log(items.value[1].logined,form.value['r-18'],!(!form.value['r-18']))
})
onMounted(function(){
  CheckLogin()
  const temp = GetSetting()
  console.log(temp)
  temp.then(res=>{
    form.value=res
    console.log(form.value)
  })
})

</script>

<style lang="less" scoped>
@import "src/assets/style/Color.less";
@import "src/assets/style/menu.less";
.tempheight{
  height:60px
}
</style>