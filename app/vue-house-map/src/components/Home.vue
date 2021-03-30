<template>
  <el-container id="app">
    <Map ref="map" />
    <Direction ref="directionPanel" />

    <el-upload
      class="overButton"
      ref="uploads"
      accept=".xls, .xlsx, .csv"
      action
      :on-change="upload"
      :show-file-list="false"
      :auto-upload="false"
    >
      <el-button size="small" type="primary" icon="el-icon-upload"
        >打开excel</el-button
      >
    </el-upload>

    <el-button
      size="small"
      type="primary"
      class="overButton2"
      @click="dialogFormVisible = true"
      >筛选</el-button
    >

    <el-button size="small" class="overButton3" @click="onClickDirection"
      >通勤</el-button
    >

    <el-dialog title="筛选" v-model="dialogFormVisible" width="600px">
      <el-form :model="form">
        <el-form-item label="售价" :label-width="formLabelWidth">
          <el-radio-group class="leftRadio" v-model="form.radioPriceSelected">
            <el-radio
              v-for="arr in form.radioPriceArr"
              :label="arr.value"
              :key="arr.value"
            >
              {{ arr.label }}
            </el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="房型" :label-width="formLabelWidth">
          <el-radio-group class="leftRadio" v-model="form.houseTypeSelected">
            <el-radio
              v-for="arr in form.houseTypeArr"
              :label="arr.value"
              :key="arr.value"
            >
              {{ arr.label }}
            </el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="面积" :label-width="formLabelWidth">
          <el-radio-group class="leftRadio" v-model="form.houseSizeSelected">
            <el-radio
              v-for="arr in form.houseSizeArr"
              :label="arr.value"
              :key="arr.value"
            >
              {{ arr.label }}
            </el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogFormVisible = false">取 消</el-button>
          <el-button type="primary" @click="onSearch">确 定</el-button>
        </span>
      </template>
    </el-dialog>

    <p id="searchValue" style="position: absolute; left: 190px">
      {{ searchText }}
    </p>
  </el-container>
</template>

<script>
import Map from "./Map.vue";
import PoiSearchBox from "./PoiSearchBox.vue";
import Direction from "./Direction.vue";
import XLSX from "xlsx";

export default {
  name: "home",
  components: {
    Map,
    Direction,
  },
  data() {
    return {
      houseList: [],
      dialogFormVisible: false,
      formLabelWidth: "50",
      maxLimitShowCount: 10 * 1000, // 10K
      form: {
        radioPriceSelected: -1, // 200万以下
        radioPriceArr: [
          { value: "0", label: "全部" },
          { value: "1", label: "200万以下" },
          { value: "2", label: "200-250万" },
          { value: "3", label: "250万以上" },
        ],
        houseTypeSelected: -1, // 1室
        houseTypeArr: [
          { value: "0", label: "全部" },
          { value: "1", label: "一室" },
          { value: "2", label: "二室" },
          { value: "3", label: "三室" },
          { value: "4", label: "四室及以上" },
        ],
        houseSizeSelected: -1, // 40平米以下
        houseSizeArr: [
          { value: "0", label: "全部" },
          { value: "1", label: "40㎡以下" },
          { value: "2", label: "40-60㎡" },
          { value: "3", label: "60-90㎡" },
          { value: "4", label: "90㎡以上" },
        ],
      },

      searchText: "无",
    };
  },
  methods: {
    filterHouse(element) {
      var f = true;
      var s = true;
      var t = true;

      if (this.form.radioPrice > 0) {
        // 1:x<200  2:200<x<250  3:250<x
        if (this.form.radioPrice == 1) {
          f = element.TotalPrice <= 200 * 10000;
        } else if (this.form.radioPrice == 2) {
          f =
            element.TotalPrice > 200 * 10000 &&
            element.TotalPrice <= 250 * 10000;
        } else if (this.form.radioPrice == 3) {
          f = element.TotalPrice > 250 * 10000;
        }
      }

      if (this.form.houseType > 0) {
        if (this.form.houseType == 1) {
          s = element.ListHouseType.indexOf("1室") > -1;
        } else if (this.form.houseType == 2) {
          s = element.ListHouseType.indexOf("2室") > -1;
        } else if (this.form.houseType == 3) {
          s = element.ListHouseType.indexOf("3室") > -1;
        } else {
          s = true;
        }
      }

      if (this.form.houseSize > 0) {
        if (this.form.houseSize == 1) {
          t = element.ListHouseSize <= 40;
        } else if (this.form.houseSize == 2) {
          t = element.ListHouseSize > 40 && element.ListHouseSize <= 60;
        } else if (this.form.houseSize == 3) {
          t = element.ListHouseSize > 60 && element.ListHouseSize <= 90;
        } else {
          t = element.ListHouseSize > 90;
        }
      }

      return f && s && t;
      //return element.Tag.indexOf("近地铁") > -1;
    },
    upload(file, fileList) {
      let files = { 0: file.raw };
      const loading = this.$loading({
        lock: true,
        text: "Loading",
        spinner: "el-icon-loading",
        background: "rgba(0, 0, 0, 0.7)",
      });

      this.loadExcel(files);
      // setTimeout(() => {
      loading.close();
      // }, 2000);
    },
    loadExcel(files) {
      var that = this;
      that.houseList.splice(0, that.houseList.length); //清空数组
      //console.log(files);
      if (files.length <= 0) {
        this.$Message.error("请选择文件");
        return false;
      } else if (!/\.(csv|xls|xlsx)$/.test(files[0].name.toLowerCase())) {
        this.$Message.error("上传格式不正确，请上传xls或者xlsx格式");
        return false;
      }

      const fileReader = new FileReader();
      fileReader.onload = function (ev) {
        // try {
        const data = ev.target.result;
        const workbook = XLSX.read(data, {
          type: "binary",
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
        ws.forEach((element) => {
          that.houseList.push(element);
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
      houseList.forEach((element) => {
        if (_this.filterHouse(element)) {
          i++;
          if (i >= _this.maxLimitShowCount) {
            this.$notify({
              title: "警告",
              message:
                "最大加载" + _this.maxLimitShowCount + "条，请调整范围后重试！",
              type: "warning",
            });
            return;
          }
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
              title: element.Title + "\r\n" + element.Link,
              id: i,
              full: element,
              style: 0,
            });
          }
        }
      });

      _this.$refs.map.addMassMarks(massMarsList);

      //resolve();
      //});
    },
    // 打开通勤面板
    onClickDirection() {
      this.$refs.directionPanel.drawer = true;
    },
    onSearch() {
      this.dialogFormVisible = false;

      //"价格：全部，房型：全部，大小：全部",
      const f = this.form;
      this.searchText =
        "价格：" +
        f.radioPriceArr[f.radioPriceSelected].label +
        "，房型：" +
        f.houseTypeArr[f.houseTypeSelected].label +
        "，大小：" +
        f.houseSizeArr[f.houseSizeSelected].label;
      var _this = this;
      setTimeout(() => {
        _this.loadToMap(_this.houseList);
      }, 200);
    },
  },
};
</script>

<style>
.el-header {
  background-color: #fbfbfb;
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
.overButton2 {
  position: absolute;
  top: 10px;
  left: 120px;
}
.overButton3 {
  position: absolute;
  top: 10px;
  right: 180px;
}
.searchBox {
  width: 600px;
  height: auto;
  position: absolute;
  top: 50px;
  background-color: #fbfbfb;
  border-width: 1px;
  border-style: solid;
  border-color: #f1f1f1; /*黄 透明 透明 */
}
.leftRadio {
  float: left;
  margin-top: 10px;
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
.searchPoiBox {
  position: absolute;
  width: 300px;
  top: 50px;
  right: 10px;
}
</style>