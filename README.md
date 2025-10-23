# 使用须知 Read Me First

<p align="center">

<img style="width :256px" src="./build/appicon.png">
<br>
<a href="https://github.com/ManInM00N/Go_Pixiv/releases"> <img src="https://img.shields.io/github/v/release/ManInM00N/Go_Pixiv" alt=""/> </a>
<a href="https://golang.google.cn/dl/"> <img src="https://img.shields.io/github/go-mod/go-version/ManInM00N/Go_Pixiv"> </a>
<a href="https://wails.io/zh-Hans/docs/gettingstarted/installation"> <img src="https://img.shields.io/badge/wails-v3.0.0--alpha36-red"> </a>
<a href="https://github.com/ManInM00N/Go_Pixiv/blob/master/LICENSE"><img src="https://img.shields.io/github/license/ManInM00N/Go_Pixiv"> </a>

</p>

- [x] 升级wails框架版本 :tada:
- [x] 优化预览缓存机制,设置图片代理 :wave:
  - [x] 图片详细展示 
- [x] 优化下载线程池 :confused:
  - [ ] 中断任务
- [ ] 以图搜图 :alien:
- [x] 托盘化 :zap:
  - [ ] 桌面通知 
- [x] 小说下载功能 :monocle_face:
  - [x] 小说查看功能 
  - [x] 小说系列下载功能 :wheelchair:
- [x] GIF下载功能 :tada:
  - [ ] 优化GIF Encoder

## ！！！不提供代理

代理问题需要自己解决
关注页面需要有自己的pixiv账号(cookie),若个人行为导致账号封禁概不负责<br>
保证代理无误的情况下，不添加cookie也可以下载，但是不能下载和观看r18和部分图片<br>
下载速度取决于你的代理，理论上参数可以调很快但是会报429，而且太快会被封443端口，所以不建议下太快，默认设置间隔时间1000ms

## 制作接口参考

<https://github.com/daydreamer-json/pixiv-ajax-api-docs/tree/main><br>

## 安装和调试
You should Install Go >= 1.24.5 and Node.js >= v20

wails3 安装
```sh
    go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-alpha.36
```

调试
```sh
    wails3 task dev
```
for more details see Taskfile.yml


## 🐛 问题反馈

如果遇到问题，请提交 Issue 到 [GitHub](https://github.com/ManInM00N/Go_Pixiv/issues)。
