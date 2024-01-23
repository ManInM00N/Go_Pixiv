<h1>使用须知 Read Me First!!!</h1>
<br>
<p align="center">

<img src="https://github.com/ManInM00N/Go_Pixiv/blob/master/build/appicon.png">

![License](https://img.shields.io/badge/LICENSE-GPL3.0-green)
![Static Badge](https://img.shields.io/badge/Wails-2.7.1-%2332CD99)
![Static Badge](https://img.shields.io/badge/Go-%3E%3D1.21.5-%23007FFF)
</p>


<h1>
！！！不提供代理！！！<br>
梯子代理问题需要自己解决<br>
关注页面需要有自己的pixiv账号(cookie),若个人行为导致账号封禁概不负责<br>
保证代理无误的情况下，不添加cookie也可以下载，但是不能下载和观看r18图片<br>
应用内置网页,点击图片的名字可以打开对应网页,所以rank和follow两个页面其实没必要，但是做都做了索性稍微写了写<br>
下载速度取决于你的代理，理论上参数可以调很快但是会报429，而且太快会被封443端口，所以不建议下太快，默认设置间隔时间1000ms
</h1>
<h2>
制作参考:https://github.com/daydreamer-json/pixiv-ajax-api-docs/tree/main<br>
</h2>
<h2>配置设定</h2>
<h3>在settings.yml中</h3>
<p>
proxy:你本地梯子的代理ip后面的端口，可以从你梯子的配置中得到，这个不会配的话我无能为力,拿v2ray的配置方法举例，端口就是http后面的数字：<br>
<img src="https://github.com/ManInM00N/go-pixiv/blob/master/assets/proxy.png"><br>
cookie:打开登录后的pixiv网页，在电脑网页按F12，从应用程序一栏中Cookie的PHPSESSID的值<br>
<img src="https://github.com/ManInM00N/go-pixiv/blob/master/assets/cookie1.png"><br>
<img src="https://github.com/ManInM00N/go-pixiv/blob/master/assets/cookie2.png"><br>
<img src="https://github.com/ManInM00N/go-pixiv/blob/master/assets/cookie3.png"><br>
r-18:true启用，false禁用   懂得都懂，不启用是无法看到r-18的图片的<br>
minlikelimit:下载图片的点赞数限制 小于的不下载<br>
downloadposition:图片储存位置，如果目标位置没有文件夹则会改成此目录下的Download文件夹(自动创建)<br>
retry429: 429状态码等待时间<br>
downloadinterval: 下载间隔时间<br>
retryinterval: 请求重试间隔<br>
同目录下会有一个cache文件夹，里面存放预览图，如果太多了手动删除即可，此时页面中的图片会丢失，刷新后重新下载预览图
</p>