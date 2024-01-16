
<template>
  <el-container
    class="Main"
  >
    <Menu_List/>
  </el-container>

</template>
<script>
import Menu_List from "./components/Menu_List.vue";
export default {
  components:{
    Menu_List,
  },
  data(){
    return{

    }
  }
};
</script>
<style lang='less' scoped>
@import "assets/style/variable.less";


</style>