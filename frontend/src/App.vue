
<template>
  <el-container >
    <el-aside
        name="menu" id="menu"
        class="menu"
        style="height:100%; width: 60px"
    >
      <el-menu
          style=" height: 100%"
          :default-active="'1'"
          :selectedKeys="[current]"
          :theme="theme"
          @click="handleClick"
      >
        <el-menu-item
            v-for="item in items"
            :key="item.id.toString()"
            :id="item.id"
            class="menu_item"
        >
          <el-icon size="30px" ></el-icon>
        </el-menu-item>
        <el-menu-item
            :key = "items.length+1"
            class="menu_item"
        >
          <el-icon size="30px"><HomeFilled /></el-icon>
        </el-menu-item>
      </el-menu>
      <el-footer>

      </el-footer>
    </el-aside>
    <el-main class="View" id = "View" name="View">
      <DateChoose></DateChoose>
    </el-main>
<!--    <el-footer></el-footer>-->
  </el-container>

</template>
<script>
import DateChoose from "./components/DateChoose.vue";

export default {
  components:{
    // Tab
    DateChoose
  },
  data () {
    return {
      current: '1',
      theme: 'dark',
      items:[
        {id: 1,iconmsg :"<HomeFilled />"},
        {id: 2,iconmsg :"<Menu />"},
        {id: 3,iconmsg :"<Menu />"},
      ]
    }
  },
  methods: {
    handleClick (e) {
      console.log('click ', e)
      this.current = e.key
    },
    changeTheme ({ target }) {
      this.theme = target.checked ? 'light' : 'dark'
    },
  },
};
</script>
<style lang='less' scoped>
.menu{
  width: 60px;
  height: 100%;
}
.menu_item{
  display: flex;
  align-items: center;
  //padding-left: 10px;
}
.View{
  width: calc(100% - 60px);
  height: 100%;
  display: flex;
  position: relative;
}
.Main{
  display: block;
  width: 100%;
  height: 100%;
}
</style>