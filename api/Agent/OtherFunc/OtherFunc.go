package OtherFunc

import (
	"EFunc/utils"
	"github.com/gin-gonic/gin"
	"server/Service/Ser_AppInfo"
	"server/Service/Ser_AppUser"
	"server/Service/Ser_Log"
	"server/new/app/logic/common/agent"
	"server/new/app/logic/common/agentLevel"
	"server/structs/Http/response"
	DB "server/structs/db"
)

type Api struct{}

type 结构请求_修改用户绑定信息 struct {
	AppId int    `json:"AppId"`
	User  string `json:"User"`
	Key   string `json:"Key"`
}

// 修改软件用户绑定信息
func (a *Api) SetAppUserKey(c *gin.Context) {
	var 请求 结构请求_修改用户绑定信息
	err := c.ShouldBindJSON(&请求)
	if err != nil {
		response.FailWithMessage("提交参数错误:"+err.Error(), c)
		return
	}

	if !Ser_AppInfo.AppId是否存在(请求.AppId) {
		response.FailWithMessage("应用不存在", c)
		return
	}

	局_可操作AppId := agent.L_agent.Id取代理可操作应用AppId列表(c, c.GetInt("Uid"))
	if !utils.S数组_整数是否存在(局_可操作AppId, 请求.AppId) {
		response.FailWithMessage("无该应用操作权限,请联系上级授权该应用任意制卡卡类,获取应用权限", c)
		return
	}

	AppUserid := Ser_AppUser.User或卡号取Id(请求.AppId, 请求.User)

	if AppUserid == 0 {
		response.FailWithMessage("用户不存在", c)
		return
	}

	局_用户详情, err2 := Ser_AppUser.Id取详情(请求.AppId, AppUserid)
	if err2 != nil {
		response.FailWithMessage(err2.Error(), c)
		return
	}

	err = Ser_AppUser.Set绑定信息(请求.AppId, 局_用户详情.Uid, 请求.Key)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	局_信息 := "修改绑定信息 '" + 局_用户详情.Key + "'  ->  '" + 请求.Key + "'"

	Ser_Log.Log_写代理操作日志(c.GetInt("Uid"), agentLevel.L_agentLevel.Q取Id代理级别(c, c.GetInt("Uid")), 请求.AppId, AppUserid, Ser_AppUser.Id取User(请求.AppId, AppUserid), DB.D代理功能_修改用户绑定, c.ClientIP(), 局_信息)
	response.OkWithMessage("操作成功", c)
	return
}
