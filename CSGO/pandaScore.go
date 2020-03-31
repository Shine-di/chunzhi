/**
* @Author: D-S
* @Date: 2020/3/20 10:37 上午
 */

package CSGO

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)


type BaseData struct {
	Model
	GameId      int64
	Name        string `gorm:"not null; column:name" json:"name"`
	NameEnglish string `gorm:"not null; column:name_english" json:"name_english"`
	ShortName   string `gorm:"not null; column:short_name" json:"short_name"`
	Icon        string `gorm:"not null; column:icon" json:"icon"`
	Description string `gorm:"not null; column:description" json:"description"`
	From   uint `json:"from"`
}

//基础数据匹配规则
type ExBaseData struct {
	ExId  string         `json:"ex_id"` //外部Id
	From  int            `json:"from"`  //来源平台
	PId   int64          `json:"p_id"`  //内部Id(人工维护)
	Alias pq.StringArray `json:"alias"` //名称匹配池 初始name nameEn shortName(人工维护)
}


//外部联赛
type ExRateLeague struct {
	ExBaseData
	BaseData
}

//外部系列赛
type ExRateSerise struct {
	ExBaseData
	StartTime  int64  `gorm:"not null; column:start_time" json:"start_time"`
	LeagueId   int64  `gorm:"not null; column:league_id" json:"league_id"`
	SeasonInfo string `gorm:"not null; column:season_info" json:"season_info"` //BO3
	TeamIds  pq.StringArray `json:"team_ids"`
}

//外部team
type ExRateTeam struct {
	ExBaseData
}

//内部玩法
type RateGroup struct {
	BaseData
	Stage          int            `json:"stage"`          //第几场 全局玩法值为100
}

//外部玩法
type ExRateGroup struct {
	BaseData
	Stage          int            `json:"stage"`          //第几场 全局玩法值为100
	PId int64 `json:"p_id"`
	Level  uint   `json:"level"`   // 1 联赛 2系列赛 3小局
	ExGroupStatus    int `json:"ex_group_status"`         //玩法状态  开盘/锁盘/封盘/结算
	GroupStatus int `json:"group_status"`
}

//内部玩法分类
type RateGroupType struct {
	BaseData
	Level  uint   `json:"level"`   // 1 联赛 2系列赛 3小局
	Wight int `json:"wight"` //排序权重
}



//内部选项
type RateItem struct {
	RateGroupId int64 `json:"rate_group_id"`
	ValueType uint `json:"value_type"` //1.常量  2.战队名称 3.数值 4.常量+数值
	Value          string         `json:"value"`          //玩法选项的值 %s %d-%d %s  胜
}

//外部选项
type ExRateItem struct {
	ExBaseData
	RateGroupId string `json:"rate_group_id"`
	Value          string         `json:"value"`          //玩法选项的值
	Rate           string         `json:"rate"`           //选项的指数
	InitRate       string         `json:"initRate"`       //选项的初始指数
	IsWin          bool           `json:"isWin"`          //是否获胜
}


type RateItem struct {
	gorm.Model
	LeagueId       int            `json:"leagueId"`       //外部联赛id
	SeriesId       int64          `json:"seriesId"`       //外部系列赛id
	ItemId         int64          `json:"itemId"`         //玩法选项id
	GroupId        int64          `json:"groupId"`        //玩法id
	GroupName      int64          `json:"groupName"`      //玩法名字
	GroupShortName string         `json:"groupShortName"` //玩法名字简写
	GroupStatus    int            `json:"status"`         //玩法状态 开盘/锁盘/封盘/结算
	From           int            `json:"from"`           //来源平台
	Value          string         `json:"value"`          //玩法选项的值
	Rate           string         `json:"rate"`           //选项的指数
	InitRate       string         `json:"initRate"`       //选项的初始指数
	IsWin          bool           `json:"isWin"`          //是否获胜
	PLeagueId      int64          `json:"pLeagueId"`      //潘多拉联赛id
	PSeriesId      int64          `json:"pSeriesId"`      //潘多拉系列赛id
	PItemId        int64          `json:"pItemId"`        //潘多拉
	Var            pq.StringArray `json:"var"`            //
}

type RateGroupRelation struct {
	gorm.Model
	From     int    `json:"from"`     //来源平台
	Name     string `json:"name"`     //外部名称
	PName    string `json:"pName"`    //内部名称
	PGroupId int    `json:"pGroupId"` //内部玩法id
}

type RateGroupType struct {
	gorm.Model
	GameId uint8  `json:"game_id"` //游戏id
	Name   string `json:"name"`    //玩法名称
	Level  uint   `json:"level"`   // 1 联赛 2系列赛 3小局
	Status int    `json:"status"`  //玩法状态 1可用 2禁用
	Wight int `json:"wight"` //排序权重
}

//玩法选项
type RateGroupItemType struct {
	gorm.Model
	GroupTypeId    int    `json:"group_type_id"`
	Title          string `json:"title"`
	CategoryTypeId uint   `json:"category_type_id"` //选项分类id
	Template       string `json:"template"`         //模板 egg.  %s,%s   WE 战队2:0 IG战队  %s %s:%s %s    WE战队胜  %s战队胜
}

//玩法选项分类
type RateGroupItemCategoryType struct {
	gorm.Model
	Title string `json:"title"`
}