
<template>
  <div ref="totop" class="totop" @click="scrollToTop" title="返回顶部">
    <svg t="1681820716964" class="icon" style="display:block" viewBox="0 0 1024 1024" version="1.1"
         xmlns="http://www.w3.org/2000/svg" p-id="6409" width="32" height="32">
      <path
          d="M866.7 96H157.3c-17.8 0-32.2 14.4-32.2 32.2 0 17.8 14.4 32.2 32.2 32.2h709.5c17.8 0 32.2-14.4 32.2-32.2 0-17.8-14.4-32.2-32.3-32.2zM512 225c-16.5 0-33 6.3-45.6 18.9L169.9 540.3c-5.8 5.8-9.4 13.9-9.4 22.8s3.6 17 9.4 22.8c5.8 5.8 13.9 9.4 22.8 9.4s17-3.6 22.8-9.4l264.2-264.2v574c0 17.8 14.4 32.2 32.2 32.2 17.8 0 32.2-14.4 32.2-32.2v-574l264.2 264.2c5.8 5.8 13.9 9.4 22.8 9.4s17-3.6 22.8-9.4c5.8-5.8 9.4-13.9 9.4-22.8s-3.6-17-9.4-22.8L557.6 243.9C545 231.3 528.5 225 512 225z"
          fill="#ffffff" p-id="6410"></path>
    </svg>
  </div>
</template>

<script setup lang="ts">
import {defineComponent, onMounted, ref,toRef} from "vue";

const totop = ref();

const scrollToTop = () => {
  window.scrollTo({
    // top: document.documentElement.offsetHeight, //回到底部
    top: 0, //回到顶部
    left: 0,
    behavior: "smooth", //smooth 平滑；auto:瞬间
  });
};

onMounted(() => {
  // 页面滚动窗口监听事件
  window.onscroll = function () {
    // 获取浏览器卷去的高度
    let high = document.documentElement.scrollTop || document.body.scrollTop; //兼容各浏览器
    if (high >= 500) {
      totop.value.style.display = "block";
    } else {
      totop.value.style.display = "none";
    }
  };
});
</script>

<style scoped lang="less">
.totop {
  display: none;
  position: fixed;
  bottom: 120px;
  right: 80px;
  background-color: #b0b0b0;
  padding: 10px;
  border-radius: 50%;

  &:hover {
     cursor: pointer;
     background-color: #00000088;
   }
}
</style>