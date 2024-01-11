<template>
  <div class="demo-date-picker">
    <div class="block">
      <span class="demonstration">Picker with quick options</span>
      <el-date-picker
          v-model="value2"
          type="date"
          placeholder="Pick a day"
          :disabled-date="disabledDate"
          :shortcuts="shortcuts"
          size="large"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'

const value2 = ref('')

const shortcuts = [
  {
    text: 'Today',
    value: new Date(),
  },
  {
    text: 'Yesterday',
    value: () => {
      const date = new Date()
      date.setTime(date.getTime() - 3600 * 1000 * 24)
      return date
    },
  },
  {
    text: 'A week ago',
    value: () => {
      const date = new Date()
      date.setTime(date.getTime() - 3600 * 1000 * 24 * 7)
      return date
    },
  },
]

const disabledDate = (time: Date) => {
  return time.getTime() > Date.now()
}
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
