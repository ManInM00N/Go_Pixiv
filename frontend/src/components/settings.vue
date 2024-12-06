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
                <el-input type="text" oninput="value=value.replace(/[^\d.]/g,'')" v-model="form.likelimit" />
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
// const props = defineProps({
//
//   form:{
//     type:DAO.Settings,
//   }
// })

function UpLoad() {
    console.log(form, form)
    axios.post("http://127.0.0.1:7234/api/update", {
        prefix: form.prefix,
        proxy: form.proxy,
        cookie: form.cookie,
        r_18: form.r_18,
        downloadposition: form.downloadposition,
        likelimit: form.likelimit,
        retry429: form.retry429,
        downloadinterval: form.downloadinterval,
        retryinterval: form.retryinterval,
        differauthor: form.differauthor
    }, {

    }).then(res => {

    }).catch(error => {

    })
}
</script>

<style lang="less" scoped>
@import "../assets/style/variable.less";
</style>
