<template>
    <div class="date-picker">
        <div class="block">
            <el-date-picker v-model="selectedDate" format="YYYY-MM-DD" value-format="YYYY-MM-DD" type="date"
                :disabled="lock" placeholder="Pick a day" :disabled-date="disabledDate" size="large"
                :clearable="false" />
        </div>
    </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { padStart } from "../assets/js/Time.js"
import emitter from "../assets/js/Pub.js";

const props = defineProps({
    re: {
        type: Boolean,
        default: false,
    },
    lock: {
        type: Boolean,
        default: false,
    }
});
function disabledDate(time: Date) {
    return time.getTime() > calculateNearestNoon();
}
const calculateNearestNoon = () => {
    const now = new Date();
    const noon = new Date(now);
    noon.setHours(12, 0, 0, 0);
    return now.getHours() >= 12 ? new Date(noon.getTime() - 24 * 60 * 60 * 1000) : new Date(noon.getTime() - 48 * 60 * 60 * 1000);
};
const selectedDate = ref(formatDateToYYYYMMDD(calculateNearestNoon()))
function formatDateToYYYYMMDD(date) {
    const year = date.getFullYear();
    const month = padStart(date.getMonth() + 1, 2, '0');
    const day = padStart(date.getDate(), 2, '0');
    return `${year}-${month}-${day}`;
}
const emitsEventList = defineEmits(["DownloadByRank"])

const re = ref(props.re)
defineExpose(
    { selectedDate }
)
</script>

<style scoped lang="less">
.date-picker {
    display: flex;
    width: 100%;
    padding: 0;
    flex-wrap: wrap;

    &.block {
        padding: 30px 0;
        text-align: center;
        border-right: solid 1px var(--el-border-color);
        flex: 1;
    }

    &.demonstration {
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
