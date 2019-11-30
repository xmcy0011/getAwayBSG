# 逃离北上广


[![Home](https://img.shields.io/badge/link-项目主页-brightgreen.svg)](https://jinnrry.github.io/getAwayBSG/)
[![Link](https://img.shields.io/badge/link-python实现-blue.svg)](https://github.com/jinnrry/getAwayBSG/tree/python)
[![Downloads](https://img.shields.io/github/downloads/jinnrry/getAwayBSG/total)](https://img.shields.io/github/downloads/jinnrry/getAwayBSG/total)
[![forks](https://img.shields.io/github/forks/jinnrry/getAwayBSG?style=flat)](https://img.shields.io/github/forks/jinnrry/getAwayBSG?style=flat)
[![starts](https://img.shields.io/github/stars/jinnrry/getAwayBSG)](https://img.shields.io/github/stars/jinnrry/getAwayBSG)
[![license](https://img.shields.io/github/license/jinnrry/getAwayBSG)](https://img.shields.io/github/license/jinnrry/getAwayBSG)
[![issues](https://img.shields.io/github/issues/jinnrry/getAwayBSG)](https://img.shields.io/github/issues/jinnrry/getAwayBSG)
[![version](https://img.shields.io/github/release/jinnrry/getAwayBSG)](https://img.shields.io/github/release/jinnrry/getAwayBSG)



> **注意！**\
> 1.本项目仅供学习研究，禁止用于任何商业项目\
> 2.运行的时候为被爬方考虑下！尽量不要爬全站。请在配置文件中设置你需要的城市爬取即可！\
> 3.[项目主页](https://jinnrry.github.io/getAwayBSG/)里面有现成数据，不需要你自己动手运行爬虫 

## 进度

PS：原项目地址https://jinnrry.github.io/getAwayBSG/  
> 如果你是一个正准备逃离北上广等一线城市却又找不到去处的IT人士，或许这个项目能给你点建议。  
> 通过爬虫，抓取了链接、智联的工作，租房，二手房一系列数据，为你提供各城市的宏观分析数据

- [x] 打印logo
- [x] 链家二手房抓取
    - [x] 多线程
    - [x] mongodb
    - [x] 抓取房屋概要(单线程)/详情(多线程)
    - [x] 按区抓取城市
    - [x] 记录城市名
- [ ] 链家租房抓取
- [ ] 智联招聘



## 安装编译

```bash
cat /etc/redhat-release # 作者Linux环境 x64
# CentOS Linux release 7.6.1810 (Core)
```

### mongodb
1.创建仓库文件
```bash
vi /etc/yum.repos.d/mongodb-org-3.4.repo

#  然后复制下面配置,保存退出
[mongodb-org-3.4]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/$releasever/mongodb-org/3.4/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-3.4.asc
```

2.安装配置
```bash
yum install -y mongodb-org # 安装

vim /etc/mongod.conf

# 更改数据存储目录，默认/var/log/mongodb
dbPath: /data/db
```

3.启动
```bash
systemctl restart mongod # 启动mongodb
systemctl enable mongod # 开启启动，可选
```

### golang
```bash
yum install golang # 安装go

# 配置go root和 go path
vim /etc/profile
#golang
export GOROOT=/usr/lib/golang
export GOPATH=/data/go  # 这个是存放代码位置
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

source /etc/profile # 生效
go env # 确认
```

### 编译
```bash
cd /data/go # go path 的路径
mkdir src && cd src
mkdir github.com && cd github.com
git clone https://github.com/xmcy0011/getAwayBSG.git 
cd getAwayBSG
go build # 编译
cp getAwayBSG bin/ # 把程序拷贝到bin目录
./getAwayBSG # 运行，运行前请先确认config.yaml里面的配置信息
```

## 运行

1.配置要抓取的城市 config.yaml  
2.运行
```bash
./getAwayBSG # 第一次可选择1，抓取概要列表后，会自动抓取详情，如果中途中断，下次运行选2即可
```

```shell script
  _______  _______ .___________.     ___      ____    __    ____      ___      ____    ____ .______        _______.  _______ 
 /  _____||   ____||           |    /   \     \   \  /  \  /   /     /   \     \   \  /   / |   _  \      /       | /  _____|
|  |  __  |  |__   `---|  |----`   /  ^  \     \   \/    \/   /     /  ^  \     \   \/   /  |  |_)  |    |   (----`|  |  __  
|  | |_ | |   __|      |  |       /  /_\  \     \            /     /  /_\  \     \_    _/   |   _  <      \   \    |  | |_ | 
|  |__| | |  |____     |  |      /  _____  \     \    /\    /     /  _____  \      |  |     |  |_)  | .----)   |   |  |__| | 
 \______| |_______|    |__|     /__/     \__\     \__/  \__/     /__/     \__\     |__|     |______/  |_______/     \______| 
                                                                                                                             
2019-11-30 17:05:21.472 [INFO]  [getAwayBSG/main.go:45] dbAddress=mongodb://106.14.172.35:27017,dbName=crawler
每次抓取都是全量的!
Usage of ./getAwayBSG:
  -config string
        设置配置文件 (default "./config.yaml")
  -help
        显示帮助
  -lianjia_ershou
        抓取链家二手房数据
  -lianjia_zufang
        抓取链家租房数据
  -zhilian
        抓取智联招聘数据
2019-11-30 17:05:21.473 [INFO]  [getAwayBSG/main.go:72] 请选择任务
2019-11-30 17:05:21.473 [INFO]  [getAwayBSG/main.go:73] 1.爬取链家二手房（全量）
2019-11-30 17:05:21.473 [INFO]  [getAwayBSG/main.go:74] 2.爬取链家二手房（详情）
```

## Contact

email: xmcy0011@sina.com