<template>
  <div>
    <el-button size="small" @click="this.drawer = true" style="margin-left: 5px"
      >通勤</el-button
    >

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
        @clickSearch="clickSearch"
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
            <strong :name="o.id">无</strong>
          </el-col>
          <el-col :span="12">
            <p>驾车</p>
            <strong :name="o.id">无</strong>
          </el-col>
        </el-row>
      </el-card>
    </el-drawer>
  </div>
</template>

<script>
import PoiSearchBox from "./PoiSearchBox.vue";
import { geo, getTransit } from "../api/api.js";

export default {
  components: {
    PoiSearchBox,
  },
  data() {
    return {
      drawer: false, // 通勤时间工具
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
      this.clearLastResult();
      this.directionCalc(this.$refs.poiSearchBox.search.address);
    },
    async directionCalc(originAddress) {
      let location = await geo(originAddress);
      console.log("addr=" + originAddress + ",location=" + location);

      for (let i = 0; i < this.destination.length; i++) {
        let dest = this.destination[i].location;
        // 公交
        let results = await getTransit(location, dest, 0);
        if (results.length > 0) {
          let duration = results[0].duration;
          // results[0] 代表最快的方案
          let hour = parseInt(duration / 60 / 60);

          let min = parseInt(duration / 60) % 60;
          let text = min + "分钟";
          console.log(duration + " " + hour + " " + min);
          if (hour > 0) {
            text = hour + "小时" + text;
          }

          let kmeters = (results[0].walking_distance / 1000).toFixed(2);
          text += "，步行：" + kmeters + "千米";

          var elementArr = document.getElementsByName("dest" + i);
          console.log("dest" + i + ",element.length=" + elementArr.length);
          if (elementArr.length > 0) {
            elementArr[0].innerText = text;
          }
        }

        // 驾车
      }
    },
  },
};
</script>
