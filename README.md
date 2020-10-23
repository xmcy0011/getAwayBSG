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

> 如果你是一个正准备逃离北上广等一线城市却又找不到去处的 IT 人士，或许这个项目能给你点建议。  
> 通过爬虫，抓取了链接、智联的工作，租房，二手房一系列数据，为你提供各城市的宏观分析数据

- [x] 打印 logo
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

## 数据分析

下载 mongodb compass:https://www.mongodb.com/download-center/compass?jmp=docs

用法简介：  
1.Documents 选项（点击 OPTIONS，显示更多选项）

```mongodb
FILTER：查询条件
{City:"cs"} # 显示所有长沙的二手房

PROJECT：选择显示那些字段，1显示，0排除
{City:1,Title:1} # 只显示城市和标题

SORT：排序，-1,1代表倒序或顺序
{Size:-1} # 大小从大到小排序
```

2.Aggregations（举例，聚合查询 sh-上海户型，并倒序展示聚合结果）

```mongodb
1.点击 ADD STAGE，依次增加下面的命令：
$project:  # 要显示那些字段
{
  City:1,
  BaseAttr:1
}

$unwind:  # 展开数组为扁平化的，便于聚合查询
{
  path: "$BaseAttr", # 这是一个数组，unwind命令可以把它展开，便于聚合0位置表示户型
  includeArrayIndex: 'arrayIndex' # 给数组索引起个别名
}

$match:   # 过滤
{
  $and:[
    {arrayIndex:0}, # 上面的
    {City:"sh"},    # 上海
  ]
}

$group:  # 聚合分组统计
{
  _id: "$BaseAttr",
  count: {
    $sum: 1        # 新建count字段，$sum表示出现一次$BaseAttr.0 就 + 1
  }
}

2.点击Save下拉 -> Create a view，此时显示聚合结果
3.SORT里填入：{count:-1}，从大到小，显示最多的二手房户型
```

常用命令
- $gt:大于
- $lt:小于
- $gte:大于或等于
- $lte:小于或等于
- 模糊查询（包含）：{Title:{$regex:"公寓"},City:"sh"}

推荐资源：

1. 【最有用】官方帮助手册 https://docs.mongodb.com/manual/reference/operator/aggregation-pipeline/
2. 菜鸟教程 https://www.runoob.com/mongodb/mongodb-aggregate.html
3. 《MongoDB 权威指南第 2 版》.pdf，详细了解每个参数含义

*******
**PS：下面的语句，直接copy进入./mongo 里面使用**
*******

### 最受欢迎户型

```mongodb
db.lianjia.aggregate([
{
      $project:{
            City:1,
            BaseAttr:1
      }
},
{
      $unwind:{
            path: "$BaseAttr",
            includeArrayIndex: 'arrayIndex'
      }
},
{
      $match:{
            $and:[
                  {arrayIndex:0},
                  {City:"sh"}
            ]
      }
},
{
      $group:{
            _id: "$BaseAttr",
            count: {
                  $sum: 1
            }
      }
},
{
	$sort: { count: -1 },
      $limit: 30
}
])
```

### 各区域均价

1.按区域（0:"上海房产网" 1:"上海" **2:"浦东"** 3:"张江"）
```mongodb
db.lianjia.aggregate([
{
      $project:{
            City:1,
            AreaName:1,
            UnitPrice:1
      }
},
{
      $unwind:{
            path: "$AreaName",
            includeArrayIndex: 'arrayIndex'
      }
},
{
      $match:{
            $and:[
                  {arrayIndex:2},
                  {City:"sh"}
            ]
      }
},
{
      $group:{
            _id: "$AreaName",
            count: {$sum: 1},
            avgUintPrice:{$avg:"$UnitPrice"}
      }
},
{
      $project:{
            avgUintPriceCeil:{$ceil:"$avgUintPrice"},
            AreaName:1,
            count:1
      }
},
{
	$sort: { avgUintPriceCeil: -1 }
}
])
```

2.按区域（0:"上海房产网" 1:"上海" 2:"浦东" **3:"张江"**）
```mongodb
db.lianjia.aggregate([
{
      $project:{
            City:1,
            AreaName:1,
            UnitPrice:1
      }
},
{
      $unwind:{
            path: "$AreaName",
            includeArrayIndex: 'arrayIndex'
      }
},
{
      $match:{
            $and:[
                  {arrayIndex:3},
                  {City:"sh"}
            ]
      }
},
{
      $group:{
            _id: "$AreaName",
            count: {$sum: 1},
            avgUintPrice:{$avg:"$UnitPrice"}
      }
},
{
      $project:{
            avgUintPriceCeil:{$ceil:"$avgUintPrice"},
            AreaName:1,
            count:1
      }
},
{
	$sort: { avgUintPriceCeil: 1 }
}
])
```

### 地理空间分析

#### 火车站1公里附近房源
高德坐标拾取：https://lbs.amap.com/console/show/picker  
mongodb docs：https://docs.mongodb.com/manual/reference/operator/query/nearSphere/    

1.建空间索引 
```mongo
# mongodb 提供的地图索引有两种，分别是 2d 和 2dsphere
# 2d 索引通过二维平面记录点坐标，支持在平面几何中计算距离，而 2dsphere 则支持在球面上进行距离的计算，并且支持 mongodb 的所有地理空间查询方法
db.lianjia.ensureIndex({"Location":"2dsphere"})
```

芙蓉广场(为例)： 
```go
[112.984947,28.195951]
```

MongoDB Compass:
```mongodb
{
  'Location': {
    $nearSphere: {
      $geometry: {type: 'Point', coordinates: [121.493443,31.161054]},
      $maxDistance: 1000
    }
  },'TotalPrice':{'$lt':2000000}
}
```

Shell:
```mongodb
db.lianjia_2.find({
    Location: {
      $nearSphere: {
        $geometry: {
          type: "Point",
          coordinates: [112.984947, 28.195951]
        },
        $maxDistance: 2000
      }
    }
  }).limit(2);
```

## 其他

1.mongodb导出excel
```shell
mongoexport -d crawler -c lianjia_sh_20200911 -f ListCrawlTime,ListHouseType,ListHouseSize,ListHouseWhat,City,ListHouseFloor,Tag,Title,UnitPrice,Link,ListVillageName,ListAreaName,ListHouseBorn,TotalPrice,ListHouseOrientations,ListHouseDecorate,AreaName,BaseAttr,BeOnlineTime,CompletedInfo,DecorateInfo,DetailCrawlTime,DirectionInfo,FloorInfo,HouseRecordLJ,RoomInfo,Size,TransactionAttr,VillageName,FormattedAddress,Location --type=csv -o ./lianjia_sh_20200911.csv
```

2.数据分析软件推荐：Tableau

3.下载
```shell
scp -P 51219 root@106.14.172.35:/data/getAwayBSG/bin/data/lianjia_sh_20200729.csv .
```

## Contact

email: xmcy0011@sina.com
