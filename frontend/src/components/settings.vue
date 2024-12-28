<template>
    <el-main>
        <el-form label-width="100px" :model="form" style="max-width: 460px" class="Background_Opacity">
            <el-form-item label="Localhost:">
                <el-input v-model="form.prefix" disabled="true" />
            </el-form-item>
            <el-form-item label="代理端口">
                <el-input v-model="form.proxy">
                </el-input>
            </el-form-item>
            <el-form-item label="Cookie">
                <el-input v-model="form.cookie">
                </el-input>
            </el-form-item>
            <el-form-item label="下载位置">
                <el-input v-model="form.downloadposition" />
            </el-form-item>
            <el-form-item>

            </el-form-item>
            <el-form-item label="请求自省时间(ms)">
                <el-input type="number" v-model="form.retry429" />
            </el-form-item>
            <el-form-item label="请求失败重试间隔(ms)">
                <el-input type="number" v-model="form.retryinterval" />
            </el-form-item>
            <el-form-item label="下载间隔(ms)">
                <el-input type="number" v-model="form.downloadinterval" />
            </el-form-item>
            <el-form-item label="作品收藏数限制">
                <el-input type="number" oninput="value=value.replace(/[-]/g,'')" v-model="form.likelimit" />
            </el-form-item>
            <el-form-item>
                <el-checkbox label="是否启用R-18" v-model="form.r_18" />
                <el-checkbox label="下载图片对作者分类" v-model="form.differauthor" />
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

import { form } from "../assets/js/configuration.js"
import { defineComponent, onMounted, ref, reactive } from "vue";
import emitter from "../assets/js/Pub.js"
import axios from "axios";

function UpLoad() {
    console.log(form, form.value)
    axios.post("http://127.0.0.1:7234/api/update", {
        prefix: form.value.prefix,
        proxy: form.value.proxy,
        cookie: form.value.cookie,
        r_18: form.value.r_18,
        downloadposition: form.value.downloadposition,
        likelimit: Number(form.value.likelimit),
        retry429: form.value.retry429,
        downloadinterval: form.value.downloadinterval,
        retryinterval: form.value.retryinterval,
        differauthor: form.value.differauthor,
        expired_time: form.value.expired_time,
    }, {
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(res => {
        form.value.prefix = res.data.setting.prefix
        form.value.proxy = res.data.setting.proxy
        form.value.cookie = res.data.setting.cookie.toString()
        form.value.r_18 = res.data.setting.r_18
        form.value.downloadposition = res.data.setting.downloadposition
        form.value.likelimit = res.data.setting.likelimit
        form.value.retry429 = res.data.setting.retry429
        form.value.downloadinterval = res.data.setting.downloadinterval
        form.value.retryinterval = res.data.setting.retryinterval
        form.value.differauthor = res.data.setting.differauthor
        form.value.expired_time = res.data.setting.expired_time

    }).catch(error => {

    })
}
</script>

<style lang="less" scoped>
@import "../assets/style/variable.less";
</style>
