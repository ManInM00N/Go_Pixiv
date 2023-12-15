<script setup>
import HelloWorld from './components/HelloWorld.vue'</script>
<template>
  <div class="app-container">
    <!-- 标签切换 -->
    <div class="tags_box">
      <div v-for="(item,index) in messageList" :key="item.id" class="tags">
        <div :class="{isActive:index==active}" @click="handelToogel(index)">{{item}}</div>
      </div>
    </div>
    <!-- 内容部分 -->
    <div class="con_box">
      <div v-show="active == 0">内容一</div>
      <div v-show="active == 1">内容二</div>
      <div v-show="active == 2">内容三</div>
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      messageList: ["标签一", "标签二", "标签三"],
      active: 0,
    };
  },
  methods: {
    handelToogel(index) {
      this.active = index;
    },
  },
};
</script>
<style lang="less" scoped>
.tags_box {
  width: 50%;
  margin: 20px auto;
  display: flex;
}
.tags div {
    padding: 10px 20px;
    cursor: pointer;
    margin: 0 15px;
}
.isActive {
  border-bottom: 2px solid #2ed573;
}
.con_box{
  width: 50%;
  margin: 20px auto;
}

</style>