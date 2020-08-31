package modes

import (
	"fmt"
	"time"
)

/*
 * 描述: 用户管理认证表
 * status  0正常 1注销
 ********************************************************************/
type AdminUser struct {
	Id     int64  `json:"emp_id" xorm:"id"`           // 员 工 id
	RoleId int64  `json:"emp_role_id" xorm:"role_id"` // 角 色 id
	TeamId int64  `json:"emp_team_id" xorm:"team_id"` // 小 组 id
	ProjId int64  `json:"emp_proj_id" xorm:"proj_id"` // 项 目 id
	Name   string `json:"emp_name" xorm:"name"`       // 员工名称
	Phone  string `json:"emp_phone" xorm:"phone"`     // 账    号
	Pass   string `json:"-" xorm:"pass"`              // 密    码
	TQAcc  int    `json:"emp_tq_acc" xorm:"tq_acc"`   // TQ座席编号
	TQPw   string `json:"emp_tq_pw" xorm:"tq_pw"`     // TQ 秘 密
	Start  int    `json:"emp_start" xorm:"start"`     // 状    态
	At     int64  `json:"pro_at" xorm:"at"`           // 创建时间
}

func (this *AdminUser) TableName() string {
	return "employee"
}

func (this *AdminUser) Save() (int64, error) {
	this.At = time.Now().Unix()
	return Db(0).Insert(this)
}

func (this *AdminUser) Get() (bool, error) {
	if _, err := Db(0).Get(this); err != nil {
		return false, err
	}
	if this.RoleId%100 == 0 {
		return true, nil
	}
	if this.RoleId%100 == 4 {
		return true, nil
	}
	return false, nil

}

func (this *AdminUser) update(where string, field string) (int64, error) {
	return Db(0).Where(where).Cols(field).Update(this)
}

func (this *AdminUser) IdSet(field string) (int64, error) {
	where := fmt.Sprintf("id = %d", this.Id)
	return this.update(where, field)
}
