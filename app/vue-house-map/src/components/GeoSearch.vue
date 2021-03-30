<template>
  <el-row>
    <el-col></el-col>
    <el-col :span="6">
      <el-input
        v-model="search"
        @focus="focus"
        @blur="blur"
        @keyup.enter="onClickSearch"
        placeholder="搜索地点"
      >
        <el-button style="margin-left: 10px" @click="onClickSearch"
          >确定</el-button
        >
      </el-input>

      <!---设置z-index优先级,防止被走马灯效果遮挡-->
      <el-card
        @mouseenter="enterSearchBoxHanlder"
        v-if="isSearch"
        class="box-card"
        id="search-box"
        style="position: relative; z-index: 15"
      >
        <dl v-if="isSearchList">
          <dd
            v-for="search in searchList"
            :key="search.name"
            @click="this.search = search.name"
          >
            {{ search.name }}（{{ search.address }}）
          </dd>
        </dl>
      </el-card>
    </el-col>
  </el-row>
</template>

<script>
import { queryPlace } from "../api/api.js";

export default {
  data() {
    return {
      search: "", //当前输入框的值
      isFocus: false, //是否聚焦
      searchList: [], //搜索返回数据,
    };
  },
  computed: {
    isSearchList() {
      return this.isFocus && this.search;
    },
    isSearch() {
      return this.isFocus;
    },
  },
  methods: {
    focus() {
      this.isFocus = true;
    },
    blur() {
      var self = this;
      this.searchBoxTimeout = setTimeout(function () {
        self.isFocus = false;
      }, 300);
    },
    enterSearchBoxHanlder() {
      clearTimeout(this.searchBoxTimeout);
    },
    onClickSearch() {
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