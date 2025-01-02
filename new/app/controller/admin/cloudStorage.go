package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"server/new/app/controller/Common"
	"server/new/app/logic/common/cloudStorage"
	"server/new/app/models/common"
	"server/new/app/models/request"
	. "server/new/app/models/response"
	"server/structs/Http/response"
)

// CloudStorage
// @MenuName 扩展
// @ModuleName 云存储
type CloudStorage struct {
	Common.Common
}

func NewCloudStorageController() *CloudStorage {
	var C = CloudStorage{}
	return &C
}

// Index
// @action 文件列表
// @show  1
func (C *CloudStorage) GetList(c *gin.Context) {
	var 请求 struct {
		request.List
		Path      string `json:"Path"`
		Delimiter string `json:"Delimiter"`
	}
	if !C.ToJSON(c, &请求) {
		return
	}
	var 局_list []common.W文件对象详情
	var err error
	局_list, err = cloudStorage.L_云存储.H获取文件列表(c, 请求.Path, 请求.Delimiter)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(GetList{List: 局_list, Count: int64(len(局_list))}, "操作成功", c)
	return
	//继续对接前端
}

// Index
// @action 取文件上传授权
// @show  1
func (C *CloudStorage) GetUpToken(c *gin.Context) {
	var 请求 struct {
		Path string `json:"Path"`
	}
	if !C.ToJSON(c, &请求) {
		return
	}
	var 局_list common.W文件上传凭证
	var err error
	局_list, err = cloudStorage.L_云存储.Q取文件上传授权(c, 请求.Path)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(局_list, "操作成功", c)
	return
	//继续对接前端
}

// Index
// @action 移动文件
// @show  1
func (C *CloudStorage) MoveTo(c *gin.Context) {
	var 请求 struct {
		Path1 string `json:"Path1"`
		Path2 string `json:"Path2"`
	}
	if !C.ToJSON(c, &请求) {
		return
	}
	var 局_list common.W文件上传凭证
	var err error
	err = cloudStorage.L_云存储.W文件移动(c, 请求.Path1, 请求.Path2)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(局_list, "操作成功", c)
	return
}

// Delete
// @action 删除
// @show  2
func (C *CloudStorage) Delete(c *gin.Context) {
	var 请求 struct {
		Path []string `json:"Path"`
	}
	if !C.ToJSON(c, &请求) {
		return
	}
	var 局_list common.W文件上传凭证

	_, err := cloudStorage.L_云存储.W文件删除(c, 请求.Path)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(局_list, "操作成功", c)
	return
}

// X下载
// @action X下载
// @show  2
func (C *CloudStorage) Download(c *gin.Context) {
	var 请求 struct {
		Path string `json:"Path"`
	}
	if !C.ToJSON(c, &请求) {
		return
	}

	下载地址, err := cloudStorage.L_云存储.X下载(c, 请求.Path)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(下载地址, "操作成功", c)
	return
}

// Q取外链地址
// @action Q取外链地址
// @show  2
func (C *CloudStorage) GetDownloadUrl(c *gin.Context) {
	var 请求 struct {
		Path     string `json:"Path"`
		LongTime int64  `json:"LongTime"`
	}
	if !C.ToJSON(c, &请求) {
		return
	}

	下载地址, err := cloudStorage.L_云存储.Q取外链地址(c, 请求.Path, 请求.LongTime)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(下载地址, "操作成功", c)
	return
}

// Q取基础信息
// @action Q取基础信息
// @show  2
func (C *CloudStorage) GetBaseInfo(c *gin.Context) {
	var 请求 struct {
	}
	if !C.ToJSON(c, &请求) {
		return
	}

	基础信息, err := cloudStorage.L_云存储.Q取基础信息(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var 局_MAP gin.H
	_ = json.Unmarshal([]byte(基础信息), &局_MAP)
	response.OkWithDetailed(局_MAP, "操作成功", c)
	return
}

/*
// Info
// @action 查询
// @show  2
func (C *CloudStorage) Info(c *gin.Context) {
	var 请求 request.Id
	if !C.ToJSON(c, &请求) {
		return
	}

	var S = service.S_CloudStorage{}
	tx := *global.GVA_DB
	var info db.DB_Cron_log
	info, err := S.Info(&tx, 请求.Id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	response.OkWithDetailed(info, "操作成功", c)
	return
}


// DeleteBatch
// @action 删除批量维护
// @show  2
func (C *CloudStorage) DeleteBatch(c *gin.Context) {
	var 请求 struct {
		Type    int    `json:"Type" binding:"required,min=1"`
		Keyword string `json:"Keyword" `
	}
	if !C.ToJSON(c, &请求) {
		return
	}

	var 影响行数 int64
	var S = service.S_CloudStorage{}
	tx := *global.GVA_DB
	影响行数, err := S.DeleteType(&tx, 请求.Type, 请求.Keyword)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功,数量"+strconv.FormatInt(影响行数, 10), c)
	return
}
*/
