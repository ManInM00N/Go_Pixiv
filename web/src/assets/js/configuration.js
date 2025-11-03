import axios from "axios";
import { ref } from "vue"
import {CheckLogin, OpenFileFolder} from "../../../bindings/main/internal/pixivlib/ctl.js";
import {ElNotification} from "element-plus";
export const defaultConf = {
    prefix: 'http://127.0.0.1:7234',
    proxy: '',
    useproxy: false,
    imageEngine:{
        sauceNaoConf:{
            api_key: "",
            numbers: 10,
            testmode: false,  // 测试模式（不消耗配额）
            db: 999,          // 搜索数据库：999=全部, 5=Pixiv, 9=Danbooru等
            hide: 0,          // 隐藏成人内容：0=不隐藏, 1=隐藏, 2=部分隐藏
            dedupe: 2,        // 去重级别：0=关闭, 1=同站, 2=全局
        }
    },
    pixivConf:{
        cookie: '',
        r_18: false,
        downloadposition: 'Download',
        likelimit: 0,
        retry429: 2000,
        downloadinterval: 500,
        retryinterval: 1000,
        differauthor: false,
        expired_time: 7,
        logined: false,
    },

}
let prePost = {}
export let form = ref(Object.assign({}, defaultConf))
axios.get("http://127.0.0.1:7234/api/getsetting").then(res => {
    console.log(form.value)

    form.value.prefix = res.data.setting.prefix
    form.value.proxy = res.data.setting.proxy
    form.value.useproxy = res.data.setting.useproxy

    form.value.pixivConf.cookie = res.data.setting.pixivConf.cookie.toString()
    form.value.pixivConf.r_18 = res.data.setting.pixivConf.r_18
    form.value.pixivConf.downloadposition = res.data.setting.pixivConf.downloadposition
    form.value.pixivConf.likelimit = Number(res.data.setting.pixivConf.likelimit)
    form.value.pixivConf.retry429 = res.data.setting.pixivConf.retry429
    form.value.pixivConf.downloadinterval = res.data.setting.pixivConf.downloadinterval
    form.value.pixivConf.retryinterval = res.data.setting.pixivConf.retryinterval
    form.value.pixivConf.differauthor = res.data.setting.pixivConf.differauthor
    form.value.pixivConf.expired_time = res.data.setting.pixivConf.expired_time

    form.value.imageEngine.sauceNaoConf.numbers = res.data.setting.imageEngine.sauceNaoConf.numbers
    form.value.imageEngine.sauceNaoConf.api_key = res.data.setting.imageEngine.sauceNaoConf.api_key

    prePost = Object.assign({}, form.value)
    prePost.pixivConf = Object.assign({}, prePost.pixivConf)
    console.log(res.data)
    console.log(form.value)
    if (form.value.cookie != "") {
        CheckLogin()
    }
}).catch(error => {
    console.log(error)
})

const reconnectFields = ['prefix', 'proxy', 'pixivConf', 'useproxy']


export async function ChooseFolder(){
    form.value.pixivConf.downloadposition =  await OpenFileFolder()

}

export async function updateSettings(){
    let needRecon = false
    await axios.post("http://127.0.0.1:7234/api/update", {
        prefix: form.value.prefix,
        proxy: form.value.proxy,
        useproxy: Boolean(form.value.useproxy),
        pixivConf:{
            cookie: form.value.pixivConf.cookie,
            r_18: Boolean(form.value.pixivConf.r_18),
            downloadposition: form.value.pixivConf.downloadposition,
            likelimit: Number(form.value.pixivConf.likelimit),
            retry429: Number(form.value.pixivConf.retry429),
            downloadinterval: Number(form.value.pixivConf.downloadinterval),
            retryinterval: Number(form.value.pixivConf.retryinterval),
            differauthor: Boolean(form.value.pixivConf.differauthor),
            expired_time: Number(form.value.pixivConf.expired_time),
        },
        imageEngine:{
            sauceNaoConf:{
                api_key: form.value.imageEngine.sauceNaoConf.api_key,
                numbers: form.value.imageEngine.sauceNaoConf.numbers,
            }
        },

    }, {
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(res => {
        for (const field of reconnectFields) {
            let oldValue = prePost[field]
            console.log(field)
            let newValue = res.data.setting[field]
            if (field === 'pixivConf') {
                oldValue = prePost.pixivConf.cookie
                newValue = res.data.setting.pixivConf.cookie.toString()
                console.log(oldValue,newValue)
            }
            if (oldValue !== newValue) {
                needRecon = true
                break
            }
        }
        console.log(needRecon)

        let nxt = !needRecon && form.value.pixivConf.logined

        form.value.prefix = res.data.setting.prefix
        form.value.proxy = res.data.setting.proxy
        form.value.useproxy = res.data.setting.useproxy

        form.value.pixivConf = {
            cookie: res.data.setting.pixivConf.cookie.toString(),
            r_18: res.data.setting.pixivConf.r_18,
            downloadposition: res.data.setting.pixivConf.downloadposition,
            likelimit: Number(res.data.setting.pixivConf.likelimit),
            retry429: res.data.setting.pixivConf.retry429,
            downloadinterval: res.data.setting.pixivConf.downloadinterval,
            retryinterval: res.data.setting.pixivConf.retryinterval,
            differauthor: res.data.setting.pixivConf.differauthor,
            expired_time: res.data.setting.pixivConf.expired_time,
            logined: nxt,
        }


        form.value.imageEngine = {
            sauceNaoConf:{
                api_key: res.data.setting.imageEngine.sauceNaoConf.api_key,
                numbers: res.data.setting.imageEngine.sauceNaoConf.numbers,
            }
        }

        prePost = Object.assign({}, form.value)
        prePost.imageEngine = Object.assign({}, prePost.imageEngine)
        prePost.imageEngine.sauceNaoConf = Object.assign({}, prePost.imageEngine.sauceNaoConf)
        prePost.pixivConf = Object.assign({}, prePost.pixivConf)
        console.log(res)
    }).catch(error => {
        console.log(error)
        ElNotification({
            type: "error",
            title: "保存失败",
            message: "error",
            position: 'bottom-right',
            duration: 5000,
        })
    }).finally(() => {

        // CheckLogin()
    })
    return needRecon
}


export const waterFallConf ={
    breakpoints: {
        1750:{
            rowPerView: 5,
        },
        1400:{
            rowPerView: 4,
        },
        1050: {
            rowPerView: 3,
        },
        700: {
            rowPerView: 2,
        },
        400: {
            rowPerView: 1,
        }
    },
    width : 300,
    gutter : 20,
}