/*jshint esversion: 6 */

const axios = require("axios");
const kMapKey = "e8819cde9b68966210cb6ff2bf4e76d7";
const kDefaultCity = "shanghai";

/**
 * 查询POI地址对应的经纬度
 * @export
 * @param {*} address：地址
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
export async function getTransit(origin, destination, strategy) {
  console.log(
    "getTransit origin=" +
      origin +
      ",dest=" +
      destination +
      ",strategy=" +
      strategy
  );
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
