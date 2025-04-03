package controller

import (
	. "EFunc/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/buffer"
	"server/Service/Ser_AppInfo"
	"server/Service/Ser_KaClass"
	"server/Service/Ser_User"
	"server/global"
	"server/new/app/controller/Common"
	"server/new/app/logic/common/agent"
	"server/structs/Http/response"
	DB "server/structs/db"
	"strconv"
)

type AgentUser struct {
	Common.Common
}

func NewAgentUserController() *AgentUser {
	return &AgentUser{}
}

// 销售统计
func (C *AgentUser) GetKaSalesStatistics(c *gin.Context) {
	var 请求 struct {
		AppId        int      `json:"AppId"`
		Num          int      `json:"Num"`
		RegisterTime []string `json:"RegisterTime"`
		UseTime      []string `json:"UseTime"`
		KaClassId    int      `json:"KaClassId"`
		AgentLv      int      `json:"AgentLv"`
		AgentName    string   `json:"AgentName"`
	}
	//解析失败
	if !C.ToJSON(c, &请求) {
		return
	}
	var info = struct {
		appInfo   DB.DB_AppInfo
		DB_Ka     []DB.DB_Ka
		总数        int64
		局_制卡人     []string
		卡类id名称map map[int]string
	}{
		卡类id名称map: make(map[int]string), // 新增初始化代码 // [!code ++]
	}

	tx := *global.GVA_DB

	if 请求.AppId < 10000 && 请求.AppId != 0 {
		response.FailWithMessage("AppId请输>=10000的整数", c)
		return
	}
	info.appInfo = Ser_AppInfo.App取App详情(请求.AppId)
	//和管理员的逻辑略微不一样 不指定,就查全部
	if 请求.AgentName != "" {
		局_代理详情, ok := Ser_User.User取详情(请求.AgentName)
		if !ok {
			response.FailWithMessage("代理不存在", c)
			return
		}
		if 局_代理详情.UPAgentId != c.GetInt("Uid") && 局_代理详情.Id != c.GetInt("Uid") {
			response.FailWithMessage("无权限查看其它代理信息", c)
			return
		}

		info.局_制卡人 = append(info.局_制卡人, 请求.AgentName)
	} else {
		info.局_制卡人 = append(info.局_制卡人, c.GetString("User")) // 先添加自身
		info.局_制卡人 = append(info.局_制卡人, agent.L_agent.Q取下级代理数组_user(c, []int{c.GetInt("Uid")})...)
	}
	局_DB := tx.Model(DB.DB_Ka{})

	if 请求.AppId != 0 {
		局_DB.Where("AppId = ?", 请求.AppId)
	}

	if 请求.Num == 1 || 请求.Num == 2 {
		switch 请求.Num {
		case 1: //已经使用
			局_DB.Where("Num = NumMax")
		case 2: //未使用过
			局_DB.Where("Num < NumMax")
		}
	}
	if 请求.RegisterTime != nil && len(请求.RegisterTime) == 2 && 请求.RegisterTime[0] != "" && 请求.RegisterTime[1] != "" {
		制卡开始时间, _ := strconv.Atoi(请求.RegisterTime[0])
		制卡结束时间, _ := strconv.Atoi(请求.RegisterTime[1])
		局_DB.Where("RegisterTime > ?", 制卡开始时间).Where("RegisterTime < ?", 制卡结束时间+86400)
	}

	if 请求.UseTime != nil && len(请求.UseTime) == 2 && 请求.UseTime[0] != "" && 请求.UseTime[1] != "" {
		使用开始时间, _ := strconv.Atoi(请求.UseTime[0])
		使用结束时间, _ := strconv.Atoi(请求.UseTime[1])
		局_DB.Where("UseTime > ?", 使用开始时间).Where("UseTime < ?", 使用结束时间+86400)
	}
	if 请求.KaClassId != 0 {
		局_DB.Where("KaClassId = ?", 请求.KaClassId)
	}
	//直接限制 仅代理自身和直属下级
	var 局_制卡人数组 = []string{c.GetString("User")}
	局_制卡人数组 = append(局_制卡人数组, agent.L_agent.Q取下级代理数组_user(c, []int{c.GetInt("Uid")})...)

	//直接限制 仅代理自身和直属下级
	局_临时数组 := make([]string, 0, len(info.局_制卡人))
	for i := range info.局_制卡人 {
		if S数组_是否存在(局_制卡人数组, info.局_制卡人[i]) {
			局_临时数组 = append(局_临时数组, info.局_制卡人[i])
		}
	}
	info.局_制卡人 = 局_临时数组

	//过滤制卡人
	if len(info.局_制卡人) == 0 {
		info.局_制卡人 = append(info.局_制卡人, c.GetString("User")) // 先添加自身
	}

	局_DB.Where("RegisterUser in (?)", info.局_制卡人)

	//Count(&总数) 必须放在where 后面 不然值会被清0
	err := 局_DB.Count(&info.总数).Find(&info.DB_Ka).Error
	//fmt.Println("用户总数%d", 总数, DB_LinksToken)
	if err != nil {
		response.FailWithMessage("查询失败,参数异常"+err.Error(), c)
		return
	}

	if info.总数 == 0 {
		response.FailWithMessage("查询失败,无符合条件卡号,无法统计", c)
		return
	}

	局_制卡人列表 := []string{}
	for _, v := range info.DB_Ka {
		if S数组_取文本出现次数(局_制卡人列表, v.RegisterUser) == 0 {
			//如果是自身,放到首位
			if v.RegisterUser == c.GetString("User") {
				局_制卡人列表 = append([]string{v.RegisterUser}, 局_制卡人列表...)
			} else {
				局_制卡人列表 = append(局_制卡人列表, v.RegisterUser)
			}
		}
	}

	var 局_快速文本对象 buffer.Buffer
	for _, item制卡人User := range 局_制卡人列表 {
		局_卡类map := make(map[int]int) // [!code ++]
		for i := range info.DB_Ka {
			if info.DB_Ka[i].RegisterUser == item制卡人User {
				//卡类id 统计数量
				if _, ok := 局_卡类map[info.DB_Ka[i].KaClassId]; ok {
					局_卡类map[info.DB_Ka[i].KaClassId]++
				} else {
					局_卡类map[info.DB_Ka[i].KaClassId] = 1
				}
			}
		}
		局_快速文本对象.AppendString("\n========[代理]:" + item制卡人User + "========")
		局_快速文本对象.AppendString("\n")
		局_快速文本对象.AppendString("\n[应用名称]:" + info.appInfo.AppName + "")
		for 局_卡类id, 局_卡类id数量 := range 局_卡类map {
			info.卡类id名称map[局_卡类id] = Ser_KaClass.Id取Name(局_卡类id)
			局_快速文本对象.AppendString("   [" + info.卡类id名称map[局_卡类id] + "]:" + strconv.Itoa(局_卡类id数量) + "")
		}
		局_快速文本对象.AppendString("\n-----------------------------------------")
		局_计数 := 1
		for i := range info.DB_Ka {
			if info.DB_Ka[i].RegisterUser == item制卡人User {
				局_快速文本对象.AppendString("\n" + strconv.Itoa(局_计数) + ":" + info.DB_Ka[i].Name)
				局_快速文本对象.AppendString("  [制卡时间]:" + S时间_时间戳到时间(info.DB_Ka[i].RegisterTime))
				局_快速文本对象.AppendString("  [使用时间]:" + S时间_时间戳到时间(info.DB_Ka[i].UseTime))
				局_快速文本对象.AppendString("  [卡类名称]:" + info.卡类id名称map[info.DB_Ka[i].KaClassId])
				局_快速文本对象.AppendString("  [状态]:" + S三元(info.DB_Ka[i].Status == 1, "正常", "冻结"))
				局_快速文本对象.AppendString("  [管理备注]:" + info.DB_Ka[i].AdminNote)
				局_快速文本对象.AppendString("  [代理备注]:" + info.DB_Ka[i].AgentNote)
				局_计数++
			}
		}
		局_快速文本对象.AppendString("\n")
	}

	response.OkWithDetailed(局_快速文本对象.String(), "获取成功", c)
	return
}
