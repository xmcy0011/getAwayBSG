<template>
  <el-container id="app">
    <el-header>
      <el-upload
        class="overButton"
        ref="upload"
        accept=".xls, .xlsx"
        action
        :on-change="upload"
        :show-file-list="false"
        :auto-upload="false"
      >
        <el-button slot="trigger" size="small" type="primary">加载excel</el-button>
      </el-upload>
    </el-header>
    <el-main>
      <Map />
    </el-main>
  </el-container>
</template>

<script>
import Map from "./Map.vue";

export default {
  name: "home",
  components: {
    Map
  },
  methods: {
    upload() {
      console.log("file", file);
      console.log("fileList", fileList);
      let files = { 0: file.raw };
      this.loadExcel(files);
    },
    loadExcel(files) {
      var that = this;
      console.log(files);
      if (files.length <= 0) {
        //如果没有文件名
        return false;
      } else if (!/\.(xls|xlsx)$/.test(files[0].name.toLowerCase())) {
        this.$Message.error("上传格式不正确，请上传xls或者xlsx格式");
        return false;
      }

      const fileReader = new FileReader();
      fileReader.onload = ev => {
        try {
          const data = ev.target.result;
          const workbook = XLSX.read(data, {
            type: "binary"
          });
          const wsname = workbook.SheetNames[0]; //取第一张表
          const ws = XLSX.utils.sheet_to_json(workbook.Sheets[wsname]); //生成json表格内容
          console.log(ws);
          // that.peopleArr = [];//清空接收数据
          // if(that.peopleArr.length == 1 && that.peopleArr[0].roleName == "" && that.peopleArr[0].enLine == ""){
          //     that.peopleArr = [];
          // }
          //重写数据
          try {
          } catch (err) {
            console.log(err);
          }

          this.$refs.upload.value = "";
        } catch (e) {
          return false;
        }
      };
      fileReader.readAsBinaryString(files[0]);
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