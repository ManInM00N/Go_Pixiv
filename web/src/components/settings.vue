<template>
    <el-main>
        <el-form label-width="100px" :model="form" style="max-width: 500px" class="Background_Opacity">
            <el-form-item label="Localhost:" label-width=150>
                <el-input v-model="form.prefix" disabled="true" />
            </el-form-item>
            <el-form-item label="代理端口" label-width=150>
                <el-input v-model="form.proxy" :disabled="!form.useproxy">
                </el-input>
            </el-form-item>
            <el-form-item label="Cookie" label-width=150>
                <el-input v-model="form.cookie">
                </el-input>
            </el-form-item>
            <el-form-item label="下载位置" label-width=150>
                <el-input v-model="form.downloadposition" />
            </el-form-item>
            <el-form-item>

            </el-form-item>
            <el-form-item label="请求自省时间(ms)" label-width=150>
                <el-input type="number" v-model="form.retry429">
                    <template #append>ms</template>
                </el-input>
            </el-form-item>
            <el-form-item label="失败重试间隔(ms)" label-width=150>
                <el-input type="number" v-model="form.retryinterval">
                    <template #append>ms</template>
                </el-input>
            </el-form-item>
            <el-form-item label="下载间隔(ms)" label-width=150>
                <el-input type="number" v-model="form.downloadinterval">
                    <template #append>ms</template>
                </el-input>
            </el-form-item>
            <el-form-item label="作品收藏数限制" label-width=150>
                <el-input type="number" oninput="value=value.replace(/[-]/g,'')" v-model="form.likelimit" />
            </el-form-item>
          <el-form-item label="缓存过期时间" label-width=150>
            <el-input type="number" v-model="form.expired_time" >
              <template #append>day(s)</template>
            </el-input>
          </el-form-item>
            <el-form-item>
                <el-checkbox label="是否启用R-18" v-model="form.r_18" />
                <el-checkbox label="下载图片对作者分类" v-model="form.differauthor" />
                <el-checkbox label="是否启用本地代理" v-model="form.useproxy" />
            </el-form-item>
            <el-form-item>
                <el-button @click="UpLoad">
                    Update
                </el-button>
            </el-form-item>
        </el-form>
    </el-main>
</template>

<script lang="ts" setup>
import { Events } from "@wailsio/runtime";
import { ElNotification } from "element-plus";
import {form, updateSettings} from "../assets/js/configuration.js"
import axios from "axios";
import { CheckLogin } from "../../bindings/main/internal/pixiv/ctl.js";

function UpLoad() {
    console.log(form, form.value)
    form.value.logined = false
    ElNotification({
        type: "info",
        title: "INFO",
        message: "Login ......",
        position: 'bottom-right',
        duration: 3000,
    })
    updateSettings()

}
</script>

<style lang="less" scoped>
@import "../assets/style/variable.less";
</style>
