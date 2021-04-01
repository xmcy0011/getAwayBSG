<template>
  <div>
    <el-input
      v-model="search.name"
      @focus="focus"
      @blur="blur"
      @input="onSearch"
      @keyup.enter="onSearch"
      placeholder="搜索地点"
    >
      <template #append>
        <el-button
          icon="el-icon-search"
          @click="onClickSearch"
          tooltip="搜索"
        ></el-button>
        <el-button
          icon="el-icon-document-copy"
          style="margin-left: 2px"
          @click="onClickCopy"
        ></el-button>
      </template>
    </el-input>

    <!---设置z-index优先级,防止被走马灯效果遮挡-->
    <el-card
      @mouseenter="enterSearchBoxHanlder"
      v-if="isSearch"
      id="search-box"
      style="position: relative; z-index: 15; cursor: pointer"
    >
      <dl v-if="isSearchList">
        <dd
          v-for="search in searchList"
          :key="search.name"
          @click="this.search = search"
        >
          {{ search.name }}
          <!-- ({{ search.address }}）-->
        </dd>
      </dl>
    </el-card>
  </div>
</template>

<script>
import { queryPlace } from "../api/api.js";

export default {
  data() {
    return {
      search: { name: "" }, //当前输入框的值
      isFocus: false, //是否聚焦
      searchList: [], //搜索返回数据,
    };
  },
  computed: {
    isSearchList() {
      return this.isFocus && this.search.name != "";
    },
    isSearch() {
      return this.isFocus;
    },
  },
  methods: {
    focus() {
      if (this.search.name != "") {
        this.searchPoi();
      }
      this.isFocus = true;
    },
    blur() {
      var self = this;
      this.searchBoxTimeout = setTimeout(function () {
        self.isFocus = false;
        self.searchList = [{ name: "" }];
      }, 300);
    },
    enterSearchBoxHanlder() {
      clearTimeout(this.searchBoxTimeout);
    },
    // 点击确定按钮
    onClickSearch() {
      this.$emit("clickSearch", this.search);
    },
    // 点击拷贝
    onClickCopy() {
      this.$emit("clickCopy");
    },
    onSearch() {
      this.searchPoi();
    },
    // 搜索POI
    async searchPoi() {
      let pois = await queryPlace(this.search);
      console.log(pois);
      this.searchList = pois;
    },
  },
};
</script>