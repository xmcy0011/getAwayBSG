<template>
  <div>
    <!-- <el-button size="small" @click="this.drawer = true" style="margin-left: 5px"
      >通勤</el-button
    > -->
    <el-drawer
      title="通勤时间测算"
      v-model="drawer"
      :direction="rtl"
      :before-close="drawerClose"
      destroy-on-close
      style="overflow: scroll"
    >
      <PoiSearchBox
        ref="poiSearchBox"
        v-show="drawer"
        @clickSearch="clickSearch"
        @clickCopy="clickCopy"
        style="margin-left: 10px; margin-right: 10px"
      />

      <el-card
        class="box-card"
        style="margin-top: 20px"
        v-for="o in destination"
        body-style="padding:15px;"
        :key="o"
      >
        <template #header>
          <div class="card-header" style="height: 15px">
            <span>{{ o.name }}</span>
          </div>
        </template>
        <el-row style="margin-top: -20px">
          <el-col :span="12">
            <p>公交</p>
            <p :name="o.id">无</p>
          </el-col>
          <el-col :span="12">
            <p>驾车</p>
            <p :name="o.id">无</p>
          </el-col>
        </el-row>
      </el-card>
    </el-drawer>
  </div>
</template>

<script>
import PoiSearchBox from "./PoiSearchBox.vue";
import { geo, getTransitIntegrated, getTransitDriving } from "../api/api.js";

export default {
  components: {
    PoiSearchBox,
  },
  data() {
    return {
      drawer: false, // 通勤时间工具
      searchResultCopyText: "",
      destination: [
        {
          id: "dest0",
          name: "哔哩哔哩",
          location: "121.506414,31.309352",
        },
        {
          id: "dest1",
          name: "拼多多",
          location: "121.425928,31.219444",
        },
        {
          id: "dest2",
          name: "万达股份",
          location: "121.519372,31.077439",
        },
        {
          id: "dest3",
          name: "莉莉丝",
          location: "121.39253,31.170216",
        },
      ],
    };
  },
  methods: {
    // 通勤时间面板关闭
    drawerClose() {
      this.drawer = false;
    },
    clearLastResult() {
      for (let i = 0; i < this.destination.length; i++) {
        var elementArr = document.getElementsByName("dest" + i);
        if (elementArr.length >= 2) {
          elementArr[0].innerText = "查询中";
          elementArr[1].innerText = "查询中";
        }
      }
    },
    // 点击搜索
    clickSearch(poi) {
      if (this.$refs.poiSearchBox.search.name == "") {
        this.$notify({
          title: "警告",
          message: "请输入地址",
          type: "warning",
          //position: "bottom-right",
        });
        return;
      }
      this.clearLastResult();
      this.directionCalc(this.$refs.poiSearchBox.search.address);
    },
    copy(value) {
      // 动态创建 textarea 标签
      const textarea = document.createElement("textarea");
      // 将该 textarea 设为 readonly 防止 iOS 下自动唤起键盘，同时将 textarea 移出可视区域
      textarea.readOnly = "readonly";
      textarea.style.position = "absolute";
      textarea.style.left = "-9999px";
      // 将要 copy 的值赋给 textarea 标签的 value 属性
      textarea.value = value;
      // 将 textarea 插入到 body 中
      document.body.appendChild(textarea);
      // 选中值并复制
      textarea.select();
      textarea.setSelectionRange(0, textarea.value.length);
      document.execCommand("Copy");
      document.body.removeChild(textarea);
    },
    // 点击拷贝
    clickCopy() {
      if (this.searchResultCopyText == "") {
        this.$notify({
          title: "警告",
          message: "请搜索后再进行复制！",
          type: "warning",
        });
        return;
      }
      this.copy(this.searchResultCopyText);
      this.$notify({
        title: "拷贝到剪贴板成功",
        dangerouslyUseHTMLString: true,
        message: "<p>" + this.searchResultCopyText + "</p>",
      });
    },
    async directionCalc(originAddress) {
      let location = await geo(originAddress);
      console.log("addr=" + originAddress + ",location=" + location);

      const getTimeText = function (timeSpan) {
        let hour = parseInt(timeSpan / 60 / 60);

        let min = parseInt(timeSpan / 60) % 60;
        let text = min + "分钟";
        if (hour > 0) {
          text = hour + "小时" + text;
        }

        return text;
      };

      // 出发地：originAddress
      // 目的地：
      // bilibili\t公交：xx（）\t驾车：xx（）\n
      this.searchResultCopyText = "出发地：" + originAddress + "\n目的地：";

      for (let i = 0; i < this.destination.length; i++) {
        let dest = this.destination[i].location;

        this.searchResultCopyText += "\n" + this.destination[i].name;

        // 公交查询
        let results = await getTransitIntegrated(location, dest, 0);
        if (results.length > 0) {
          // results[0] 代表最快的方案
          let duration = results[0].duration;
          let text = getTimeText(duration);

          let kmeters = (results[0].walking_distance / 1000).toFixed(2);
          text += " 步行：" + kmeters + "千米";

          var elementArr = document.getElementsByName("dest" + i);
          console.log("dest" + i + ",element.length=" + elementArr.length);
          if (elementArr.length > 0) {
            elementArr[0].innerText = text;
          }
          this.searchResultCopyText += "\t公交：\t" + text + "\t";
        }

        // 驾车
        results = await getTransitDriving(location, dest, 0);
        if (results.length > 0) {
          // results[0] 代表最快的方案
          let duration = results[0].duration;
          let text = getTimeText(duration);

          let kmeters = (results[0].distance / 1000).toFixed(2);
          text += " " + kmeters + "千米";

          var elementArr = document.getElementsByName("dest" + i);
          console.log("dest" + i + ",element.length=" + elementArr.length);
          if (elementArr.length > 0) {
            elementArr[1].innerText = text;
          }
          this.searchResultCopyText += "\t驾车：\t" + text;
        }
      }
    },
  },
};
</script>

<style scoped>
.card-header {
  font-weight: bold;
}
</style>