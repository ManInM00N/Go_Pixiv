<h1>使用须知 Read Me First!!!</h1>
<h1>
！！！不提供代理！！！
<br>
梯子代理问题需要自己解决
<br>
需要有自己的pixiv账号,若个人行为导致账号封禁概不负责<br>
保证代理无误的情况下，不添加cookie也可以下载<br>
<br>
下载速度取决于你的代理，理论上参数可以调很快但是会给429，而且太快会被封443端口，所以不建议下太快，默认设置间隔时间1000ms
</h1>
<h2>
目前没有设计gui，暂时使用fyne做简陋的读入，作者现在大二前端写的依托，所以等寒假摸完了vue，再换成wails在做成包发布<br>
<br>
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
r-18:true启用，false禁用   懂得都懂<br>
minlikelimit:下载图片的点赞数限制 小于的不下载<br>
downloadposition:图片储存位置，如果目标位置没有文件夹则会改成此目录下的Download文件夹(自动创建)<br>
retry429: 429状态码等待时间<br>
downloadinterval: 下载间隔时间<br>
retryinterval: 请求重试间隔<br>
</p>