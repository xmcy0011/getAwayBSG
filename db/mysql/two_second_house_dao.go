package mysql

import (
	"github.com/getAwayBSG/configs"
	"time"
)

type TwoSecondHoseDao struct {
	tableName string
}

// 二手房信息
type HouseInfo struct {
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

func init() {
	DefaultTwoSecondHoseDao.tableName = configs.ConfigInfo.MysqlDbName

	// create table

}

// 增加一条房屋概要信息
func (t *TwoSecondHoseDao) Add(title string, totalPrice int, unitPrice int, link string) {
	dbMaster := DefaultManager.GetDbMaster()
	if dbMaster != nil {
		//sql := fmt.Sprintf("insert into %s(Title,ListCrawlTime,DetailStatus,TotalPrice,UnitPrice,Link," +
		//	"created,updated) values(%d,'%s',%d,%d,%d,%d,'%s',%d,%d,%d,%d,%d)",
		//	t.tableName, msgId, clientMsgId, fromId, toId, 0, msgType, msgData, cim.CIMResCode_kCIM_RES_CODE_OK,
		//	cim.CIMMsgFeature_kCIM_MSG_FEATURE_DEFAULT, cim.CIMMsgStatus_kCIM_MSG_STATUS_NONE, timeStamp, timeStamp)
		//_, err := dbMaster.Exec(sql)
		//if err != nil {
		//	logger.Sugar.Errorf("exec failed,sql:%s,error:%s", sql, err.Error())
		//	return 0, err
		//}
		//
		//session.Exec()
	}
}
