
<template>
  <el-container
    class="Main"
  >
    <Menu_List/>
    <router-view/>
  </el-container>

</template>
<script>
import Menu_List from "./components/Menu_List.vue";
export default {
  components:{
    Menu_List,
  },

};
</script>
<style lang='less' scoped>
@import "assets/style/variable.less";
.Main{
  height: calc(100vh - 20px);
}

</style>