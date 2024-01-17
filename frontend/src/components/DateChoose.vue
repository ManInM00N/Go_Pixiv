<template>
  <div class="date-picker">
    <div class="block">
      <el-date-picker
          v-model="selectedDate"
          type="date"
          placeholder="Pick a day"
          :disabled-date="disabledDate"
          size="large"
          value-format="YYYY-MM-DD"
      />
    </div>
  </div>
</template>

<script lang="ts" >
import { ref} from 'vue'

// const selectedDate = ref('')
export default {
  data() {
    return {
      selectedDate:''
    }
  },
  methods:{
    disabledDate(time: Date) {
      return time.getTime() >  calculateNearestNoon();
    }
  }
}

// 计算距离现在最近的正午时间
const calculateNearestNoon = () => {
  const now = new Date();
  const noon = new Date(now);
  noon.setHours(12, 0, 0, 0);
  return now.getHours() >= 12 ? new Date(noon.getTime()-24*60*60*1000) : new Date(noon.getTime()- 48 * 60 * 60 * 1000);
};


</script>

<style scoped lang="less">
.date-picker {
  display: flex;
  width: 100%;
  padding: 0;
  flex-wrap: wrap;
  &.block{
    padding: 30px 0;
    text-align: center;
    border-right: solid 1px var(--el-border-color);
    flex: 1;
   }
  &.demonstration{
    display: block;
    color: var(--el-text-color-secondary);
    font-size: 14px;
    margin-bottom: 20px;
  }
}


.date-picker .block:last-child {
  border-right: none;
}

</style>
