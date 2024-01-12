<template>
  <el-aside
      name="menu" id="menu"
      class="asidemenu"
      style="width: 60px"
  >
    <el-menu
        :default-active="activeIndex"
        :selectedKeys="[current]"
        :theme="theme"
        class="vertical-menu"
        mode="vertical"
        :router="true"
        default-openeds="maindownload"
    >
      <el-container class="top-items">
        <el-menu-item
            v-for="item in items"
            :key="item.key"
            :id="item.id"
            :index="'/'+item.key"
            class="menu_item"
            @select="handleMenuSelect"
            index="/MainDownload"
        >
          <el-container class="item_body">
            <el-icon class="Icon">
              <component :is="item.iconmsg"/>
            </el-icon>
          </el-container>
        </el-menu-item>
      </el-container>
      <el-container class="placeholder-item" ></el-container>
      <el-menu-item
          :key = "userself.key"
          :id = "userself.id"
          :index="'/'+userself.key"
          class="menu_item"
      >
        <el-container>
          <el-image class="headImg" src="frontend/assets/images/no_profile.png" />
        </el-container>
      </el-menu-item>
    </el-menu>

  </el-aside>
  <el-main class="View" id = "View" name="View">
    <router-view></router-view>
    <DateChoose/>
  </el-main>
</template>

<script>
import DateChoose from "./DateChoose.vue";
export default {
  name: "Menu_List",
  components:{
    DateChoose,
  },
  data () {
    return {
      activeIndex:"/maindownload",
      current: '1',
      theme: 'dark',
      items:[
        {id: 1,iconmsg :"HomeFilled",key:"maindownload"},
        {id: 2,iconmsg :"StarFilled",key:"follow"},
        {id: 3,iconmsg :"Histogram",key:"rank"},
        {id: 4,iconmsg :"Search",key:"search"},
        {id: 5,iconmsg :"Setting",key:"setting"},
      ],
      userself: {
        id:6,key: "user"
      }
    }
  },
  methods: {
    handleMenuSelect(index) {
      this.activeIndex = index;
      console.log("ee",this.activeIndex)
    },
    changeTheme ({ target }) {
      this.theme = target.checked ? 'light' : 'dark'
    },
  },
}
</script>

<style lang="less" scoped>
@import "src/assets/style/menu.less";
</style>