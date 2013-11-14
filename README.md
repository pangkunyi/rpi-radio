rpi-radio
=========

基于百度PCS+Raspberry Pi来实现的音乐电台

###依赖库
```
github.com/pangkunyi/baidu-pcs
github.com/gorilla/mux
```

###安装
```
1, git clone git@github.com:pangkunyi/rpi-radio.git
2, make dep
3, make install
4, make run
```

###配置

1， 在用户目录下创建并编辑.baidu_pcs.cfg.json文件，内容如下：
```
{"client_id":"9xdPohPviuXXXXXXXXXXX", "secret_id":"drahBeGpC9XXXXXXXXXXXXX", "open_dir":"/apps/kunyi"}
```
client_id和secret_id分别为你在百度开放平台创建应用时的API KEY和SECRET KEY, open_dir为应用申请PCS API时允许访问的目录

2， 获取授权码
```
后续补上
```

3, 在用户目录下创建并编辑.baidu_pcs.auth_data.json文件，内容如下：
```
{"access_token":"3.6ea190686a21cb61XXXXXXXXXXXX"}
```
access_token为步骤2获取的授权码

###运行
```
make run
```

###访问
```
http://127.0.0.1:8808/
```
由于偷懒，8808端口硬编码了
