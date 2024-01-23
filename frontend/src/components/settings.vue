<template>
  <el-main>
    <el-form
        label-width="100px"
        :model="form"
        style="max-width: 460px"
        class="Background_Opacity"
    >
      <el-form-item
      label="Localhost:"
      >
        <el-input
            v-model="$props.form.prefix"
            disabled="true"
        />
      </el-form-item>
      <el-form-item
          label="代理端口"
      >
        <el-input
          v-model="$props.form.proxy"
        >
        </el-input>
      </el-form-item>
      <el-form-item
          label="Cookie"
      >
        <el-input
            v-model="$props.form.cookie"
        >
        </el-input>
      </el-form-item>
      <el-form-item
          label="下载位置"
      >
        <el-input v-model="$props.form.downloadposition"/>
      </el-form-item>
      <el-form-item>

      </el-form-item>
      <el-form-item
          label="请求自省时间(ms)"
      >
        <el-input
            type="number"

            v-model="$props.form.retry429"
        />
      </el-form-item>
      <el-form-item
          label="请求失败重试间隔(ms)"
      >
        <el-input
            type="number"

            v-model="$props.form.retryinterval"
        />
      </el-form-item>
      <el-form-item
          label="下载间隔(ms)"
      >
        <el-input
            type="number"
            v-model="$props.form.downloadinterval"
        />
      </el-form-item>
      <el-form-item
          label="作品收藏数限制"
      >
        <el-input
            type="number"
            v-model="$props.form.minlikelimit"
        />
      </el-form-item>
      <el-form-item>
        <el-checkbox
            label="是否启用R-18"
            v-model="$props.form['r-18']"
        />
        <el-checkbox
            label="下载图片对作者分类"
            v-model="$props.form.differauthor"
        />
      </el-form-item>
      <el-form-item>
        <el-button
          @click="UpLoad"
        >
          Update
        </el-button>
      </el-form-item>
    </el-form>
  </el-main>
</template>

<script lang="ts" setup>
import {EventsOn} from "../../wailsjs/runtime/runtime.js";
import {defineComponent, onMounted, ref,reactive} from "vue";
import {UpdateSetting,GetSetting} from "../../wailsjs/go/main/App.js";
import {DAO} from "../../wailsjs/go/models.ts";
import emitter from "../assets/js/Pub.js"
const props = defineProps({
  form:{
    type:DAO.Settings,
  }
})
function UpLoad(){
  console.log(props.form)
  UpdateSetting(props.form)
  emitter.emit("Relogin")
}
</script>

<style lang="less" scoped>
@import "../assets/style/variable.less";
</style>