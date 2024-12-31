<template>
    <el-container style="margin-left: 15px;">
        <el-main style="
      display: flex;
        flex-direction: column;
        height: 90%;
      ">
            <div style="width: 100%;height: 150px;">
                <el-tabs type="card">
                    <el-tab-pane label="Pid/AuthorId">
                    </el-tab-pane>
                    <el-tab-pane label="Rank">
                    </el-tab-pane>
                    <el-tab-pane label="Novel">
                    </el-tab-pane>
                </el-tabs>
            </div>
            <el-row>
                <el-col :span="16" />
                <el-col :span="8">
                    下载模式
                    <el-select v-model="now" ref="mode" class="m-2" size="default" @change="changetype"
                        style="width:150px">
                        <el-option v-for="(item) in modes" :key="item.value" :label="item.label" :value="item.value" />
                    </el-select>
                </el-col>
            </el-row>
            <el-divider class="Divide_Line" style="color:#52616b" />
            <el-row>
                <el-col :span="6">
                    <el-input v-model="inputValue" size="large" placeholder="Pid/AuthorId" clearable
                        type="number"></el-input>
                </el-col>
                <el-col :span="2" />
                <el-col :span="8">
                    <el-select v-model="period" ref="mode2" class="m-2" size="large" style="width:150px"
                        @change="changetype2">
                        <el-option v-for="(item) in options" :key="item.value" :label="item.label" :value="item.value"
                            :disabled="!item.disabled" />
                    </el-select>
                </el-col>
                <el-col :span="8">
                    <date-choose key="main" ref="dateSelect"></date-choose>
                </el-col>
            </el-row>
            <el-row style="height: 20px" />
            <el-row>
                <el-col :span="1" />
                <el-col :span="15" class="Tre">
                    <el-progress stroke-width="28" striped striped-flow :duration="10" :percentage="percent"
                        text-inside="true" style="color:black">
                        <span style="font-size:16px;color:#c9d6df" v-if="queue.length > 0">{{ queue[0].value }} {{
                            percent
                            }}%</span>
                    </el-progress>
                </el-col>
                <el-col :span="2" />
                <el-col :span="6">
                    <el-button style="" id="bt" type="success" size="large" @click="Download" :disabled="wait">
                        Download
                        <el-icon size="large">
                            <Download />
                        </el-icon>
                    </el-button>
                    <el-text class="Tre">
                        {{ queuenow }}
                    </el-text>
                </el-col>
            </el-row>
            <br>
            <el-row class="Get_Remain">
                <el-col :span="1" />
                <el-col :span="15" class="terminal-text">
                    <el-scrollbar class="text Micro" style="height:500px">
                        <p v-for="item in logs">
                            {{ item }}
                        </p>
                    </el-scrollbar>
                </el-col>
                <el-col :span="1" />
                <el-col :span="7">
                    <el-table :data="queue" :cell-class-name="cellStyle" class="Half_light Tre queueTable"
                        style="height:490px">
                        <template #empty>
                            No Task
                        </template>
                        <el-table-column label="TaskQueue" prop="value" style="height:60px;">
                            <template #default="scope">
                                <span style="color:#c9d6df">{{ scope.row.value }}</span>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-col>
            </el-row>
        </el-main>
        <el-footer style="height: 10%;">
            <el-text type="danger" id="time">
                此软件为免费开源，如果是购买获得请退款举报 {{ timeElement }}
            </el-text>
        </el-footer>
    </el-container>
</template>

<script lang="ts" setup>
import DateChoose from "./DateChoose.vue";
import { onMounted, ref } from "vue";
import { ElMessage } from "element-plus";
import axios from "axios";
import { timeElement } from "../assets/js/Time.js"
import { Events } from "@wailsio/runtime"
const dateSelect = ref(null)
import { ws, form } from "../assets/js/configuration.js";

onMounted(() => {
    ws.value.onmessage = (event) => {
        // res.value = event.data;
        handleMessage(JSON.parse(event.data));
    };
})
const logs = ref([]);
const rows = ref(10); // 可根据需要调整展示的行数
const cellStyle = ({ rowIndex }) => {
    if (rowIndex === 0) {
        return 'Xbord'
    }
}
function handleMessage(data) {
    if (data.type == 1) {
        percent.value = data.newnum
    } else if (data.type == 2) {
        queue.value.shift()
    } else if (data.type == 3) {
        queue.value.push(data.newtask)
    }
}
const mode = ref('')
function changetype(data) {
    console.log(data)
    now.value = data
}
function changetype2(data) {
    console.log(data)
    period.value = data
}
const percent = ref(0);
const now = ref("Pid")
const queue = ref([
])
const options = ref([
    {
        value: "daily",
        label: "Daily",
        disabled: true
    },
    {
        value: "weekly",
        label: "Weekly",
        disabled: true
    },
    {
        value: "monthly",
        label: "Monthly",
        disabled: true
    },
    {
        value: "daily_r18",
        label: "Daily_R18",
        disabled: form.value.r_18,
    },
    {
        value: "weekly_r18",
        label: "Weekly_R18",
        disabled: form.value.r_18,
    },
]);
const modes = ref([
    {
        value: "Pid",
        label: "By Pid",
    },
    {
        value: "Author",
        label: "By AuthorId",
    },
    {
        value: "Rank",
        label: "By Rank",
    },
]);
const inputValue = ref('');
const period = ref("daily");

Events.On("UpdateProcess", function (newnum) {
    // console.log(newnum, newnum.data[0][0])
    percent.value = newnum.data[0][0];
});
Events.On("Push", function (newmsg) {
    console.log(newmsg, newmsg.data[0][0])
    queue.value.push({ value: newmsg.data[0][0] })
});
Events.On("Pop", function () {
    queue.value.shift()
});
Events.On("UpdateTerminal", function (newmsg) {
    console.log(newmsg)
    logs.value.push(newmsg.data[0][0])
    if (logs.value.length > 50) {
        logs.value.pop()
    }
})
function Download() {
    console.log("Downloading ", now.value)
    if (now.value != "Rank" && inputValue.value === '') {
        return
    }
    ws.value.send(
        JSON.stringify({
            type: now.value,
            id: inputValue.value,
            period: period.value,
            time: dateSelect.value.selectedDate
        })
    )
    inputValue.value = ''
    return
}

</script>
<style lang="less" scoped>
@import "../assets/style/font.less";
@import "../assets/style/variable.less";
@import "../assets/style/color.less";

.terminal-text {
    background-color: rgba(0, 0, 0, 0.3);
    font-family: monospace;
    border: none;
    white-space: pre-wrap;
}

.text {
    background: rgba(@quartz, 0.1);
}

.queueTable {
    -webkit-background-clip: text;
    //opacity: 0.5;
}

.No_Background {
    background: rgba(ff, ff, ff, 0.3);
    border: 2px solid #CD7F32;
}

/deep/ .el-table {
    thead {
        color: #fff;
        font-weight: 500;
        background: linear-gradient(to right, rgba(#dde7f2, 0.5), rgba(#878ecd, 0.5)) !important;

        & th {
            background-color: transparent;
        }

        & tr {
            background-color: transparent;
        }
    }
}
</style>
