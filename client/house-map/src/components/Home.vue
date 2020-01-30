<template>
  <el-container id="app">
    <el-header>
      <el-upload
        class="overButton"
        ref="uploads"
        accept=".xls, .xlsx, .csv"
        action
        :on-change="upload"
        :show-file-list="false"
        :auto-upload="false"
      >
        <el-button slot="trigger" size="small" type="primary">加载excel</el-button>
      </el-upload>
    </el-header>
    <el-main>
      <Map ref="map" />
    </el-main>
  </el-container>
</template>

<script>
import Map from "./Map.vue";
import XLSX from "xlsx";

export default {
  name: "home",
  components: {
    Map
  },
  data() {
    return {
      houseList: []
    };
  },
  methods: {
    filterHouse(element) {
      // 总价150万 - 200万
      return (
        element.TotalPrice < 200 * 10000 && element.TotalPrice > 150 * 10000
      );
    },
    upload(file, fileList) {
      let files = { 0: file.raw };
      this.loadExcel(files);
    },
    loadExcel(files) {
      var that = this;
      that.houseList.splice(0, that.houseList.length); //清空数组
      //console.log(files);
      if (files.length <= 0) {
        //如果没有文件名
        return false;
      } else if (!/\.(csv|xls|xlsx)$/.test(files[0].name.toLowerCase())) {
        this.$Message.error("上传格式不正确，请上传xls或者xlsx格式");
        return false;
      }

      const fileReader = new FileReader();
      fileReader.onload = function(ev) {
        // try {
        const data = ev.target.result;
        const workbook = XLSX.read(data, {
          type: "binary"
        });
        const wsname = workbook.SheetNames[0]; //取第一张表
        const ws = XLSX.utils.sheet_to_json(workbook.Sheets[wsname]); //生成json表格内容
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
        //console.log(ws);

        var i = 0;
        ws.forEach(element => {
          if (that.filterHouse(element)) {
            i++;
            if (i < 100) {
              //console.log(element.Title);
              //that.houseList.push(element);
            }
            that.houseList.push(element);
          }
        });

        that.loadToMap(that.houseList);
        that.$refs.uploads.value = "";
        // } catch (e) {
        //   console.log(e);
        //   return false;
        // }
      };
      fileReader.readAsBinaryString(files[0]);
    },
    loadToMap(houseList) {
      var _this = this;
      _this.$refs.map.clear();

      //return new Promise(function(resolve, reject) {
      var i = 0;
      var massMarsList = [];
      houseList.forEach(element => {
        i++;
        if (element.Location != undefined && element.Location.length >= 2) {
          var tempStr = element.Location.substring(
            1,
            element.Location.length - 2
          );
          var lonLat = tempStr.split(",");
          var lon = lonLat[0];
          var lat = lonLat[1];

          massMarsList.push({
            lnglat: [parseFloat(lon), parseFloat(lat)],
            title: element.Title+"\r\n"+element.Link,
            id: i,
            style: 0
          });
        }
      });

      _this.$refs.map.addMassMarks(massMarsList);

      //resolve();
      //});
    }
  }
};
</script>

<style>
.el-header {
  background-color: #b3c0d1;
  color: #333;
  text-align: center;
}
.el-aside {
  background-color: #d3dce6;
  color: #333;
  text-align: center;
}
.el-main {
  background-color: #e9eef3;
  color: #333;
  text-align: center;
}
.overButton {
  position: absolute;
  top: 10px;
  left: 10px;
}
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  width: 100%;
  height: 100%;
}
</style>