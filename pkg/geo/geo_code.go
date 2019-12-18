package geo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/getAwayBSG/pkg/logger"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const filePath = "pkg/geo/subway.json"                  // 地铁站点名字
const outputGeocodeFilePath = "pkg/geo/subway-geo.json" // 地铁经纬输出文件路径
const geoCodeUrl = "https://restapi.amap.com/v3/geocode/geo"
const geoCodeKey = "e8819cde9b68966210cb6ff2bf4e76d7"

type SubwayInfo struct {
	City string
	Line string
	Name string
}

type LatLon struct {
	FormattedAddress string `json:"formatted_address"`
	Location         string `json:"location"`
}
type GeoResult struct {
	Status   string   `json:"status"`
	Geocodes []LatLon `json:"geocodes"`
}

func StartGeoCode() {
	geoCodeCitySubway("上海", outputGeocodeFilePath)
}

func loadCitySubway(cityName string) ([]*SubwayInfo, error) {
	fs, err := os.Open(filePath)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}
	defer fs.Close()

	subways := make([]*SubwayInfo, 0)

	reader := bufio.NewReader(fs)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		// city,line,name
		arr := strings.Split(string(line), ",")
		if arr[0] == cityName {
			temp := &SubwayInfo{
				City: arr[0],
				Line: arr[1],
				Name: arr[2],
			}
			subways = append(subways, temp)
		}
	}
	return subways, nil
}

func geoCodeCitySubway(city string, outputFilePath string) {
	subwayArr, err := loadCitySubway(city)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	fs, err := os.Create(outputFilePath)
	if err != nil {
		logger.Sugar.Error(err)
		return
	}
	defer fs.Close()

	fsWrite := bufio.NewWriter(fs)

	for i := range subwayArr {
		time.Sleep(time.Millisecond * 100)

		if i%20 == 0 {
			fsWrite.Flush()
		}

		subway := subwayArr[i]
		//keywords := subway.City + subway.Line + subway.Name + "地铁站"
		keywords := subway.Name + "地铁站"
		res, err := http.Get(fmt.Sprintf("%s?address=%s&key=%s&city=%s", geoCodeUrl, keywords, geoCodeKey, city))
		if err != nil {
			logger.Sugar.Error(err)
			continue
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logger.Sugar.Error(err)
			continue
		}

		jsonInfo := &GeoResult{}
		err = json.Unmarshal(body, &jsonInfo)
		if err != nil {
			logger.Sugar.Error(err)
		} else {
			if jsonInfo.Status == "1" {
				// 可能转码错误，请手动拾取地理坐标：https://lbs.amap.com/console/show/picker
				if len(jsonInfo.Geocodes) > 0 && strings.Contains(jsonInfo.Geocodes[0].FormattedAddress, subway.Name+"地铁站") {
					logger.Sugar.Debugf("poi:%s,location:%s", subway, jsonInfo.Geocodes[0].Location)
					_, err := fsWrite.WriteString(fmt.Sprintf("%s,%s,%s,%s\n", subway.City, subway.Line,
						subway.Name, jsonInfo.Geocodes[0].Location))
					if err != nil {
						logger.Sugar.Error(err)
						break
					}
				} else {
					if len(jsonInfo.Geocodes) > 0 {
						logger.Sugar.Errorf("地理编码错误，poi:%s,location:%s,formatted_address:%s", subway,
							jsonInfo.Geocodes[0].Location, jsonInfo.Geocodes[0].FormattedAddress)
					} else {
						logger.Sugar.Errorf("地理编码错误，poi:%s,没有结果", subway)
					}

					_, err := fsWrite.WriteString(fmt.Sprintf("%s,%s,%s,\n", subway.City, subway.Line, subway.Name))
					if err != nil {
						logger.Sugar.Error(err)
						break
					}
				}
			} else {
				logger.Sugar.Error("status is not 1")
			}
		}
	}

	fsWrite.Flush()
	fs.Sync()
}
