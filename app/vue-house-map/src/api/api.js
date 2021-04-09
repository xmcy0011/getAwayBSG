/*jshint esversion: 6 */

const axios = require("axios");
const kMapKey = "e8819cde9b68966210cb6ff2bf4e76d7";
const kDefaultCity = "shanghai";

/**
 * 查询POI地址对应的经纬度
 * @export
 * @param {*} address：结构化地址
 * @returns
 */
export async function geo(address) {
  return new Promise((resolve) => {
    // Make a request for a user with a given ID
    axios
      .get("https://restapi.amap.com/v3/geocode/geo", {
        params: {
          key: kMapKey,
          address: address,
          city: kDefaultCity,
        },
      })
      .then(function(response) {
        if (response.data.status == 1) {
          let geocodes = response.data.geocodes[0];
          resolve(geocodes.location);
        }
      })
      .catch(function(error) {
        // handle error
        console.log(error);
      });
  });
}

/**
 * 搜索POI建议地址列表
 * @export
 * @param {*} keywords：关键字
 * @returns
 */
export async function queryPlace(keywords) {
  return new Promise((resolve) => {
    axios
      .get("https://restapi.amap.com/v3/place/text", {
        params: {
          key: kMapKey,
          keywords: keywords,
          city: kDefaultCity,
          // types: 1,
          offset: 5, // 分页大小
        },
      })
      .then(function(response) {
        if (response.data.status == 1) {
          let pois = response.data.pois;
          //console.log(pois);
          resolve(pois);
        }
      })
      .catch(function(error) {
        // handle error
        console.log(error);
      });
  });
}

/**
 * 公交路径规划，获取起点和终点的时间
 * @export
 * @param {*} origin：起点，格式为：121,31
 * @param {*} destination：终点
 * @param {*} strategy：公交换乘策略，0：最快捷模式，1：最经济模式，2：最少换乘模式，3：最少步行模式，5：不乘地铁模式
 * @returns [{cost: "4.0", duration: "3230", walking_distance: "1579", distance: "18359"}]
 * cost：费用
 * duration：预期时间，秒
 * walking_distance：步行距离，米
 * distance：总距离，米
 */
export async function getTransitIntegrated(origin, destination, strategy) {
  return new Promise((resolve) => {
    axios
      .get("https://restapi.amap.com/v3/direction/transit/integrated", {
        params: {
          key: kMapKey,
          origin: origin,
          destination: destination,
          city: kDefaultCity,
          cityd: kDefaultCity, // 跨城公交规划时的终点城市
          extensions: "base", // 返回基本结果
          strategy: strategy,
          nightflag: 0, // 不计算夜班车
        },
      })
      .then(function(response) {
        console.log(response);
        let data = response.data;
        if (data.status == 1) {
          if (data.infocode == "10000") {
            let route = data.route;
            let transitsArr = route.transits;
            resolve(transitsArr);
          } else {
            console.log(response.data);
          }
        }
      })
      .catch(function(error) {
        // handle error
        console.log(error);
      });
  });
}

/**
 * 驾车路径规划，获取起点和终点的时间
 * @export
 * @param {*} origin：起点，格式为：121,31
 * @param {*} destination：终点
 * @param {*} strategy：策略：
1，费用优先，不走收费路段，且耗时最少的路线
2，距离优先，不考虑路况，仅走距离最短的路线，但是可能存在穿越小路/小区的情况
3，速度优先，不走快速路，例如京通快速路（因为策略迭代，建议使用13）
4，躲避拥堵，但是可能会存在绕路的情况，耗时可能较长
5，多策略（同时使用速度优先、费用优先、距离优先三个策略计算路径）。
其中必须说明，就算使用三个策略算路，会根据路况不固定的返回一~三条路径规划信息。
6，速度优先，不走高速，但是不排除走其余收费路段
7，费用优先，不走高速且避免所有收费路段
8，躲避拥堵和收费，可能存在走高速的情况，并且考虑路况不走拥堵路线，但有可能存在绕路和时间较长
9，躲避拥堵和收费，不走高速
 * @returns [{distance: "22124", duration: "2298", traffic_lights: "13"}]
 */
export async function getTransitDriving(origin, destination, strategy) {
  return new Promise((resolve) => {
    axios
      .get("https://restapi.amap.com/v3/direction/driving", {
        params: {
          key: kMapKey,
          origin: origin,
          destination: destination,
          strategy: strategy,
          extensions: "base",
          province: "沪", // 用于判断是否限行
        },
      })
      .then(function(response) {
        console.log(response);
        let data = response.data;
        if (data.status == 1) {
          if (data.infocode == "10000") {
            let route = data.route;
            let paths = route.paths;
            resolve(paths);
          } else {
            console.log(response.data);
          }
        }
      })
      .catch(function(error) {
        // handle error
        console.log(error);
      });
  });
}
