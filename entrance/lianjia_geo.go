package entrance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getAwayBSG/pkg/configs"
	"github.com/getAwayBSG/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var client *mongo.Client
var cityNames map[string]string
// key:villageName:cityName
var cacheGeoTables map[string]LatLon

const geoCodeUrl = "https://restapi.amap.com/v3/geocode/geo"
const geoCodeKey = "e8819cde9b68966210cb6ff2bf4e76d7"

type LatLon struct {
	FormattedAddress string `json:"formatted_address"`
	Location         string `json:"location"`
}
type GeoResult struct {
	Status   string   `json:"status"`
	Count    string   `json:"count"`
	GeoCodes []LatLon `json:"geocodes"`
}
type GeoHouseInfo struct {
	ListAreaName string `json:"ListAreaName"`
	Link         string `json:"Link"`
	City         string `json:"City"`
	Title        string `json:"Title"`
}

var (
	totalCount   = 0
	successCount = 0
	startTime    = time.Now()
)

func init() {
	cacheGeoTables = make(map[string]LatLon, 0)
	cityNames = make(map[string]string, 0)
	cityNames["sh"] = "上海"
	cityNames["cs"] = "长沙"
	cityNames["bj"] = "北京"

	cityNames["hz"] = "杭州"
	cityNames["cd"] = "成都"
	cityNames["wh"] = "武汉"
	cityNames["su"] = "苏州"
	cityNames["nb"] = "宁波"
	cityNames["gz"] = "广州"
	cityNames["sz"] = "深圳"

	cityNames["cq"] = "重庆"
	cityNames["zz"] = "郑州"
	cityNames["nj"] = "南京"
	cityNames["xa"] = "西安"
	cityNames["tj"] = "天津"
	cityNames["sy"] = "沈阳"
	cityNames["qd"] = "青岛"
	cityNames["dg"] = "东莞"
	cityNames["km"] = "昆明"
}

func initMongodb() error {
	// 初始化中间表，把小区经纬度存入到mongodb，没有命中才调用高德API，提供性能和速度
	tempClient, err := mongo.NewClient(options.Client().ApplyURI(configs.ConfigInfo.DbRrl + "/" + configs.ConfigInfo.DbDatabase))
	if err != nil {
		return err
	}

	err = tempClient.Connect(context.Background())
	if err != nil {
		return err
	}

	err = tempClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}

	client = tempClient
	return nil
}

// 经纬度表
func dbGeoCollection() *mongo.Collection {
	db := client.Database(configs.ConfigInfo.DbDatabase)
	return db.Collection("geo")
}

// 二手房表
func dbLjCollection() *mongo.Collection {
	db := client.Database(configs.ConfigInfo.DbDatabase)
	return db.Collection(configs.ConfigInfo.TwoHandHouseCollection)
}

// 获取小区对应的经纬度
func getVillageLatLon(cityName, villageName string) (lat, lon float32, formattedAddress string, err error) {
	// 先找缓存
	items, ok := cacheGeoTables[villageName+":"+cityName]
	if ok {
		latLon := strings.Split(items.Location, ",")
		lon, err := strconv.ParseFloat(latLon[0], 10)
		if err != nil {
			logger.Sugar.Error(err)
			return 0, 0, "", nil
		}

		lat, err := strconv.ParseFloat(latLon[1], 10)
		if err != nil {
			logger.Sugar.Error(err)
			return 0, 0, "", nil
		}
		return float32(lat), float32(lon), items.FormattedAddress, nil
	}

	// 再找数据库
	table := dbGeoCollection()
	if table != nil {
		r := table.FindOne(context.Background(), bson.M{"villageName": villageName, "cityName": cityName})
		row := bson.M{}
		err := r.Decode(&row)
		if err != nil { //not find
			return 0, 0, "", nil
		}

		lat := row["lat"].(float64)
		lon := row["lon"].(float64)
		formattedAddress := row["formattedAddress"].(string)

		// 加入到缓存中
		location := LatLon{}
		// 6位小数点即可
		location.Location = fmt.Sprintf("%.6f", lon) + "," + fmt.Sprintf("%.6f", lat)
		location.FormattedAddress = formattedAddress
		cacheGeoTables[villageName+":"+cityName] = location
		return float32(lat), float32(lon), formattedAddress, nil
	}

	return 0, 0, "", errors.New("db not connect")
}

// 添加小区对应的经纬度
func addVillageLatLon(cityName, villageName string, lat, lon float32, formattedAddress string) {
	table := dbGeoCollection()
	if table != nil {
		tempLon := fmt.Sprintf("%.6f", lon)
		tempLat := fmt.Sprintf("%.6f", lat)
		fmtLon, _ := strconv.ParseFloat(tempLon, 10)
		fmtLat, _ := strconv.ParseFloat(tempLat, 10)

		_, err := table.InsertOne(context.Background(), bson.M{"villageName": villageName, "cityName": cityName,
			"lat": fmtLat, "lon": fmtLon, "formattedAddress": formattedAddress})
		if err != nil {
			logger.Sugar.Error(err)
		} else {
			// 加入到缓存中
			location := LatLon{}
			// 6位小数点即可
			location.Location = tempLon + "," + tempLat
			location.FormattedAddress = formattedAddress
			cacheGeoTables[villageName+":"+cityName] = location
		}
	}
}

// 链家二手房，补充房源经纬度信息
func updateLjVillageLatLon(link string, lon, lat float32, formattedAddress string) {
	lj := dbLjCollection()
	if lj != nil {
		location := make([]float32, 0)
		location = append(location, lon)
		location = append(location, lat)
		_, err := lj.UpdateOne(context.Background(), bson.M{"Link": link}, bson.M{"$set": bson.M{"Location": location, "FormattedAddress": formattedAddress}})
		if err != nil {
			if !strings.Contains(err.Error(), "multiple write errors") {
				logger.Sugar.Errorf("数据库插入失败:%s", err.Error())
			}
		}
	}
}

// 高德地理编码小区(http api)，获取对应经纬度
func getLocationFromUrl(cityName, villageName string) (lat float32, lon float32, formattedAddress string, err error) {
	url := ""
	searchName, ok := cityNames[cityName]
	if ok {
		url = fmt.Sprintf("%s?address=%s&key=%s&city=%s", geoCodeUrl, villageName, geoCodeKey, searchName)
	} else {
		url = fmt.Sprintf("%s?address=%s&key=%s", geoCodeUrl, villageName, geoCodeKey)
		//logger.Sugar.Warn("can't find cityName in cityNames(map[string]string),geo code search in whole china.")
	}

	res, err := http.Get(url)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, 0, "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, 0, "", err
	}

	jsonInfo := &GeoResult{}
	err = json.Unmarshal(body, &jsonInfo)
	if err != nil {
		logger.Sugar.Error(err)
		return 0, 0, "", err
	}

	if jsonInfo.Status != "1" {
		logger.Sugar.Errorf("geocode error:%s,villageName=%s,city=%s", jsonInfo.Status, villageName, cityName)
	} else {
		// 可能转码错误，请手动拾取地理坐标：https://lbs.amap.com/console/show/picker
		if jsonInfo.Count != "0" && len(jsonInfo.GeoCodes) > 0 {
			// 没有办法判断，先以高德返回的为主吧
			//if strings.Contains(jsonInfo.GeoCodes[0].FormattedAddress, villageName) {
			formattedAddress = jsonInfo.GeoCodes[0].FormattedAddress
			geoArr := strings.Split(jsonInfo.GeoCodes[0].Location, ",")
			if len(geoArr) < 2 {
				logger.Sugar.Errorf("error location format:%s,villageName=%s,city=%s",
					jsonInfo.GeoCodes[0].Location, villageName, cityName)
			} else {
				lon, err := strconv.ParseFloat(geoArr[0], 10)
				if err != nil {
					return 0, 0, "", err
				}
				lat, err := strconv.ParseFloat(geoArr[1], 10)
				if err != nil {
					return 0, 0, "", err
				}
				return float32(lat), float32(lon), formattedAddress, nil
			}
			//} else {
			//	return 0, 0, errors.New("result is not correct,formattedAddress=" +
			//		jsonInfo.GeoCodes[0].FormattedAddress + ",villageName=" + villageName)
			//}
		} else {
			return 0, 0, "", errors.New("have no result,status=" + jsonInfo.Status + ",GeoCodes.len=" + strconv.Itoa(len(jsonInfo.GeoCodes)))
		}
	}

	return 0, 0, "", errors.New("error")
}

// 地理编码某个小区
func geocode(cityName, villageName, link string) (lonOut, latOut float32, err error) {
	lat, lon, formattedAddress, err := getVillageLatLon(cityName, villageName)
	if err != nil {
		return 0, 0, err
	}

	if lat == 0 && lon == 0 {
		lat, lon, formattedAddress, err = getLocationFromUrl(cityName, villageName)
		if err != nil {
			return 0, 0, err
		}
		// 把小区的经纬度缓存到对应的表
		addVillageLatLon(cityName, villageName, lat, lon, formattedAddress)
	}

	// 更新房源经纬度
	updateLjVillageLatLon(link, lon, lat, formattedAddress)
	return lon, lat, err
}

// 优雅退出（退出信号）
func waitElegantExit(signalChan chan os.Signal) {
	go func() {
		<-signalChan

		timeDiff := time.Now().Sub(startTime)
		logger.Sugar.Infof("主动退出，成功：%d，失败：%d，总数：%d，用时：%.2f 秒=%.2f 分=%.2f 时", successCount, totalCount-successCount,
			totalCount, timeDiff.Seconds(), timeDiff.Minutes(), timeDiff.Hours())
		logger.Sugar.Info("如果需要继续补充，请再次执行程序，选择3即可！")

		os.Exit(2)
	}()
}

// 使用高德地理编码服务
// 补充链家所有二手房小区的经纬度，便于查询分析靠近地铁站1公里的房源
func StartGeocodeLJ() {
	err := initMongodb()
	if err != nil {
		logger.Sugar.Errorf("mongodb connect error:%s", err)
		return
	}

	signalChan := make(chan os.Signal, 1)
	// 注册CTRL+C：被打断通道,syscall.SIGINT
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	waitElegantExit(signalChan)

	//isAll := true

	lj := dbLjCollection()
	if lj != nil {
		logger.Sugar.Infof("查找没有经纬度的房源中，文档名=%s...", configs.ConfigInfo.TwoHandHouseCollection)
		cursor, err := lj.Find(context.Background(), bson.M{"Location": bson.M{"$exists": false}})
		defer cursor.Close(context.Background())

		if err != nil {
			logger.Sugar.Error(err)
			return
		}

		// 一次性加载结果到内存
		logger.Sugar.Info("读取结果到内存...")
		var houseArr = make([]GeoHouseInfo, 0)

		//if isAll {
		err = cursor.All(context.Background(), &houseArr)
		if err != nil {
			logger.Sugar.Error(err)
			return
		}
		//} else {
		//	timeoutContext, _ := context.WithTimeout(context.Background(), time.Duration(time.Second))
		//	const testLimitCount = 1000
		//	for i := 0; i < testLimitCount; i++ {
		//		if cursor.Next(timeoutContext) {
		//			geo := GeoHouseInfo{}
		//			err := cursor.Decode(&geo)
		//			if err != nil {
		//				logger.Sugar.Error(err)
		//			} else {
		//				houseArr = append(houseArr, geo)
		//			}
		//		} else {
		//			break
		//		}
		//		if i != 0 && i%100 == 0 {
		//			logger.Sugar.Infof("已读取%d", i)
		//		}
		//	}
		//}
		totalCount = len(houseArr)

		logger.Sugar.Infof("3秒后，开始反地理编码，总数=%d，如果需要中断请按下Ctrl+C", totalCount)
		time.Sleep(time.Duration(time.Second * 3))
		for i := range houseArr {
			item := houseArr[i]

			if item.City == "" || item.ListAreaName == "" {
				logger.Sugar.Errorf("没有城市或小区信息，城市：%s,小区：%s,标题：%s,Link：%s", item.City, item.ListAreaName, item.Title, item.Link)
				continue
			}

			lon, lat, err := geocode(item.City, item.ListAreaName, item.Link)
			if err != nil {
				logger.Sugar.Errorf("城市：%s,小区：%s,标题：%s,Link：%s 补充经纬度错误：%s", item.City, item.ListAreaName, item.Title, item.Link, err.Error())
			} else {
				successCount++
				logger.Sugar.Infof("[%d/%d]城市：%s,小区：%s,标题：%s,经纬度:%f,%f", i, totalCount, item.City, item.ListAreaName, item.Title, lon, lat)
			}

			//time.Sleep(10 * time.Millisecond)
		}
	}

	timeDiff := time.Now().Sub(startTime)
	logger.Sugar.Infof("补充结束,成功：%d，失败：%d，总数：%d，用时：%.2f 秒=%.2f 分=%.2f 时", successCount, totalCount-successCount,
		totalCount, timeDiff.Seconds(), timeDiff.Minutes(), timeDiff.Hours())
}
