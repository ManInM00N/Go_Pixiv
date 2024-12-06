import axios from "axios";
import { ref } from "vue"
import { CheckLogin } from "../../../bindings/main/ctl.js";
export let form = ref({
    prefix: '',
    proxy: '',
    cookie: '',
    r_18: false,
    downloadposition: 'download',
    likelimit: 0,
    retry429: 2000,
    downloadinterval: 500,
    retryinterval: 1000,
    differauthor: false,

})
axios.get("http://127.0.0.1:7234/api/getsetting").then(res => {
    console.log(form.value)
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
    console.log(res.data)
    console.log(res.data.setting.cookie.toString(), res.data.setting.cookie)
    console.log(form.value)
    if (form.value.cookie != "") {
        CheckLogin()
    }
}).catch(error => {
    console.log(error)
})
