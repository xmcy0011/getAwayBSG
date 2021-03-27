<template>
  <div class="full">
    <div id="map-container"></div>
    <div class="mapStyle">
      <el-row>
        <el-select
          size="small"
          style="width: 100px"
          v-model="curMapStyle"
          placeholder="请选择"
          @change="_mapStyleSelectedChange"
        >
          <el-option-group
            v-for="group in options3"
            :key="group.label"
            :label="group.label"
          >
            <el-option
              v-for="item in group.options"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-option-group>
        </el-select>

        <el-button
          size="small"
          @click="_onClickRangingTool"
          style="margin-left: 5px"
          >测距</el-button
        >
        <el-button
          size="small"
          @click="_onClickTimeQueryTool"
          style="margin-left: 5px"
          >通勤</el-button
        >
      </el-row>
      <el-row style="background-color: white; margin-top: 5px">
        <div>
          <el-checkbox
            v-model="itAreaChecked"
            @change="_loadItArea(itAreaChecked)"
            >科技园</el-checkbox
          >
          <el-checkbox
            v-model="itCompanyChecked"
            @change="_loadItCompany(itCompanyChecked)"
            >IT公司</el-checkbox
          >
        </div>
      </el-row>

      <el-drawer
        title="通勤时间测算"
        v-model="drawer"
        :direction="rtl"
        :before-close="drawerClose"
        destroy-on-close
        style="overflow: scroll"
      >
        <el-row>
          <el-input
            v-model="input"
            placeholder="请输入出发地"
            style="width: 400px; margin-left: 10px"
          ></el-input>
          <el-button style="margin-left: 10px">确定</el-button>
        </el-row>

        <el-card
          class="box-card"
          style="margin-top: 20px"
          v-for="o in destination"
          body-style="padding:15px;"
          :key="o"
        >
          <template #header>
            <div class="card-header" style="height:15px;">
              <span>{{ o.name }}</span>
            </div>
          </template>
          <el-row style="margin-top:-20px;">
            <el-col :span="12">
              <p>公交</p>
              <strong>1小时1分钟</strong>
            </el-col>
            <el-col :span="12">
              <p>驾车</p>
              <strong>54分钟</strong>
            </el-col>
          </el-row>
        </el-card>
      </el-drawer>
    </div>
  </div>
</template>

<script>
// @ is an alias to /src
import AMap from "AMap";
//const shellEle = require("electron").shell; // electron shell
//const Store = require("electron-store");    // electron storage

export default {
  name: "home",
  data() {
    return {
      map: null,
      massMarks: null, // 海量点
      ruler: null, // 测距工具
      drawer: false, // 通勤时间工具
      destination: [
        {
          name: "哔哩哔哩",
          lng: 121.506414,
          lat: 31.309352,
        },
        {
          name: "拼多多",
          lng: 121.425928,
          lat: 31.219444,
        },
        {
          name: "万达股份",
          lng: 121.519372,
          lat: 31.077439,
        },
        {
          name: "莉莉丝",
          lng: 121.39253,
          lat: 31.170216,
        },
      ],
      infoWindow: null,
      infoWindowMarker: null,
      options3: [
        {
          label: "推荐",
          options: [
            {
              value: "normal",
              label: "标准",
            },
            {
              value: "dark",
              label: "幻影黑",
            },
          ],
        },
        {
          label: "自定义",
          options: [
            {
              value: "macaron",
              label: "马卡龙",
            },
            {
              value: "graffiti",
              label: "涂鸦",
            },
            {
              value: "fresh",
              label: "草色青",
            },
            {
              value: "blue",
              label: "靛青蓝",
            },
            {
              value: "whitesmoke",
              label: "远山黛",
            },
            {
              value: "darkblue",
              label: "极夜蓝",
            },
            {
              value: "light",
              label: "月光银",
            },
            {
              value: "grey",
              label: "雅士灰",
            },
          ],
        },
      ],
      curMapStyle: "normal",
      curHouseMarker: null, // 当前激活的房源
      itCompanyAreaList: [], // 上海科技园,
      itCompanyList: [], // 上海It公司
      itAreaChecked: true,
      itCompanyChecked: true,
    };
  },
  mounted() {
    this.map = new AMap.Map("map-container", {
      mapStyle: "amap://styles/normal",
      zoom: 15,
      center: [121.480331, 31.153403],
    });
    AMap.plugin("AMap.ToolBar", () => {
      // 异步加载插件
      let toolbar = new AMap.ToolBar({
        position: "RB",
      });
      this.map.addControl(toolbar);
    });
    AMap.plugin(["AMap.Scale"], () => {
      var scale = new AMap.Scale();
      this.map.addControl(scale);
    });
    //加载距离测量插件
    let _this = this;
    AMap.plugin(["AMap.RangingTool"], function () {
      _this.ruler = new AMap.RangingTool(_this.map);
    });

    document.onkeydown = function (event) {
      var e = event || window.event || arguments.callee.caller.arguments[0];
      if (e && e.keyCode == 37) {
        console.log("key 37 press");
      } else if (e && e.keyCode == 38) {
        console.log("key 38 press");
      }
    };

    setTimeout(() => {
      // 加载著名上海科技园
      this._loadItArea();
      this._loadItCompany();
    }, 1 * 1000);
  },
  methods: {
    // 使用默认浏览器打开
    openUrl(link) {
      //shellEle.openExternal(link);
    },
    /**
     * 添加一个覆盖物
     * @param {[list]}：海量点列表
     [{
        lnglat: [116.405285, 39.904989], //点标记位置
        title: 'beijing1',
        id:1,
        style: 0
     },{
       lnglat: [116.405285, 39.904989], //点标记位置
        title: 'beijing2',
        id:2,
        style: 0
     }]
     */
    addMassMarks(list) {
      let _this = this;
      // 创建样式对象
      let so = [
        {
          url: require("../assets/wujiaoxing_red.png"), // 图标地址
          size: new AMap.Size(11, 11), // 图标大小
          anchor: new AMap.Pixel(5, 5), // 图标显示位置偏移量，基准点为图标左上角
        },
        {
          url: require("../assets/wujiaoxing_black.png"), // 图标地址
          size: new AMap.Size(11, 11), // 图标大小
          anchor: new AMap.Pixel(5, 5), // 图标显示位置偏移量，基准点为图标左上角
        },
        {
          url: require("../assets/wujiaoxing_yellow.png"), // 图标地址
          size: new AMap.Size(11, 11), // 图标大小
          anchor: new AMap.Pixel(5, 5), // 图标显示位置偏移量，基准点为图标左上角
        },
      ];

      list.forEach((item) => {
        // 收藏，则标记
        let type = _this._checkHouseIsCollect(item.full.HouseRecordLJ);
        if (type == 1) {
          // like
          item.style = 2;
        } else if (type == 2) {
          // dislike
          item.style = 1;
        }
      });

      let massMarks = new AMap.MassMarks(list, {
        zIndex: 111, // 海量点图层叠加的顺序
        cursor: "pointer",
        zooms: [3, 19], // 在指定地图缩放级别范围内展示海量点图层
        style: so, // 设置样式对象
      });

      //let marker = new AMap.Marker({ content: " ", map: this.map });
      let lastClickTime = new Date().valueOf();

      // 创建 infoWindow 实例
      let infoWindow = new AMap.InfoWindow();

      // 将海量点添加至地图实例
      massMarks.setMap(this.map);
      //mouseover
      massMarks.on("click", function (e) {
        // ListHouseType: "2室1厅"
        // ListHouseSize: 73.24
        // ListHouseWhat: "板楼"
        // City: "sh"
        // ListHouseFloor: "高楼层(共18层)"
        // Tag: "["近地铁","VR房源","房本满五年"]"
        // Title: "星颂南北卧室 一手动迁税费低 婚房精装住的少 全留"
        // UnitPrice: 38914
        // Link: "https://sh.lianjia.com/ershoufang/107102079791.html"
        // ListVillageName: "曹路"
        // ListAreaName: "星颂家园"
        // ListHouseBorn: "2011年建"
        // TotalPrice: 2850000
        // ListHouseOrientations: "南"
        // ListHouseDecorate: "精装"
        // AreaName: "["上海房产网","上海","浦东","曹路","星颂家园","当前房源"]"
        // BaseAttr: "["房屋户型:2室1厅1厨1卫","所在楼层:高楼层 (共18层)","建筑面积:73.24㎡","户型结构:平层","套内面积:暂无数据","建筑类型:板楼","房屋朝向:南","建筑结构:钢混结构","装修情况:精装","梯户比例:两梯四户","配备电梯:有","产权年限:70年"]"
        // CompletedInfo: "2011年建/板楼"
        // HouseRecordLJ: 107102079791
        // TransactionAttr: "["交易权属:动迁安置房","上次交易:2012-10-22","房屋用途:普通住宅","房屋年限:满五年","产权所属:共有","抵押信息:无抵押","房本备件:已上传房本照片"]"
        // FormattedAddress: "上海市浦东新区星颂家园"
        let h = e.data.full;
        _this.curHouseMarker = e;

        // 信息窗体的内容
        let content = ["<div style='font-size:12px;'>"];
        content.push("标题：" + h.Title);
        content.push("编号：" + h.HouseRecordLJ + " " + h.ListHouseBorn);
        content.push(
          "房型：" +
            h.ListHouseType +
            ", " +
            h.ListHouseOrientations +
            ", " +
            h.ListHouseDecorate
        );
        content.push(
          "价格：" +
            h.ListHouseType +
            "元(单价)" +
            ", " +
            h.TotalPrice +
            "元(总价)"
        );
        content.push(
          "链接：<span id='houseLick' data='" +
            h.Link +
            "' style='cursor:hand;color:blue;'><u>" +
            h.Link +
            "</u></span>"
        );
        content.push("区域：" + h.ListVillageName + " " + h.ListAreaName);
        content.push("大小：" + h.ListHouseSize + " 平米," + h.ListHouseWhat);
        content.push("楼层：" + h.ListHouseFloor);
        content.push(
          "标签：" +
            h.Tag.substring(1, h.Tag.length - 2)
              .split(",")
              .join(" ")
              .replace(/"/g, "")
        );
        content.push(
          "地址1：<br/>" +
            h.AreaName.substring(1, h.AreaName.length - 2)
              .split(",")
              .join(" ")
              .replace(/"/g, "")
        );
        content.push("地址2：" + h.FormattedAddress);
        content.push(
          "详情：<br/>" +
            h.BaseAttr.substring(1, h.BaseAttr.length - 2)
              .split(",")
              .join(" ")
              .replace(/"/g, "")
        );
        content.push(
          "其他：<br/>" +
            h.TransactionAttr.substring(1, h.TransactionAttr.length - 2)
              .split(",")
              .join(" ")
              .replace(/"/g, "")
        );
        content.push(
          "<button id='btnCollect' houseId='" +
            h.HouseRecordLJ +
            "'>收藏</button><button id='btnDisCollect' houseId='" +
            h.HouseRecordLJ +
            "' style='margin-left:5px;'>不看</button>"
        );
        content.push("</div>");
        // 打开信息窗体
        infoWindow.setPosition(e.data.lnglat);
        infoWindow.setContent(content.join("<br/>"));
        infoWindow.open(_this.map);
        setTimeout(
          (link) => {
            let span = document.getElementById("houseLick");
            span.onclick = function () {
              console.log("span click" + this.getAttribute("data"));
              _this.openUrl(this.getAttribute("data"));
            };

            let btnLike = document.getElementById("btnCollect");
            if (btnLike != null) {
              btnLike.onclick = function () {
                console.log("like:" + this.getAttribute("houseId"));
                _this._saveOrUpdateHouse(this.getAttribute("houseId"), 1);
                _this.curHouseMarker.data.style = 2;
              };
            }
            let btnDislike = document.getElementById("btnDisCollect");
            if (btnDislike != null) {
              btnDislike.onclick = function () {
                console.log("dislike:" + this.getAttribute("houseId"));
                _this._saveOrUpdateHouse(this.getAttribute("houseId"), 2);
                _this.curHouseMarker.data.style = 1;
              };
            }
          },
          1000,
          e.data.Link
        );

        // black
        e.data.style = 1;
      });
      this.massMarks = massMarks;
    },
    // 清楚所有覆盖物
    clear() {
      this.map.clearMap();
      if (this.massMarks != null) {
        this.massMarks.clear();
      }
    },
    // 收藏 collectType:0,delete,1:like,2:dislike
    _saveOrUpdateHouse(hourseId, collectType) {
      // if (collectType == 1 || collectType == 2) {
      //   const store = new Store();
      //   store.set("house_" + hourseId, collectType);
      //   console.log(store.get("house_" + hourseId));
      // } else {
      //   store.delete("house_" + hourseId);
      // }
    },
    // 是否收藏
    _checkHouseIsCollect(hourseId) {
      //const store = new Store();
      //return store.get("house_" + hourseId);
      return false;
    },
    // 测距
    _onClickRangingTool() {
      console.log("_onClickRangingTool");
      if (this.ruler != null) {
        this.ruler.turnOn();
      }
    },
    // 通勤
    _onClickTimeQueryTool() {
      this.drawer = true;
    },
    // 通勤时间面板关闭
    drawerClose() {
      this.drawer = false;
    },
    // 地图自定义主题改变
    _mapStyleSelectedChange(value) {
      this.map.setMapStyle("amap://styles/" + value);
    },
    // 加载 上海著名IT科技园，张江、漕河泾、五角场、紫竹
    _loadItArea(add) {
      // 121.594377,31.206623 浦东新区张江高科技园区
      // 121.397769,31.170644 徐汇区漕河泾开发区
      // 121.513906,31.304645 杨浦区五角场创智天地
      // 121.45126,31.02279 闵行区紫竹科学园区
      let _this = this;
      // clear
      this.itCompanyAreaList.forEach((element) => {
        _this.map.remove(element);
      });
      this.itCompanyAreaList.slice(0, this.itCompanyAreaList.length); // clear
      if (add || add == undefined) {
        this._addArea(121.594377, 31.206623, "张江高科");
        this._addArea(121.397769, 31.170644, "漕河泾");
        this._addArea(121.513906, 31.304645, "五角场");
        this._addArea(121.45126, 31.02279, "紫竹");
      }
    },
    // 添加一个区域
    _addArea(lon, lat, title) {
      let circle = new AMap.Circle({
        center: new AMap.LngLat(lon, lat), // 圆心位置
        radius: 500, // 圆半径
        fillColor: "gray", // 圆形填充颜色
        fillOpacity: 0.2,
        strokeColor: "#fff", // 描边颜色
        strokeWeight: 2, // 描边宽度
      });
      // 创建一个 Marker 实例：
      // let marker = new AMap.Text({
      //   position: new AMap.LngLat(lon, lat),
      //   text: title
      // });
      // 小三角、和阴影
      const content = `<div class="test_triangle_border">
                        <div class="popup">
                          <em></em><span></span>${title}
                        </div>
                       </div>`;
      let marker = new AMap.Marker({
        content: content, // 自定义点标记覆盖物内容
        position: [lon, lat],
        title: title,
        offset: new AMap.Pixel(-75, -72),
      });

      this.itCompanyAreaList.push(marker);
      this.itCompanyAreaList.push(circle);

      // 将创建的点标记添加到已有的地图实例：
      this.map.add(marker);
      this.map.add(circle);
    },
    // 加载，上海著名互联网公司
    _loadItCompany(add) {
      let _this = this;
      this.itCompanyList.forEach((element) => {
        _this.map.remove(element);
      });
      this.itCompanyList.slice(0, this.itCompanyList.length); // clear
      if (add || add == undefined) {
        this._addCompany(121.506414, 31.309352, "哔哩哔哩");
        this._addCompany(121.488982, 31.255272, "网易");
        this._addCompany(121.604903, 31.179749, "陆金所");
        this._addCompany(121.605876, 31.179705, "百度");
        this._addCompany(121.62648, 31.217651, "喜马拉雅");
        this._addCompany(121.387672, 31.166692, "今日头条");
        this._addCompany(121.397336, 31.167389, "腾讯");
        this._addCompany(121.47633, 31.256971, "饿了么");
        this._addCompany(121.425928, 31.219444, "拼多多");
        this._addCompany(121.526575, 31.080261, "泛微");
        this._addCompany(121.542336, 31.27552, "UCloud");
        this._addCompany(121.384473, 31.2141, "深兰科技");
        this._addCompany(121.494233, 31.248829, "大疆创新");
        this._addCompany(121.512506, 31.306694, "声网");
        this._addCompany(121.549879, 31.226638, "蚂蚁金服");
        this._addCompany(121.603265, 31.180476, "盛大游戏");
        this._addCompany(121.531287, 31.218108, "蜻蜓FM");
        this._addCompany(121.478772, 31.205269, "点融");
        this._addCompany(121.254768, 31.330339, "小红书");
        this._addCompany(121.434285, 31.201397, "高德");
        this._addCompany(121.517926, 31.203763, "唯品会");
        this._addCompany(121.349475, 31.229919, "爱奇艺");
        this._addCompany(121.351152, 31.22074, "携程");
        this._addCompany(121.418106, 31.177024, "miHoYo");
        this._addCompany(121.363962, 31.12282, "哈啰出行");
        this._addCompany(121.26557, 31.055322, "巨人网络");
        this._addCompany(121.420766, 31.192054, "掌门一对一");
        this._addCompany(121.596878, 31.187542, "WIFI钥匙");
        this._addCompany(121.407518, 31.17156, "捞月狗");
        this._addCompany(121.392529, 31.232237, "晓黑板");
        this._addCompany(121.399583, 31.165317, "鱼泡泡");
        this._addCompany(121.580688, 31.199923, "阅文集团");
        this._addCompany(121.512966, 31.30737, "商米");
        this._addCompany(121.620012, 31.256955, "万达信息");
        this._addCompany(121.626833, 31.208961, "趣头条");
        this._addCompany(121.533813, 31.272623, "英语流利说");
        this._addCompany(121.436295, 31.18491, "轻轻家教");
        this._addCompany(121.463089, 31.02068, "Intel");
        this._addCompany(121.409258, 31.171887, "微软");
        this._addCompany(121.519372, 31.077439, "万达股份");
        this._addCompany(121.39253, 31.170216, "莉莉丝");
        this._addCompany(121.399854, 31.168252, "商汤");
      }
    },
    _addCompany(lon, lat, title) {
      // 创建一个 Marker 实例：
      const content = `<div class="test_triangle_border_2">
                        <div class="popup">
                          <em></em><span></span>${title}
                        </div>
                       </div>`;
      let marker = new AMap.Marker({
        content: content, // 自定义点标记覆盖物内容
        position: [lon, lat],
        title: title,
        offset: new AMap.Pixel(-75, -72),
      });

      this.itCompanyAreaList.push(marker);
      // 将创建的点标记添加到已有的地图实例：
      this.map.add(marker);
    },
  },
};
</script>

<style>
.full {
  width: 100%;
  height: 100%;
}
.mapStyle {
  position: absolute;
  right: 10px;
  top: 10px;
  border-radius: 5px;
  border-color: #ebebeb;
}
#map-container {
  width: 100%;
  height: 100%;
}
/*自定义图标样式*/
.test_triangle_border .popup {
  width: 50px;
  background: #027cf6;
  padding: 1px 1px 2px 1px;
  color: #ffffff;
  text-align: center;
  border-radius: 4px;
  position: absolute;
  font-size: 10px;
  top: 30px;
  left: 30px;
}
.test_triangle_border_2 .popup {
  width: 50px;
  background: #ff9966;
  padding: 1px 1px 2px 1px;
  color: #ffffff;
  text-align: center;
  border-radius: 4px;
  border-color: #ffffff;
  border-width: 1px 1px 0;
  border-style: solid;
  position: absolute;
  font-size: 10px;
  top: 30px;
  left: 30px;
}
/*向下*/
.test_triangle_border span {
  display: block;
  width: 0;
  height: 0;
  border-width: 6px 6px 0;
  border-style: solid;
  border-color: #027cf6 transparent transparent; /*黄 透明 透明 */
  position: absolute;
  bottom: -6px;
  left: 37px;
}
.test_triangle_border_2 span {
  display: block;
  width: 0;
  height: 0;
  border-width: 6px 6px 0;
  border-style: solid;
  border-color: #ff9966 transparent transparent; /*黄 透明 透明 */
  position: absolute;
  bottom: -6px;
  left: 37px;
}
/*小阴影*/
.test_triangle_border em {
  width: 10px;
  height: 5px;
  margin: 3px;
  background: #646464;
  opacity: 0.5;
  border-radius: 50% / 50%;
  position: absolute;
  bottom: -14px;
  left: 37px;
}
.test_triangle_border_2 em {
  width: 10px;
  height: 5px;
  margin: 3px;
  background: #646464;
  opacity: 0.5;
  border-radius: 50% / 50%;
  position: absolute;
  bottom: -14px;
  left: 37px;
}
</style>