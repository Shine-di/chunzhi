/**
 * @author: D-S
 * @date: 2020/4/17 8:47 下午
 */

package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	data_game_proto "gitee.com/risewinter/data-common/proto/data-game"
	"gitee.com/risewinter/data-game/library/log"
	game_equipment "gitee.com/risewinter/data-game/page/service/game-equipment"
	"github.com/jinzhu/copier"
	"io/ioutil"
	"net/http"
)

type Model struct {
	Id         int64 `json:"id"`          //主键
	CreateTime int64 `json:"create_time"` //创建时间
	UpdateTime int64 `json:"update_time"` //更新时间
	Deleted    int32 `json:"deleted"`     //是否删除 1 未删除 2 已删除
	Status     int32 `json:"status"`      //1.正常 2.禁用
	Audit      int32 `json:"audit"`       //1.未审核 2.已审核 3.已拒绝
}

type BaseData struct {
	Model
	GameId      int64  `json:"game_id"`       //游戏id
	Name        string `json:"name"`          //名称
	NameZh      string `json:"name_zh"`       //中文名称
	NameZhShort string `json:"name_zh_short"` //中文简称
	NameEn      string `json:"name_en"`       //英文名称
	NameEnShort string `json:"name_en_short"` //英文简称
	Icon        string `json:"icon"`          //头像
	Description string `json:"description"`   //描述
	Source      uint   `json:"source"`        //来源平台
}


func Sync() {
	req := &data_game_proto.RequestGameEquipmentList{
		Offset:     0,
		Limit:      1000,
		OrderValue: "id desc",
	}
	equipment, num, err := game_equipment.GameEquipmentList(context.Background(), req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var i = int(0)
	fmt.Sprintln(fmt.Printf("num ------ %v", num))
	for _, e := range equipment {
		equi := new(BaseData)
		copier.Copy(equi, e)
		equi.Name = e.Name
		equi.NameZhShort = e.ShortName
		equi.NameEnShort = e.ShortName
		equi.Icon = e.Icon
		equi.GameId = e.GameId
		equi.Description = e.Description
		i++
		fmt.Sprintln(fmt.Printf("num ------ %v-------", i))
		err, _ := Post("http://127.0.0.1:50052/api", *equi)
		if err != nil {
			fmt.Println("-----------------调用失败")
			log.Error(err.Error())
			continue
		}
		fmt.Println("---------------调用成功")

	}
	fmt.Println("-----调用完成")
}

func Post(url string, body interface{}) (error, []byte) {
	bodyByte, _ := json.Marshal(body)
	client := http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewReader(bodyByte))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("project", "space")
	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	//fmt.Println(string(rb))
	//if string(rb) != "success" {
	//	fmt.Println(string(bodyByte))
	//	return errors.New("调用失败"), nil
	//}
	return nil, rb
}
