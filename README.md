# 使用须知 Read Me First

<p align="center">

<img src="./build/appicon.png">
<br>
<a href="https://github.com/ManInM00N/Go_Pixiv/releases"> <img src="https://img.shields.io/github/v/release/ManInM00N/Go_Pixiv" alt=""/> </a>
<a href="https://golang.google.cn/dl/"> <img src="https://img.shields.io/github/go-mod/go-version/ManInM00N/Go_Pixiv"> </a>
<a href="https://wails.io/zh-Hans/docs/gettingstarted/installation"> <img src="https://img.shields.io/badge/wails-v3.0.0--alpha8.3-red"> </a>
<a href="https://github.com/ManInM00N/Go_Pixiv/blob/master/LICENSE"><img src="https://img.shields.io/github/license/ManInM00N/Go_Pixiv"> </a>

</p>

- [x] 升级wails框架版本 :tada:
- [x] 优化预览缓存机制,设置图片代理 :wave:
- [ ] 优化下载线程池 :confused:
- [ ] 以图搜图 :alien:
- [x] 托盘化 :zap:
- [x] 小说下载功能 :monocle_face:
- GIF下载暂时无法实现,目前找不到一个比较好的办法将返回的图片制作成GIF :clown_face:

# ！！！不提供代理

代理问题需要自己解决
关注页面需要有自己的pixiv账号(cookie),若个人行为导致账号封禁概不负责<br>
保证代理无误的情况下，不添加cookie也可以下载，但是不能下载和观看r18和部分图片<br>
应用内置网页,点击图片的名字可以打开对应网页,所以rank和follow两个页面其实没必要，但是做都做了索性稍微写了写<br>
下载速度取决于你的代理，理论上参数可以调很快但是会报429，而且太快会被封443端口，所以不建议下太快，默认设置间隔时间1000ms

## 制作参考

<https://github.com/daydreamer-json/pixiv-ajax-api-docs/tree/main><br>

## 配置设定

### 在settings.yml中

**proxy**:代理ip后面的端口，可以从你vpn的配置中得到，这个不会配的话我无能为力

**cookie**:打开登录后的pixiv网页，在电脑网页按F12，从应用程序一栏中Cookie的PHPSESSID的值

**r-18**:true启用，false禁用   懂得都懂，不启用是无法看到r-18的图片的

**minlikelimit**:下载图片的点赞数限制 小于的不下载

**downloadposition**:图片储存位置，如果目标位置没有文件夹则会改成此目录下的Download文件夹(自动创建)

**retry429**: 429状态码等待时间

**downloadinterval**: 下载间隔时间

**retryinterval**: 请求重试间隔
