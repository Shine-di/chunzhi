/**
 * @author: D-S
 * @date: 2020/4/16 3:26 下午
 */

package main

//
//type RatGroupPush struct {
//	Id          int64             `json:"id"` //玩法id
//	GameId      int64             `json:"game_id"`
//	TemplateId  int64             `json:"template_id"`  //模板Id
//	RelatedId   int64             `json:"related_id"`   //玩法关联id
//	RelatedType int32             `json:"related_type"` // 1 联赛ID 2 系列赛Id 3 小局ID
//	Level       int32             `json:"level"`
//	BoNo        uint8             `json:"bo_no"`   //第几局
//	NameZh      string            `json:"name"`    // 玩法名称
//	NameEn      string            `json:"name_en"` //英文
//	Wight       uint8             `json:"wight"`   //排序权重
//	Data        []RateGroupStatus `json:"data"`
//}
//type RateGroupStatus struct {
//	Source      int32          `json:"source"`
//	IsRoulette  bool           `json:"is_roulette"` //是否有滚球
//	GroupStatus int32          `json:"group_status"`
//	Items       []RateItemPush `json:"items"`
//}
//
//type RateItemPush struct {
//	Id       int64  `json:"id"`        //选项Id
//	NameZh   string `json:"name"`      // 玩法名称
//	NameEn   string `json:"name_en"`   //英文
//	Value    string `json:"value"`     //玩法选项的值
//	Rate     string `json:"rate"`      //选项的指数
//	InitRate string `json:"init_rate"` //选项的初始指数
//	TeamId   int64  `json:"team	_id"`  //队伍ID
//	IsWin    int    `json:"is_win"`    //该选项是否获胜  0 选项没有结算，没有胜负 1 胜 2 负
//}
//
//type Model struct {
//	Id         int64 `json:"id"`          //主键
//	CreateTime int64 `json:"create_time"` //创建时间
//	UpdateTime int64 `json:"update_time"` //更新时间
//	Deleted    int32 `json:"deleted"`     //是否删除 1 未删除 2 已删除
//	Status     int32 `json:"status"`      //1.正常 2.禁用
//	Audit      int32 `json:"audit"`       //1.未审核 2.已审核 3.已拒绝
//}
//
//type BaseData struct {
//	Model
//	GameId      int64  `json:"game_id"`       //游戏id
//	Name        string `json:"name"`          //名称
//	NameZh      string `json:"name_zh"`       //中文名称
//	NameZhShort string `json:"name_zh_short"` //中文简称
//	NameEn      string `json:"name_en"`       //英文名称
//	NameEnShort string `json:"name_en_short"` //英文简称
//	Icon        string `json:"icon"`          //头像
//	Description string `json:"description"`   //描述
//	Source      uint   `json:"source"`        //来源平台
//}
//
//func Sync() {
//	req := &data_game_proto.RequestGameEquipmentList{
//		Offset:     0,
//		Limit:      1000,
//		OrderValue: "id desc",
//	}
//	equipment, num, err := game_equipment.GameEquipmentList(context.Background(), req)
//	if err != nil {
//		log.Error(err)
//		return
//	}
//	var i = int(0)
//	fmt.Sprintln(fmt.Printf("num ------ %v", num))
//	for _, e := range equipment {
//		equi := new(BaseData)
//		copier.Copy(equi, e)
//		equi.Name = e.Name
//		equi.NameZhShort = e.ShortName
//		equi.NameEnShort = e.ShortName
//		equi.Icon = e.Icon
//		equi.GameId = e.GameId
//		equi.Description = e.Description
//		i++
//		fmt.Sprintln(fmt.Printf("num ------ %v------------------", i))
//		err, _ := Post("http://127.0.0.1:50052/api", *equi)
//		if err != nil {
//			fmt.Println("-----------------调用失败")
//			log.Error(err)
//			continue
//		}
//		fmt.Println("---------------调用成功")
//
//	}
//	fmt.Println("-----调用完成")
//}
//
//func Post(url string, body interface{}) (error, []byte) {
//	bodyByte, _ := json.Marshal(body)
//	client := http.Client{}
//	req, _ := http.NewRequest("POST", url, bytes.NewReader(bodyByte))
//	req.Header.Add("Content-Type", "application/json")
//	req.Header.Add("project", "space")
//	resp, err := client.Do(req)
//	if err != nil {
//		return err, nil
//	}
//	defer resp.Body.Close()
//	rb, err := ioutil.ReadAll(resp.Body)
//	fmt.Println(string(rb))
//	if string(rb) != "success" {
//		fmt.Println(string(bodyByte))
//		return errors.New("调用失败"), nil
//	}
//	return err, rb
//}
