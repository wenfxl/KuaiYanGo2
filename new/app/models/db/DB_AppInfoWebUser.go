package db

type DB_AppInfoWebUser struct {
	Id             int    `json:"id" gorm:"column:id;primarykey;autoIncrement:false"` // 关联appInfoAppId
	Status         int    `json:"status" gorm:"column:status;default:2;comment:状态(1>启用,2>停用)"`
	UrlDownload    string `json:"urlDownload"  gorm:"column:urlDownload;size:1000;comment:下载地址"`
	CaptchaLogin   int    `json:"captchaLogin"  gorm:"column:captchaLogin;default:3;comment:登陆防爆破起始次数"`
	CaptchaReg     int    `json:"captchaReg"  gorm:"column:captchaReg;default:2;comment:注册是否要验证码"`          // 1需要 2不需要
	CaptchaSendSms int    `json:"captchaSendSms"  gorm:"column:captchaSendSms;default:1;comment:发短信是否要验证码"` // 1需要 2不需要
}

func (DB_AppInfoWebUser) TableName() string {
	return "db_app_info_web_user"
}
