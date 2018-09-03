myHardSeed

A go version of hardseed(C++) 

项目源地址：https://github.com/yangyangwithgnu/hardseed

使用该程序请自备小飞机

迅雷（极速版）依旧可以使用

使用说明

命令参数解析

- -save-path:文件保存路径
  - 用法 -save-path = "C:\Download\",注意，此处路径结尾一定需要 “\”,linux下则需要"/"
  - 默认参数  windows=“D:\download\”  macOS= "/user/" linux= "~/"
- -current-task:并行的下载线程数
  - 用法 -current-task = 16,注意，若想体验告诉下载的乐趣，cpu的核心数最好超过两个，手动狗头 
  - 默认参数 16
- -core-number:程序所使用的cpu核心数
  - 用法 -core-number= 2
  - 默认参数 cpu所有核心数
- -av-class: 选择的种子分区
  - 用法 -av-class = "caoliu_asia_non_mosaicked_original"
  - 默认参数：“aicheng_asia_mosaicked”
  - 参数详表：
    - caoliu_west_original
    - caoliu_cartoon_original
    - caoliu_asia_mosaicked_original
    - caoliu_asia_non_mosaicked_original
    - caoliu_selfie
    - aicheng_west
    - aicheng_cartoon
    - aicheng_asia_mosaicked
    - aicheng_asia_non_mosaicked
    - aicheng_selfie
- -topic-range:选取当前主题下的目录文件范围
  - 用法 -topic-range = "100 1024"
  - 默认参数 “0 1024”
- -proxy：代理设置
  - 用法 -proxy="socks5://127.0.0.1:1080" 注意，小飞机一般都是此代理出口
  - 默认参数“socks5://127.0.0.1:1080”
- -like：喜欢的主题
  - 用法 -like = “name1 name2 name3”
  - 默认 “”
- -hate: 不喜欢的主题
  - 用法 -hate = "name1 name2 name3"
  - 默认“”



使用示例

- windows用户：
  - 找到myHardSeed.exe所在窗口
  - alt+d 然后输入cmd并回车运行（如果不行请用管理员权限打开cmd并运行myHardSeed.exe）
  - 输入命令，示例：myHardSeed.exe -save-path = "D:\Download\"   -like = "ABP snis"  -topic-range = "1 500"

TODO

1. 处理caoliu的快速翻页问题
2. 处理多线程报错问题: Unsolicited response received on idle HTTP channel starting with "\n"; err=<nil>
3. 添加-1为全主题检索
