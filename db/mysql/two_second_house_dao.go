package mysql

import (
	"errors"
	"fmt"
	"github.com/getAwayBSG/configs"
	"github.com/getAwayBSG/logger"
	"time"
)

type TwoSecondHoseDao struct {
	tableName string
}

// 二手房信息
type HouseInfo struct {
	Id            int       // 数据库自增编号
	Title         string    // 标题
	ListCrawlTime time.Time // 列表抓取时间
	DetailStatus  int       // 详情抓取状态
	TotalPrice    int       // 总价，单位：元
	UnitPrice     int       // 平米单价，单位：元
	Link          string    // 链家房屋详情url

	HouseRecordLJ   string  // 房源编号（唯一？）
	DetailCrawlTime string  // 详情抓取时间
	AreaName        string  // 区域
	BaseAttr        string  // 房屋基础属性，房屋户型:3室2厅1厨2卫|所在楼层:高楼层 (共32层)|建筑面积:103.92㎡|户型结构:平...
	BeOnlineTime    string  // 上架时间
	CompletedInfo   string  // 竣工时间
	DecorateInfo    string  // 装修,平层/精装
	DirectionInfo   string  // 朝向,南北
	FloorInfo       string  // 楼层信息
	RoomInfo        string  // 房屋大小
	Size            float32 // 大小
	TransactionAttr string  // 交易属性
	VillageName     string  // 小区名称
}

// 默认实例
var DefaultTwoSecondHoseDao = &TwoSecondHoseDao{}

// 初始化，创建表等
func (t *TwoSecondHoseDao) Init() {
	DefaultTwoSecondHoseDao.tableName = configs.ConfigInfo.TwoHandHouseCollection

	// create table
	dbMaster := DefaultManager.GetDbMaster()
	if dbMaster != nil {
		sql := fmt.Sprintf("SELECT table_name FROM information_schema.TABLES WHERE table_name ='%s'",
			configs.ConfigInfo.TwoHandHouseCollection)
		row := dbMaster.QueryRow(sql)

		var tableName string
		err := row.Scan(&tableName)
		if err != nil {
			sql := fmt.Sprintf("CREATE TABLE `%s` ("+
				"`id` int(11) unsigned NOT NULL AUTO_INCREMENT,"+
				"`Title` varchar(200) DEFAULT NULL,"+
				"`ListCrawlTime` datetime DEFAULT NULL,"+
				"`DetailStatus` tinyint(4) DEFAULT NULL,"+
				"`TotalPrice` int(11) DEFAULT NULL,"+
				"`UnitPrice` int(11) DEFAULT NULL,"+
				"`Link` varchar(128) DEFAULT NULL,"+
				"`HouseRecordLJ` varchar(32) DEFAULT NULL,"+
				"`DetailCrawlTime` datetime DEFAULT NULL,"+
				"`AreaName` varchar(64) DEFAULT NULL,"+
				"`BaseAttr` varchar(500) DEFAULT NULL,"+
				"`BeOnlineTime` varchar(64) DEFAULT NULL,"+
				"`CompletedInfo` varchar(64) DEFAULT NULL,"+
				"`DecorateInfo` varchar(128) DEFAULT NULL,"+
				"`DirectionInfo` varchar(32) DEFAULT NULL,"+
				"`FloorInfo` varchar(32) DEFAULT NULL,"+
				"`RoomInfo` varchar(32) DEFAULT NULL,"+
				"`Size` float DEFAULT NULL,"+
				"`TransactionAttr` varchar(128) DEFAULT NULL,"+
				"`VillageName` varchar(64) DEFAULT NULL,"+
				"PRIMARY KEY (`id`)"+
				") ENGINE=InnoDB DEFAULT CHARSET=utf8;",
				t.tableName)
			_, err := dbMaster.Exec(sql)
			if err != nil {
				logger.Sugar.Error(err)
			}
		}
	}
}

// 增加一条房屋概要信息
func (t *TwoSecondHoseDao) Add(title string, totalPrice int, unitPrice int, link string) {
	dbMaster := DefaultManager.GetDbMaster()
	if dbMaster != nil {
		var timeStr = time.Now().Format("2006-01-02 15:04:05")
		sql := fmt.Sprintf("insert into %s(Title,ListCrawlTime,DetailStatus,TotalPrice,UnitPrice,Link)"+
			" values('%s','%s',0,%d,%d,'%s')",
			t.tableName, title, timeStr, totalPrice, unitPrice, link)
		_, err := dbMaster.Exec(sql)
		if err != nil {
			logger.Sugar.Errorf("exec failed,sql:%s,error:%s", sql, err.Error())
		}
	}
}

// 更新详细信息
func (t *TwoSecondHoseDao) Update(id int, detailStatus int, info HouseInfo) error {
	dbMaster := DefaultManager.GetDbMaster()
	if dbMaster != nil {
		sql := fmt.Sprintf("update %s set DetailStatus=%d,HouseRecordLJ='%s',DetailCrawlTime='%s',AreaName='%s',"+
			"BaseAttr='%s',BeOnlineTime='%s',CompletedInfo='%s',DecorateInfo='%s',DirectionInfo='%s',FloorInfo='%s',RoomInfo='%s',Size=%f,"+
			"TransactionAttr='%s',VillageName='%s' where id=%d",
			t.tableName, detailStatus, info.HouseRecordLJ, info.DetailCrawlTime, info.AreaName,
			info.BaseAttr, info.BeOnlineTime, info.CompletedInfo, info.DecorateInfo, info.DirectionInfo, info.FloorInfo,
			info.RoomInfo, info.Size, info.TransactionAttr, info.VillageName, id)
		_, err := dbMaster.Exec(sql)
		if err != nil {
			logger.Sugar.Error(err)
		}
		return err
	}
	return errors.New("db disconnected")
}
func (t *TwoSecondHoseDao) Update2(id int, detailStatus int) error {
	dbMaster := DefaultManager.GetDbMaster()
	if dbMaster != nil {
		sql := fmt.Sprintf("update %s set DetailStatus=%d where id=%d",
			t.tableName, detailStatus, id)
		_, err := dbMaster.Exec(sql)
		if err != nil {
			logger.Sugar.Error(err)
		}
		return err
	}
	return errors.New("db disconnected")
}

func (t *TwoSecondHoseDao) GetAll(detailStatus int) ([]HouseInfo, error) {
	dbMaster := DefaultManager.GetDbMaster()
	if dbMaster != nil {
		sql := fmt.Sprintf("select Id,Title,DetailStatus,Link from %s where DetailStatus=%d", t.tableName, detailStatus)
		rows, err := dbMaster.Query(sql)
		if err != nil {
			logger.Sugar.Error(err)
			return nil, err
		}
		defer rows.Close()
		var houseArr = make([]HouseInfo, 0)

		for rows.Next() {
			var house = HouseInfo{}
			err := rows.Scan(&house.Id, &house.Title, &house.DetailStatus, &house.Link)
			if err != nil {
				logger.Sugar.Error(err)
				break
			}
			houseArr = append(houseArr, house)
		}
		return houseArr, nil
	}
	return nil, errors.New("db disconnected")
}
