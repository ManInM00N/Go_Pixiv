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
    expired_time: 7,
})
axios.get("http://127.0.0.1:7234/api/getsetting").then(res => {
    console.log(form.value)
    form.value.prefix = res.data.setting.prefix
    form.value.proxy = res.data.setting.proxy
    form.value.cookie = res.data.setting.cookie.toString()
    form.value.r_18 = res.data.setting.r_18
    form.value.downloadposition = res.data.setting.downloadposition
    form.value.likelimit = Number(res.data.setting.likelimit)
    form.value.retry429 = res.data.setting.retry429
    form.value.downloadinterval = res.data.setting.downloadinterval
    form.value.retryinterval = res.data.setting.retryinterval
    form.value.differauthor = res.data.setting.differauthor
    console.log(res.data)
    console.log(form.value)
    if (form.value.cookie != "") {
        CheckLogin()
    }
}).catch(error => {
    console.log(error)
})
export let ws = WebSocket;
const startWebSocket = () => {
    ws.value = new WebSocket("ws://127.0.0.1:7234/api/ws");

    ws.value.onopen = () => {
        console.log('WebSocket connected');
    };

    ws.value.onmessage = (event) => {
        // res.value = event.data;
        handleMessage(JSON.parse(event.data));
    };
    ws.value.onclose = () => {
        ElMessage.error("远程主机关闭")
        console.log('WebSocket closed');
    };

    ws.value.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
};
startWebSocket()
