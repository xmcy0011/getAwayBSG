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
          let lng = geocodes.location.split(",")[0];
          let lat = geocodes.location.split(",")[1];
          resolve({ lng: lng, lat: lat });
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
          console.log(pois);
          resolve(pois);
        }
      })
      .catch(function(error) {
        // handle error
        console.log(error);
      });
  });
}
