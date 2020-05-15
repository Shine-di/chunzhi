/**
 * @author: D-S
 * @date: 2020/4/16 3:26 下午
 */

package domain

import (
	"encoding/json"
	"gitee.com/risewinter/data-basic/app/domain/model/entity"
	"gitee.com/risewinter/data-basic/app/domain/repository"
	"gitee.com/risewinter/data-common/library/mysql"
	"gitee.com/risewinter/data-common/model/enum"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"net/http"
)

func Routers(engine *gin.Engine) {

	engine.POST("/api", func(context *gin.Context) {
		var equipment entity.Equipment
		err := context.Bind(&equipment)
		if err != nil {
			glog.Error(err)
			return
		}
		equipment.Deleted = enum.RESOURCE_DEL_STATUS_NORMAL
		equipment.Status = enum.RESOURCE_ENA_STATUS_NORMAL
		equipment.Audit = enum.RESOURCE_AUD_STATUS_AUDITED
		result, err := repository.NewEquipmentRepository().CreateEquipment(mysql.Conn(), &equipment)
		_ = result
		if err != nil {
			glog.Error(err)
			glog.Info("创建失败")
			b, _ := json.Marshal(equipment)
			glog.Info(string(b))
			context.JSON(http.StatusOK, "创建失败")
			return
		}
		glog.Info("成功")
		context.JSON(http.StatusOK, "success")
	})
}
